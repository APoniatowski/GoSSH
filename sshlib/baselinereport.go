package sshlib

import "fmt"

// StepResult is the structured outcome of one dispatched step for a single host.
// It mirrors the one-line OK/NOK that reportStep prints, but retains it so the
// run can compute compliance and a real process exit code.
type StepResult struct {
	Group     string // server group name this step ran against
	Label     string // human label of the step (same as reportStep's label)
	Host      string // host fqdn
	Status    string // statusOK / statusNOK (raw executor status)
	Err       error  // transport/dead-host error, if any
	Compliant bool   // computed compliance (interprets inverted steps)
}

// ComplianceReport accumulates every StepResult across all groups and stages of
// a baseline run. It is the single source of truth for the process exit code.
type ComplianceReport struct {
	Results []StepResult
}

// Add appends one StepResult per host for a dispatched step. The results map is
// host -> Result as produced by baselineFleet.dispatch.
//
// invert flips the compliance interpretation for must-not-have probes whose raw
// statusOK means the UNWANTED item is PRESENT (a violation). For a normal step
// compliance means statusOK; for an inverted step compliance means statusNOK
// (the unwanted item is ABSENT). A transport Err is never compliant.
func (r *ComplianceReport) Add(group, label string, results map[string]Result, invert bool) {
	for _, res := range results {
		compliant := res.Err == nil &&
			((!invert && res.Status == statusOK) ||
				(invert && res.Status == statusNOK))
		r.Results = append(r.Results, StepResult{
			Group:     group,
			Label:     label,
			Host:      res.Host,
			Status:    res.Status,
			Err:       res.Err,
			Compliant: compliant,
		})
	}
}

// Compliant reports whether every recorded step is compliant.
func (r *ComplianceReport) Compliant() bool {
	for _, sr := range r.Results {
		if !sr.Compliant {
			return false
		}
	}
	return true
}

// ExitCode maps the report to a process exit code.
//
// Precedence (chosen): a transport error (1) outranks a plain non-compliance
// (2). Rationale: a transport/dead-host error means we could not actually
// determine compliance for that host, which is a more severe, operator-actionable
// failure than a host that answered and was merely non-compliant. So:
//
//	1 -> any StepResult has a non-nil Err (transport/dead host)
//	2 -> otherwise, any StepResult is non-compliant
//	0 -> fully compliant
func (r *ComplianceReport) ExitCode() int {
	var sawNonCompliant bool
	for _, sr := range r.Results {
		if sr.Err != nil {
			return 1
		}
		if !sr.Compliant {
			sawNonCompliant = true
		}
	}
	if sawNonCompliant {
		return 2
	}
	return 0
}

// Summary returns a short "<ok>/<total> compliant" line. It does NOT print
// per-step detail; reportStep remains responsible for that.
func (r *ComplianceReport) Summary() string {
	total := len(r.Results)
	ok := 0
	for _, sr := range r.Results {
		if sr.Compliant {
			ok++
		}
	}
	return fmt.Sprintf("%d/%d compliant", ok, total)
}
