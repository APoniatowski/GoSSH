package sshlib

import "testing"

func TestServiceCommandBuilderNonInteractive(t *testing.T) {
	const pkg = "nginx"

	tests := []struct {
		os            string
		wantInstall   string
		wantUninstall string
	}{
		{"debian", "apt-get install -y nginx", "apt-get remove -y nginx"},
		{"ubuntu", "apt-get install -y nginx", "apt-get remove -y nginx"},
		{"centos", "yum install -y nginx", "yum remove -y nginx"},
		{"rhel", "yum install -y nginx", "yum remove -y nginx"},
		{"fedora", "dnf install -y nginx", "dnf remove -y nginx"},
		{"opensuse", "zypper --non-interactive install nginx", "zypper --non-interactive remove nginx"},
		{"sles", "zypper --non-interactive install nginx", "zypper --non-interactive remove nginx"},
		{"arch", "pacman -S --noconfirm nginx", "pacman -R --noconfirm nginx"},
		{"freebsd", "pkg install -y nginx", "pkg delete -y nginx"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.os, func(t *testing.T) {
			service := pkg
			os := tt.os

			if got := serviceCommandBuilder(&service, &os, "install"); got != tt.wantInstall {
				t.Errorf("install[%s] = %q, want %q", tt.os, got, tt.wantInstall)
			}
			if got := serviceCommandBuilder(&service, &os, "uninstall"); got != tt.wantUninstall {
				t.Errorf("uninstall[%s] = %q, want %q", tt.os, got, tt.wantUninstall)
			}
		})
	}
}
