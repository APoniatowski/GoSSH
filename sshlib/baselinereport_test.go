package sshlib

import (
	"fmt"
	"testing"
)

// TestComplianceReportCompliantAndExit exercises Compliant()/ExitCode() across
// the all-OK, single-NOK, single-Err, and Err-outranks-NOK cases.
func TestComplianceReportCompliantAndExit(t *testing.T) {
	tests := []struct {
		name          string
		results       []StepResult
		wantCompliant bool
		wantExit      int
	}{
		{
			name: "all OK",
			results: []StepResult{
				{Host: "a", Status: statusOK, Compliant: true},
				{Host: "b", Status: statusOK, Compliant: true},
			},
			wantCompliant: true,
			wantExit:      0,
		},
		{
			name:          "empty is compliant",
			results:       nil,
			wantCompliant: true,
			wantExit:      0,
		},
		{
			name: "one NOK",
			results: []StepResult{
				{Host: "a", Status: statusOK, Compliant: true},
				{Host: "b", Status: statusNOK, Compliant: false},
			},
			wantCompliant: false,
			wantExit:      2,
		},
		{
			name: "one Err",
			results: []StepResult{
				{Host: "a", Status: statusOK, Compliant: true},
				{Host: "b", Status: statusNOK, Err: fmt.Errorf("session died")},
			},
			wantCompliant: false,
			wantExit:      1,
		},
		{
			name: "Err outranks concurrent NOK",
			results: []StepResult{
				{Host: "a", Status: statusNOK, Compliant: false},
				{Host: "b", Status: statusNOK, Err: fmt.Errorf("timed out")},
			},
			wantCompliant: false,
			wantExit:      1,
		},
		{
			name: "inverted absent (NOK) is compliant",
			results: []StepResult{
				{Host: "a", Status: statusNOK, Compliant: true},
			},
			wantCompliant: true,
			wantExit:      0,
		},
		{
			name: "inverted present (OK) is non-compliant",
			results: []StepResult{
				{Host: "a", Status: statusOK, Compliant: false},
			},
			wantCompliant: false,
			wantExit:      2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := &ComplianceReport{Results: tc.results}
			if got := r.Compliant(); got != tc.wantCompliant {
				t.Errorf("Compliant() = %v, want %v", got, tc.wantCompliant)
			}
			if got := r.ExitCode(); got != tc.wantExit {
				t.Errorf("ExitCode() = %d, want %d", got, tc.wantExit)
			}
		})
	}
}

// TestComplianceReportAdd verifies Add appends one StepResult per host and that
// the recorded fields match the source Result.
func TestComplianceReportAdd(t *testing.T) {
	r := &ComplianceReport{}
	results := map[string]Result{
		"a": {Host: "a", Status: statusOK},
		"b": {Host: "b", Status: statusNOK},
	}
	r.Add("group1", "install git", results, false)
	if len(r.Results) != 2 {
		t.Fatalf("Add appended %d results, want 2", len(r.Results))
	}
	for _, sr := range r.Results {
		if sr.Group != "group1" || sr.Label != "install git" {
			t.Errorf("group/label not propagated: %+v", sr)
		}
		if sr.Host != "a" && sr.Host != "b" {
			t.Errorf("unexpected host %q", sr.Host)
		}
		// Non-inverted: OK -> compliant, NOK -> non-compliant.
		switch sr.Host {
		case "a":
			if !sr.Compliant {
				t.Errorf("host a (OK, non-inverted) should be compliant: %+v", sr)
			}
		case "b":
			if sr.Compliant {
				t.Errorf("host b (NOK, non-inverted) should not be compliant: %+v", sr)
			}
		}
	}
}

// TestComplianceReportAddInvert verifies inverted must-not-have semantics: a raw
// statusOK (unwanted item PRESENT) is non-compliant, while statusNOK (item
// ABSENT) is compliant. A transport Err is never compliant regardless of invert.
func TestComplianceReportAddInvert(t *testing.T) {
	r := &ComplianceReport{}
	results := map[string]Result{
		"present": {Host: "present", Status: statusOK},                           // unwanted item present
		"absent":  {Host: "absent", Status: statusNOK},                           // unwanted item absent
		"dead":    {Host: "dead", Status: statusNOK, Err: fmt.Errorf("no conn")}, // unreachable
	}
	r.Add("group1", "Installed? telnet", results, true)
	for _, sr := range r.Results {
		switch sr.Host {
		case "present":
			if sr.Compliant {
				t.Errorf("inverted present host should be non-compliant: %+v", sr)
			}
		case "absent":
			if !sr.Compliant {
				t.Errorf("inverted absent host should be compliant: %+v", sr)
			}
		case "dead":
			if sr.Compliant {
				t.Errorf("errored host must never be compliant: %+v", sr)
			}
		}
	}
	if got := r.ExitCode(); got != 1 {
		t.Errorf("ExitCode() = %d, want 1 (Err outranks non-compliance)", got)
	}
}

// TestComplianceReportSummary checks the "<ok>/<total> compliant" count, where an
// errored host counts as not-OK.
func TestComplianceReportSummary(t *testing.T) {
	r := &ComplianceReport{Results: []StepResult{
		{Host: "a", Status: statusOK, Compliant: true},
		{Host: "b", Status: statusOK, Compliant: true},
		{Host: "c", Status: statusNOK, Compliant: false},
		{Host: "d", Status: statusNOK, Err: fmt.Errorf("dead")},
	}}
	if got, want := r.Summary(), "2/4 compliant"; got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}
