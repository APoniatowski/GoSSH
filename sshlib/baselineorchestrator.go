package sshlib

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// stepTimeout bounds how long a single fanned-out step waits for the fleet.
// Generous on purpose: package installs over slow links are legitimate. A host
// that exceeds it is declared dead so one hung box cannot stall the fleet.
const stepTimeout = 120 * time.Second

// baselineStep is one ordered unit of work: a per-host command map plus a label
// for reporting. Hosts absent from cmds (or mapped to "") are no-ops that round.
type baselineStep struct {
	label     string
	cmds      map[string]string // host -> command
	collectAs string            // when non-empty, each host's output is written to ./collections/<host>/<collectAs>
	invert    bool              // when true, raw statusOK means non-compliant (must-not-have probes where OK=present=violation)
}

// baselineMode selects the apply vs. read-only check command builders.
type baselineMode int

const (
	modeApply baselineMode = iota
	modeCheck
)

// baselineFleet holds the live per-host workers for one server group and the
// shared results channel they report on.
type baselineFleet struct {
	workers map[string]*baselineWorker // key: host fqdn
	results chan Result
	dead    map[string]bool
}

// newBaselineFleet dials every host in sshList concurrently, waits for all dial
// attempts to finish (a WaitGroup barrier, never len(chan)), starts one worker
// goroutine per live host, and records hosts that failed to connect as dead.
func newBaselineFleet(sshList map[string]string, pool map[string]ParsedPool) *baselineFleet {
	type dialed struct {
		host string
		exec baselineExecutor
		err  error
	}
	var wg sync.WaitGroup
	out := make(chan dialed, len(sshList))
	for host := range sshList {
		h := host
		pp, ok := pool[h]
		if !ok {
			out <- dialed{host: h, err: fmt.Errorf("no connection info for %s", h)}
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			exec, err := pp.dialExecutor(h)
			out <- dialed{host: h, exec: exec, err: err}
		}()
	}
	wg.Wait()
	close(out)

	f := &baselineFleet{
		workers: make(map[string]*baselineWorker),
		results: make(chan Result, len(sshList)),
		dead:    make(map[string]bool),
	}
	for d := range out {
		if d.err != nil {
			fmt.Printf("  [%s] connection failed: %v\n", d.host, d.err)
			f.dead[d.host] = true
			continue
		}
		w := &baselineWorker{host: d.host, exec: d.exec, in: make(chan Command, 1), results: f.results}
		f.workers[d.host] = w
		go w.run()
	}
	return f
}

// newFleetWithExecutors builds a fleet from pre-supplied executors. Test-only
// seam: lets unit tests drive the orchestrator with fakes, no SSH.
func newFleetWithExecutors(execs map[string]baselineExecutor) *baselineFleet {
	f := &baselineFleet{
		workers: make(map[string]*baselineWorker),
		results: make(chan Result, len(execs)),
		dead:    make(map[string]bool),
	}
	for host, exec := range execs {
		w := &baselineWorker{host: host, exec: exec, in: make(chan Command, 1), results: f.results}
		f.workers[host] = w
		go w.run()
	}
	return f
}

// dispatch fans one step out to every live worker and gathers exactly one Result
// per worker, with a timeout. Hosts that error at the transport level or time out
// are moved into f.dead and skipped on subsequent steps, so a single hung or dead
// host can never deadlock the fleet.
func (f *baselineFleet) dispatch(cmds map[string]string, timeout time.Duration) map[string]Result {
	sent := 0
	for host, w := range f.workers {
		if f.dead[host] {
			continue
		}
		w.in <- Command{Host: host, Cmd: cmds[host]} // inbox cap 1; never blocks the fleet
		sent++
	}
	out := make(map[string]Result, sent)
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for i := 0; i < sent; i++ {
		select {
		case r := <-f.results:
			out[r.Host] = r
			if r.Err != nil {
				f.dead[r.Host] = true
			}
		case <-timer.C:
			for host := range f.workers {
				if _, got := out[host]; !got && !f.dead[host] {
					f.dead[host] = true
					out[host] = Result{Host: host, Status: statusNOK, Err: fmt.Errorf("timed out")}
				}
			}
			return out
		}
	}
	return out
}

// shutdown closes every worker inbox; each worker drains, closes its connection,
// and exits. Safe to call once after all steps are dispatched.
func (f *baselineFleet) shutdown() {
	for _, w := range f.workers {
		close(w.in)
	}
}

// liveHosts returns the still-alive subset of sshList (host -> os), excluding
// hosts that failed to connect or died mid-run.
func (f *baselineFleet) liveHosts(sshList map[string]string) map[string]string {
	live := make(map[string]string, len(sshList))
	for host, os := range sshList {
		if !f.dead[host] {
			live[host] = os
		}
	}
	return live
}

// withSudo prefixes a command with sudo unless the host logs in as root.
func withSudo(isRoot bool, cmd string) string {
	if isRoot {
		return cmd
	}
	return "sudo " + cmd
}

// buildCmds builds a per-host command map. fn receives each host's OS string.
// When isRoot is non-nil the result is sudo-wrapped per host (apply); pass
// isRoot == nil for no-sudo check steps.
func buildCmds(sshList map[string]string, isRoot map[string]bool, fn func(os string) string) map[string]string {
	cmds := make(map[string]string)
	for host, os := range sshList {
		cmd := fn(os)
		if isRoot != nil {
			cmd = withSudo(isRoot[host], cmd)
		}
		cmds[host] = cmd
	}
	return cmds
}

// collectStep writes each alive host's captured output to a local artifact file
// at ./collections/<host>/<collectAs>. The collectAs is sanitized to a base
// filename so it cannot escape the per-host directory. Write failures are warned
// about but never abort the run.
func collectStep(collectAs string, results map[string]Result, fleet *baselineFleet) {
	name := filepath.Base(filepath.Clean(collectAs))
	if name == "." || name == string(filepath.Separator) || name == ".." {
		fmt.Printf("    [collect] invalid artifact name %q, skipping\n", collectAs)
		return
	}
	for host, r := range results {
		if fleet.dead[host] || r.Err != nil {
			continue
		}
		dir := filepath.Join("collections", host)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Printf("    [%s] collect: cannot create %q: %v\n", host, dir, err)
			continue
		}
		dest := filepath.Join(dir, name)
		if err := os.WriteFile(dest, r.Output, 0o644); err != nil {
			fmt.Printf("    [%s] collect: cannot write %q: %v\n", host, dest, err)
			continue
		}
	}
}

// reportStep prints a one-line OK/NOK summary per host for a dispatched step.
func reportStep(label string, results map[string]Result) {
	if len(results) == 0 {
		return
	}
	for host, r := range results {
		state := "OK"
		if r.Status != statusOK || r.Err != nil {
			state = "NOK"
		}
		fmt.Printf("    [%s] %s: %s\n", host, label, state)
	}
}

// isRootMap derives per-host root-ness (empty or "root" username) from the pool.
// This replaces the legacy single group-wide bool, which was incorrect when a
// group mixed root and non-root logins.
func isRootMap(pool map[string]ParsedPool) map[string]bool {
	m := make(map[string]bool, len(pool))
	for host, pp := range pool {
		m[host] = pp.Username == "" || pp.Username == "root"
	}
	return m
}

// buildPool collects connection parameters for every server across all groups,
// keyed by fqdn, so workers can be created for any host in an sshList.
func buildPool(configs *yaml.MapSlice) map[string]ParsedPool {
	pool := make(map[string]ParsedPool)
	for _, groupItem := range *configs {
		groupValue, ok := groupItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
		}
		for _, serverItem := range groupValue {
			sv, ok := serverItem.Value.(yaml.MapSlice)
			if !ok {
				panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
			}
			pp := parseServer(sv)
			pp.defaulter()
			pool[pp.FQDN] = pp
		}
	}
	return pool
}

// runBaseline is the shared driver for both apply and check on one server group.
// It resolves the host list (after excludes), connects a fleet of persistent
// sessions, then walks the stages in order, fanning each step to the fleet and
// reporting results. Apply and check differ only in which stage builders run and
// the per-step command strings they emit.
func (baselineStruct *ParsedBaseline) runBaseline(serverGroupName string, configs *yaml.MapSlice, mode baselineMode) *ComplianceReport {
	report := &ComplianceReport{}
	var sshList map[string]string
	if mode == modeApply {
		sshList = baselineStruct.applyOSExcludes(serverGroupName, configs)
	} else {
		sshList = baselineStruct.checkOSExcludes(serverGroupName, configs)
	}
	if len(sshList) == 0 {
		fmt.Println("No servers to act on for this group (check excludes / pool).")
		return report
	}

	pool := buildPool(configs)
	fleet := newBaselineFleet(sshList, pool)
	defer fleet.shutdown()
	isRoot := isRootMap(pool)

	run := func(steps []baselineStep) {
		for _, st := range steps {
			results := fleet.dispatch(st.cmds, stepTimeout)
			if st.collectAs != "" {
				collectStep(st.collectAs, results, fleet)
			}
			reportStep(st.label, results)
			report.Add(serverGroupName, st.label, results, st.invert)
		}
	}

	var rebootBool bool
	if mode == modeApply {
		run(baselineStruct.applyPrereq(fleet.liveHosts(sshList), isRoot))
		run(baselineStruct.applyMustHaves(fleet.liveHosts(sshList), isRoot, &rebootBool))
		run(baselineStruct.applyMustNotHaves(fleet.liveHosts(sshList), isRoot))
		run(baselineStruct.applyFinals(fleet.liveHosts(sshList), isRoot, &rebootBool))
	} else {
		run(baselineStruct.checkPrereqs(fleet.liveHosts(sshList)))
		run(baselineStruct.checkMustHaves(fleet.liveHosts(sshList)))
		run(baselineStruct.checkMustNotHaves(fleet.liveHosts(sshList)))
	}
	return report
}
