package sshlib

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// fakeExecutor is an in-memory baselineExecutor for driving the orchestrator
// without real SSH. It records every command it runs and signals on close.
type fakeExecutor struct {
	mu     sync.Mutex
	got    []string
	runFn  func(cmd string) (string, []byte, error)
	closed chan struct{}
}

func newFake(runFn func(string) (string, []byte, error)) *fakeExecutor {
	return &fakeExecutor{runFn: runFn, closed: make(chan struct{})}
}

func (f *fakeExecutor) Run(cmd string) (string, []byte, error) {
	f.mu.Lock()
	f.got = append(f.got, cmd)
	f.mu.Unlock()
	if f.runFn != nil {
		return f.runFn(cmd)
	}
	return statusOK, nil, nil
}

func (f *fakeExecutor) Close() error {
	close(f.closed)
	return nil
}

func (f *fakeExecutor) commands() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]string, len(f.got))
	copy(out, f.got)
	return out
}

// TestDispatchFanOut: one step reaches every live worker exactly once.
func TestDispatchFanOut(t *testing.T) {
	fakes := map[string]*fakeExecutor{"a": newFake(nil), "b": newFake(nil), "c": newFake(nil)}
	execs := map[string]baselineExecutor{}
	for h, f := range fakes {
		execs[h] = f
	}
	fleet := newFleetWithExecutors(execs)

	res := fleet.dispatch(map[string]string{"a": "cmd", "b": "cmd", "c": "cmd"}, time.Second)
	fleet.shutdown()

	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}
	for h, r := range res {
		if r.Status != statusOK || r.Err != nil {
			t.Errorf("host %s: expected OK, got %q err=%v", h, r.Status, r.Err)
		}
	}
	for h, f := range fakes {
		if got := f.commands(); len(got) != 1 || got[0] != "cmd" {
			t.Errorf("host %s ran %v, want exactly [cmd]", h, got)
		}
	}
}

// TestDispatchOrdering: sequential steps arrive at each worker in order, and
// step N+1 only runs after step N's results are gathered.
func TestDispatchOrdering(t *testing.T) {
	fakes := map[string]*fakeExecutor{"a": newFake(nil), "b": newFake(nil)}
	execs := map[string]baselineExecutor{}
	for h, f := range fakes {
		execs[h] = f
	}
	fleet := newFleetWithExecutors(execs)

	for _, step := range []string{"1", "2", "3"} {
		fleet.dispatch(map[string]string{"a": step, "b": step}, time.Second)
	}
	fleet.shutdown()

	for h, f := range fakes {
		got := f.commands()
		want := []string{"1", "2", "3"}
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("host %s order %v, want %v", h, got, want)
		}
	}
}

// TestDispatchDeadHostIsolation: a hung host is declared dead after the timeout,
// the rest of the fleet still completes, and the dead host is skipped afterward.
func TestDispatchDeadHostIsolation(t *testing.T) {
	release := make(chan struct{})
	defer close(release)
	slow := newFake(func(string) (string, []byte, error) {
		<-release // block past the dispatch timeout
		return statusOK, nil, nil
	})
	fast := newFake(nil)
	fleet := newFleetWithExecutors(map[string]baselineExecutor{"slow": slow, "fast": fast})

	res := fleet.dispatch(map[string]string{"slow": "x", "fast": "x"}, 100*time.Millisecond)
	if r, ok := res["fast"]; !ok || r.Status != statusOK {
		t.Errorf("fast host should have completed OK, got %+v", res["fast"])
	}
	if !fleet.dead["slow"] {
		t.Errorf("slow host should be marked dead after timeout")
	}

	// second step must skip the dead host without blocking
	res2 := fleet.dispatch(map[string]string{"slow": "y", "fast": "y"}, 100*time.Millisecond)
	if _, ok := res2["slow"]; ok {
		t.Errorf("dead host should be skipped on subsequent steps")
	}
	if r, ok := res2["fast"]; !ok || r.Status != statusOK {
		t.Errorf("fast host should still run, got %+v", res2["fast"])
	}
	fleet.shutdown()
}

// TestShutdownClosesWorkers: shutdown drains and closes every worker.
func TestShutdownClosesWorkers(t *testing.T) {
	fakes := map[string]*fakeExecutor{"a": newFake(nil), "b": newFake(nil)}
	execs := map[string]baselineExecutor{}
	for h, f := range fakes {
		execs[h] = f
	}
	fleet := newFleetWithExecutors(execs)
	fleet.dispatch(map[string]string{"a": "x", "b": "x"}, time.Second)
	fleet.shutdown()

	for h, f := range fakes {
		select {
		case <-f.closed:
		case <-time.After(time.Second):
			t.Errorf("worker %s was not closed after shutdown", h)
		}
	}
}

// TestNoOpCommandSkipsExecutor: an empty command yields OK without invoking Run.
func TestNoOpCommandSkipsExecutor(t *testing.T) {
	f := newFake(func(string) (string, []byte, error) {
		t.Fatalf("Run should not be called for an empty command")
		return statusNOK, nil, nil
	})
	fleet := newFleetWithExecutors(map[string]baselineExecutor{"a": f})
	res := fleet.dispatch(map[string]string{"a": ""}, time.Second)
	fleet.shutdown()
	if res["a"].Status != statusOK {
		t.Errorf("no-op should report OK, got %q", res["a"].Status)
	}
}

// TestNokKeepsHostAlive: a non-compliant result (NOK, no error) must NOT drop
// the host — later steps still reach it. Only a transport error marks it dead.
// This guards the check path, where non-compliant probes are normal.
func TestNokKeepsHostAlive(t *testing.T) {
	calls := 0
	f := newFake(func(string) (string, []byte, error) {
		calls++
		return statusNOK, nil, nil // ran, non-compliant, connection healthy
	})
	fleet := newFleetWithExecutors(map[string]baselineExecutor{"a": f})

	r1 := fleet.dispatch(map[string]string{"a": "probe1"}, time.Second)
	if r1["a"].Status != statusNOK {
		t.Fatalf("step1 want NOK, got %q", r1["a"].Status)
	}
	if fleet.dead["a"] {
		t.Fatalf("a NOK result must not mark the host dead")
	}
	r2 := fleet.dispatch(map[string]string{"a": "probe2"}, time.Second)
	if _, ok := r2["a"]; !ok {
		t.Fatalf("host should still receive step2 after a NOK")
	}
	fleet.shutdown()
	if calls != 2 {
		t.Fatalf("expected 2 probes to run, got %d", calls)
	}
}

// TestTransportErrorMarksDead: an executor error (transport/session failure)
// drops the host from subsequent steps.
func TestTransportErrorMarksDead(t *testing.T) {
	f := newFake(func(string) (string, []byte, error) {
		return statusNOK, nil, fmt.Errorf("session died")
	})
	fleet := newFleetWithExecutors(map[string]baselineExecutor{"a": f})
	fleet.dispatch(map[string]string{"a": "x"}, time.Second)
	if !fleet.dead["a"] {
		t.Fatalf("transport error should mark the host dead")
	}
	r2 := fleet.dispatch(map[string]string{"a": "y"}, time.Second)
	if _, ok := r2["a"]; ok {
		t.Fatalf("dead host must be skipped on later steps")
	}
	fleet.shutdown()
}

// TestDispatchPreservesOutput: the worker carries the executor's captured output
// through to the Result so collect steps can persist it.
func TestDispatchPreservesOutput(t *testing.T) {
	f := newFake(func(string) (string, []byte, error) {
		return statusOK, []byte("hello world\n"), nil
	})
	fleet := newFleetWithExecutors(map[string]baselineExecutor{"a": f})
	res := fleet.dispatch(map[string]string{"a": "x"}, time.Second)
	fleet.shutdown()
	if string(res["a"].Output) != "hello world\n" {
		t.Fatalf("output not preserved, got %q", res["a"].Output)
	}
}

// TestCollectStepWritesArtifacts: collectStep writes each alive host's output to
// ./collections/<host>/<collectAs>, sanitizing the artifact name to a base file.
func TestCollectStepWritesArtifacts(t *testing.T) {
	dir := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	fleet := &baselineFleet{dead: map[string]bool{"dead": true}}
	results := map[string]Result{
		"h1":   {Host: "h1", Status: statusOK, Output: []byte("loadavg-data")},
		"dead": {Host: "dead", Status: statusNOK, Output: []byte("should-not-write")},
	}
	collectStep("../escape.stat", results, fleet)

	got, err := os.ReadFile(filepath.Join(dir, "collections", "h1", "escape.stat"))
	if err != nil {
		t.Fatalf("expected artifact for h1: %v", err)
	}
	if string(got) != "loadavg-data" {
		t.Fatalf("artifact content = %q", got)
	}
	if _, err := os.Stat(filepath.Join(dir, "collections", "dead")); !os.IsNotExist(err) {
		t.Fatalf("dead host should not get a collections dir")
	}
}

// TestApplyMustHavesBuilder: the builder emits one step per item, keyed by host,
// with sudo applied per-host root-ness (the OS-key + per-host-sudo fixes).
func TestApplyMustHavesBuilder(t *testing.T) {
	var b ParsedBaseline
	b.musthave.installed = []string{"git"}
	sshList := map[string]string{"h1": "ubuntu", "h2": "centos"}
	isRoot := map[string]bool{"h1": true, "h2": false}

	var reboot bool
	steps := b.applyMustHaves(sshList, isRoot, &reboot)
	if len(steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(steps))
	}
	st := steps[0]
	if len(st.cmds) != 2 {
		t.Fatalf("expected commands for 2 hosts, got %d", len(st.cmds))
	}
	if c := st.cmds["h2"]; len(c) < 5 || c[:5] != "sudo " {
		t.Errorf("non-root host should be sudo-prefixed, got %q", c)
	}
	if c := st.cmds["h1"]; len(c) >= 5 && c[:5] == "sudo " {
		t.Errorf("root host should not be sudo-prefixed, got %q", c)
	}
}
