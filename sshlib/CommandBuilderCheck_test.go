package sshlib

import (
	"strings"
	"testing"
)

func TestPolicyCommandBuilderCheckBasename(t *testing.T) {
	tests := []struct {
		name      string
		polimport string
		wantBase  string
	}{
		{"plain", "mymod.pp", "mymod.pp"},
		{"absolute path", "/etc/selinux/targeted/mymod.pp", "mymod.pp"},
		{"relative path", "policies/sub/other.pp", "other.pp"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("policyCommandBuilder panicked: %v", r)
				}
			}()
			p := &musthavepolicies{polimport: tt.polimport}
			got := p.policyCommandBuilder("check")
			if !strings.Contains(got, tt.wantBase) {
				t.Errorf("check command %q does not contain basename %q", got, tt.wantBase)
			}
			if !strings.Contains(got, "semanage") {
				t.Errorf("check command %q does not use semanage", got)
			}
		})
	}
}

func TestPolicyCommandBuilderCheckEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("policyCommandBuilder panicked on empty polimport: %v", r)
		}
	}()
	p := &musthavepolicies{polimport: ""}
	if got := p.policyCommandBuilder("check"); got != "" {
		t.Errorf("empty polimport should yield empty command, got %q", got)
	}
}

func TestFirewallCommandBuilderCheck(t *testing.T) {
	tests := []struct {
		name      string
		port      string
		protocol  string
		wantProto []string
	}{
		{"tcp", "8080", "tcp", []string{"tcp"}},
		{"tcp udp", "8080", "tcp udp", []string{"tcp", "udp"}},
		{"single udp", "53", "udp", []string{"udp"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			port := tt.port
			proto := tt.protocol
			zone := ""
			got := firewallCommandBuilder(&port, &proto, &zone, "check")

			if got == "" {
				t.Fatal("check command is empty")
			}
			if !strings.Contains(got, tt.port) {
				t.Errorf("check command %q does not contain port %q", got, tt.port)
			}
			for _, p := range tt.wantProto {
				if !strings.Contains(got, p) {
					t.Errorf("check command %q does not contain protocol %q", got, p)
				}
			}
			if strings.Contains(got, " > awk") {
				t.Errorf("check command %q still contains the ` > awk` redirect", got)
			}
			if strings.Contains(got, "> awk") {
				t.Errorf("check command %q still contains an awk redirect", got)
			}
			if !strings.Contains(got, "grep") {
				t.Errorf("check command %q does not use grep", got)
			}
			// Multi-protocol checks must be chained so all must match.
			if len(tt.wantProto) > 1 && !strings.Contains(got, " && ") {
				t.Errorf("multi-protocol check %q is not chained with &&", got)
			}
		})
	}
}
