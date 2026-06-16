package sshlib

import (
	"testing"

	"gopkg.in/yaml.v2"
)

// mapSlice builds a yaml.MapSlice from alternating key/value pairs, mirroring how
// a single server entry is decoded from pool.yml.
func mapSlice(pairs ...interface{}) yaml.MapSlice {
	if len(pairs)%2 != 0 {
		panic("mapSlice requires an even number of arguments")
	}
	ms := make(yaml.MapSlice, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		ms = append(ms, yaml.MapItem{Key: pairs[i], Value: pairs[i+1]})
	}
	return ms
}

func TestParseServerCaseInsensitiveKeys(t *testing.T) {
	cases := []struct {
		name    string
		keyName string
	}{
		{"underscore", "Key_Path"},
		{"lower", "keypath"},
		{"camel", "KeyPath"},
		{"upper", "KEYPATH"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pp := parseServer(mapSlice(tc.keyName, "/tmp/id_rsa"))
			if pp.KeyPath != "/tmp/id_rsa" {
				t.Fatalf("key %q: got KeyPath %q, want %q", tc.keyName, pp.KeyPath, "/tmp/id_rsa")
			}
		})
	}
}

func TestParseServerAllFieldsOrderIndependent(t *testing.T) {
	// Reordered, mixed-case keys; FQDN/OS/Username/Password populated.
	sv := mapSlice(
		"Port", 2222,
		"os", "ubuntu",
		"FQDN", "host.example.com",
		"Password", "secret",
		"User_Name", "admin", // not a recognized key -> ignored
		"Username", "deploy",
	)
	pp := parseServer(sv)
	if pp.FQDN != "host.example.com" {
		t.Errorf("FQDN: got %q", pp.FQDN)
	}
	if pp.OS != "ubuntu" {
		t.Errorf("OS: got %q", pp.OS)
	}
	if pp.Username != "deploy" {
		t.Errorf("Username: got %q", pp.Username)
	}
	if pp.Password != "secret" {
		t.Errorf("Password: got %q", pp.Password)
	}
	if pp.Port != 2222 {
		t.Errorf("Port: got %d, want 2222", pp.Port)
	}
}

func TestParseServerPortFromInt(t *testing.T) {
	pp := parseServer(mapSlice("FQDN", "h", "Port", 2299))
	if pp.Port != 2299 {
		t.Fatalf("Port from int: got %d, want 2299", pp.Port)
	}
}

func TestParseServerPortFromString(t *testing.T) {
	pp := parseServer(mapSlice("FQDN", "h", "Port", "2300"))
	if pp.Port != 2300 {
		t.Fatalf("Port from string: got %d, want 2300", pp.Port)
	}
}

func TestParseServerPortAbsentOrUnparseable(t *testing.T) {
	// Absent port leaves 0 for the defaulter.
	if pp := parseServer(mapSlice("FQDN", "h")); pp.Port != 0 {
		t.Errorf("absent Port: got %d, want 0", pp.Port)
	}
	// Empty/nil value (as produced by `Port:` with no value) leaves 0.
	if pp := parseServer(mapSlice("FQDN", "h", "Port", nil)); pp.Port != 0 {
		t.Errorf("nil Port: got %d, want 0", pp.Port)
	}
	// Non-numeric string leaves 0.
	if pp := parseServer(mapSlice("FQDN", "h", "Port", "abc")); pp.Port != 0 {
		t.Errorf("bad Port: got %d, want 0", pp.Port)
	}
}

func TestDefaulterDefaultsUsernameAndPort(t *testing.T) {
	pp := parseServer(mapSlice("FQDN", "h", "Password", "secret"))
	pp.defaulter()
	if pp.Username != "root" {
		t.Errorf("Username default: got %q, want root", pp.Username)
	}
	if pp.Port != 22 {
		t.Errorf("Port default: got %d, want 22", pp.Port)
	}
}

func TestDefaulterKeepsExplicitValues(t *testing.T) {
	pp := parseServer(mapSlice("FQDN", "h", "Username", "deploy", "KeyPath", "/k", "Port", 2222))
	pp.defaulter()
	if pp.Username != "deploy" {
		t.Errorf("Username: got %q, want deploy", pp.Username)
	}
	if pp.Port != 2222 {
		t.Errorf("Port: got %d, want 2222", pp.Port)
	}
	if pp.KeyPath != "/k" {
		t.Errorf("KeyPath: got %q, want /k", pp.KeyPath)
	}
}

func TestDefaulterKeyPathOnlyDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic with KeyPath set: %v", r)
		}
	}()
	pp := parseServer(mapSlice("FQDN", "h", "Key_Path", "/k"))
	pp.defaulter()
}

func TestDefaulterPasswordOnlyDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic with Password set: %v", r)
		}
	}()
	// Mirrors the integration pool: Password present, empty Key_Path.
	pp := parseServer(mapSlice("FQDN", "h", "Password", "root", "Key_Path", nil, "Port", 22))
	pp.defaulter()
	if pp.Port != 22 {
		t.Errorf("Port: got %d, want 22", pp.Port)
	}
}

func TestDefaulterNoAuthPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when both Password and Key_Path are empty")
		}
	}()
	pp := parseServer(mapSlice("FQDN", "h"))
	pp.defaulter()
}
