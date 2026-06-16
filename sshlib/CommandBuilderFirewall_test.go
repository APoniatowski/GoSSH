package sshlib

import (
	"strings"
	"testing"
)

// TestFirewallRemoveRules verifies the remove-* cases DELETE the matching rule
// rather than adding its inverse.
func TestFirewallRemoveRules(t *testing.T) {
	port := "8080"
	proto := "tcp"
	zone := ""

	t.Run("remove-open deletes allow rule", func(t *testing.T) {
		cmd := firewallCommandBuilder(&port, &proto, &zone, "remove-open")
		for _, want := range []string{
			"--remove-port=8080/tcp",
			"ufw delete allow 8080/tcp",
			"-D INPUT",
			"ACCEPT",
		} {
			if !strings.Contains(cmd, want) {
				t.Errorf("remove-open missing %q in %q", want, cmd)
			}
		}
		if strings.Contains(cmd, "ufw deny") || strings.Contains(cmd, "-j DROP") || strings.Contains(cmd, "-A INPUT") {
			t.Errorf("remove-open must not add/deny: %q", cmd)
		}
	})

	t.Run("remove-closed deletes deny rule", func(t *testing.T) {
		cmd := firewallCommandBuilder(&port, &proto, &zone, "remove-closed")
		for _, want := range []string{
			"--remove-port=8080/tcp",
			"ufw delete deny 8080/tcp",
			"-D INPUT",
			"DROP",
		} {
			if !strings.Contains(cmd, want) {
				t.Errorf("remove-closed missing %q in %q", want, cmd)
			}
		}
		if strings.Contains(cmd, "ufw allow") || strings.Contains(cmd, "-j ACCEPT") || strings.Contains(cmd, "-A INPUT") {
			t.Errorf("remove-closed must not add/allow: %q", cmd)
		}
	})
}
