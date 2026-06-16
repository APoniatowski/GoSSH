package sshlib

import (
	"strings"
	"testing"
)

// TestPrereqURLFetchSaves verifies the plain-URL download path actually persists
// the file (it previously streamed to stdout / built a malformed symlink).
func TestPrereqURLFetchSaves(t *testing.T) {
	url := "http://fileserver/payload.txt"
	no, yes := false, true

	// No cleanup: must save into the home dir by basename (curl needs -O).
	got := prereqURLFetch(&url, &no)
	if !strings.Contains(got, "curl "+url+" -O") {
		t.Errorf("no-cleanup fetch must save with curl -O, got: %q", got)
	}
	if !strings.Contains(got, "wget "+url) {
		t.Errorf("no-cleanup fetch must fall back to wget, got: %q", got)
	}

	// Cleanup: download to /tmp and symlink into home with a SPACE before ~/.
	gotC := prereqURLFetch(&url, &yes)
	if !strings.Contains(gotC, "-o /tmp/payload.txt") {
		t.Errorf("cleanup fetch must download to /tmp, got: %q", gotC)
	}
	if !strings.Contains(gotC, "ln -sfn /tmp/payload.txt ~/payload.txt") {
		t.Errorf("cleanup fetch must build a valid symlink, got: %q", gotC)
	}
	if strings.Contains(gotC, "payload.txt~/") {
		t.Errorf("malformed (space-less) symlink must not appear, got: %q", gotC)
	}
}
