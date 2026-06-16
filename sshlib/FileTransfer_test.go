package sshlib

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFileToCommand verifies the base64-over-command-channel push: the encoded
// content is present, tee is used (with sudo only when requested), the dest is
// referenced, and no legacy [FILETRANSFER] marker leaks into the command.
func TestFileToCommand(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "hello.txt")
	content := []byte("hello gossh\n")
	if err := os.WriteFile(src, content, 0o644); err != nil {
		t.Fatalf("write temp src: %v", err)
	}
	wantB64 := base64.StdEncoding.EncodeToString(content)

	tests := []struct {
		name string
		sudo bool
	}{
		{"non-sudo", false},
		{"sudo", true},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := fileToCommand(src, "/etc/gossh/hello.txt", tc.sudo)
			if err != nil {
				t.Fatalf("fileToCommand: %v", err)
			}
			if !strings.Contains(cmd, wantB64) {
				t.Errorf("command missing base64 of content: %q", cmd)
			}
			if !strings.Contains(cmd, "tee '/etc/gossh/hello.txt'") {
				t.Errorf("command must tee to dest, got: %q", cmd)
			}
			if !strings.Contains(cmd, "mkdir -p '/etc/gossh'") {
				t.Errorf("command must mkdir dest dir, got: %q", cmd)
			}
			if strings.Contains(cmd, "[FILETRANSFER]") {
				t.Errorf("legacy marker must not appear, got: %q", cmd)
			}
			hasSudo := strings.Contains(cmd, "sudo ")
			if tc.sudo && !hasSudo {
				t.Errorf("sudo=true must prefix sudo, got: %q", cmd)
			}
			if !tc.sudo && hasSudo {
				t.Errorf("sudo=false must not contain sudo, got: %q", cmd)
			}
		})
	}
}

// TestFileToCommandMissingSource returns an error for an unreadable source.
func TestFileToCommandMissingSource(t *testing.T) {
	if _, err := fileToCommand(filepath.Join(t.TempDir(), "nope"), "/tmp/x", false); err == nil {
		t.Fatalf("expected error for missing source")
	}
}
