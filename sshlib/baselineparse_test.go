package sshlib

// Golden / characterization tests for parseBaselineGroup.
//
// These lock the *current* observable behavior of the hand-rolled
// yaml.MapSlice -> ParsedBaseline parser so a future typed-unmarshal rewrite
// can be proven equivalent. They live in package sshlib because every
// ParsedBaseline field is unexported.
//
// Test inputs are built exactly the way production does it: a baseline file is
// unmarshaled into a yaml.MapSlice, then the per-group steps MapSlice (the
// value under <baselineName> -> <serverGroup>) is handed to parseBaselineGroup.

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

// groupSteps unmarshals a whole baseline-file snippet and returns the steps
// MapSlice for the first group of the first baseline, mirroring the navigation
// the real callers perform: doc[0].Value (groups) -> groups[0].Value (steps).
func groupSteps(t *testing.T, src string) yaml.MapSlice {
	t.Helper()
	var doc yaml.MapSlice
	require.NoError(t, yaml.Unmarshal([]byte(src), &doc), "fixture must be valid YAML")
	require.NotEmpty(t, doc, "fixture must contain a baseline")

	groups, ok := doc[0].Value.(yaml.MapSlice)
	require.True(t, ok, "baseline value must be a group MapSlice")
	require.NotEmpty(t, groups, "baseline must contain a group")

	steps, ok := groups[0].Value.(yaml.MapSlice)
	require.True(t, ok, "group value must be a steps MapSlice")
	return steps
}

// fullFixture is a complete, populated single-group baseline modeled on
// config/EXAMPLE-baselines.yml. It exercises every category and sub-section
// with meaningful (non-empty) content.
const fullFixture = `
Example baseline:
  Group 1:
    Exclude:
      OS:
        - debian
        - arch
      Servers:
        - Server 1
        - Server 3
    Prerequisites:
      Tools:
        - git
        - curl
      Files:
        URLs:
          - http://where/file/is
        Local:
          Source: /local/src
          Destination: /local/dst
        Remote:
          Type: nfs
          Address: 1.2.3.4
          Username: nfsuser
          Password: nfspassword
          Source: /remote/src
          Destination: /remote/dst
          Files:
            - script.sh
            - some.tar.gz
      VCS:
        URLs:
          - https://repo/one
          - https://repo/two
        Execute:
          - make
          - make install
      Script: /path/to/script
      Commands:
        - command 1
        - command 2
      Clean-up: true
    Must-Have:
      Installed:
        - httpd
        - firewalld
      Enabled:
        - httpd
        - rsyslog
      Disabled:
        - telnet
      Configured:
        httpd:
          Source:
            - /cfg/src1
            - /cfg/src2
          Destination:
            - /cfg/dst1
            - /cfg/dst2
        openssh:
          Source:
            - /ssh/src
          Destination:
            - /ssh/dst
      Users:
        webmaster:
          Groups:
            - www
          Shell: nologin
          Home-Dir: /home/webmaster
          Sudoer: false
        jim:
          Groups:
            - wheel
            - anothergroup
          Shell: bash
          Home-Dir: /home/jim
          Sudoer: true
      Policies:
        Status: Enforced
        Import: /path/to/policy
        Reboot: true
      Rules:
        Open:
          Ports:
            - 80
            - 443
          Protocols:
            - tcp
            - tcp udp
        Closed:
          Ports:
            - 8080
          Protocols:
            - tcp
        Zones:
          - public
      Mounts:
        Mount 1:
          Type: nfs
          Address: 1.2.3.4
          Username: nfsuser
          Password: nfspassword
          Source: /mnt/from
          Destination: /mnt/to
    Must-Not-Have:
      Installed:
        - nmap
      Enabled:
        - a-service
      Disabled:
        - httpd
      Users:
        - bob
        - jane
      Rules:
        Open:
          Ports:
            - 8080
            - 8443
          Protocols:
            - tcp udp
            - tcp
        Closed:
          Ports:
            - 80
            - 443
          Protocols:
            - tcp
            - tcp udp
        Zones:
          - public
      Mounts:
        - /path/to/mount1
        - /path/to/mount2
    Final:
      Scripts:
        - /path/to/script
        - /path/to/another
      Commands:
        - command 1
        - command 2
      Collect:
        Logs:
          - httpd
          - sshd
        Stats:
          - cpu
          - memory
        Files:
          - /collect/file
        Users: true
      Restart:
        Services: true
        Servers: false
`

func TestParseExclude(t *testing.T) {
	b, _ := parseBaselineGroup(groupSteps(t, fullFixture))
	assert.Equal(t, []string{"debian", "arch"}, b.exclude.osExcl)
	assert.Equal(t, []string{"Server 1", "Server 3"}, b.exclude.serversExcl)
}

func TestParsePrerequisites(t *testing.T) {
	b, _ := parseBaselineGroup(groupSteps(t, fullFixture))

	assert.Equal(t, []string{"git", "curl"}, b.prereq.tools)
	assert.Equal(t, []string{"http://where/file/is"}, b.prereq.files.urls)

	assert.Equal(t, "/local/src", b.prereq.files.local.src)
	assert.Equal(t, "/local/dst", b.prereq.files.local.dest)

	r := b.prereq.files.remote
	assert.Equal(t, "nfs", r.mounttype)
	assert.Equal(t, "1.2.3.4", r.address)
	assert.Equal(t, "nfsuser", r.username)
	assert.Equal(t, "nfspassword", r.pwd)
	assert.Equal(t, "/remote/src", r.src)
	assert.Equal(t, "/remote/dst", r.dest)
	assert.Equal(t, []string{"script.sh", "some.tar.gz"}, r.files)

	assert.Equal(t, []string{"https://repo/one", "https://repo/two"}, b.prereq.vcs.urls)
	assert.Equal(t, []string{"make", "make install"}, b.prereq.vcs.execute)

	assert.Equal(t, "/path/to/script", b.prereq.script)
	assert.Equal(t, []string{"command 1", "command 2"}, b.prereq.commands)
	assert.True(t, b.prereq.cleanup)
}

func TestParseMustHaveLists(t *testing.T) {
	b, _ := parseBaselineGroup(groupSteps(t, fullFixture))

	assert.Equal(t, []string{"httpd", "firewalld"}, b.musthave.installed)
	// Guard the recently-fixed copy-paste bug: Enabled -> enabled,
	// Disabled -> disabled (not swapped).
	assert.Equal(t, []string{"httpd", "rsyslog"}, b.musthave.enabled)
	assert.Equal(t, []string{"telnet"}, b.musthave.disabled)
}

func TestParseMustHaveConfigured(t *testing.T) {
	b, _ := parseBaselineGroup(groupSteps(t, fullFixture))
	svcs := b.musthave.configured.services

	require.Contains(t, svcs, "httpd")
	assert.Equal(t, []string{"/cfg/src1", "/cfg/src2"}, svcs["httpd"].source)
	assert.Equal(t, []string{"/cfg/dst1", "/cfg/dst2"}, svcs["httpd"].destination)

	require.Contains(t, svcs, "openssh")
	assert.Equal(t, []string{"/ssh/src"}, svcs["openssh"].source)
	assert.Equal(t, []string{"/ssh/dst"}, svcs["openssh"].destination)
}

func TestParseMustHaveUsers(t *testing.T) {
	b, _ := parseBaselineGroup(groupSteps(t, fullFixture))
	users := b.musthave.users.users

	require.Contains(t, users, "webmaster")
	assert.Equal(t, []string{"www"}, users["webmaster"].groups)
	assert.Equal(t, "nologin", users["webmaster"].shell)
	assert.Equal(t, "/home/webmaster", users["webmaster"].home)
	assert.False(t, users["webmaster"].sudoer)

	require.Contains(t, users, "jim")
	assert.Equal(t, []string{"wheel", "anothergroup"}, users["jim"].groups)
	assert.Equal(t, "bash", users["jim"].shell)
	assert.Equal(t, "/home/jim", users["jim"].home)
	assert.True(t, users["jim"].sudoer)
}

func TestParseMustHavePolicies(t *testing.T) {
	b, _ := parseBaselineGroup(groupSteps(t, fullFixture))
	assert.Equal(t, "Enforced", b.musthave.policies.polstatus)
	assert.Equal(t, "/path/to/policy", b.musthave.policies.polimport)
	assert.True(t, b.musthave.policies.polreboot)
}

func TestParseMustHaveRules(t *testing.T) {
	b, _ := parseBaselineGroup(groupSteps(t, fullFixture))
	rules := b.musthave.rules

	// YAML ints become string ports.
	assert.Equal(t, []string{"80", "443"}, rules.fwopen.ports)
	assert.Equal(t, []string{"tcp", "tcp udp"}, rules.fwopen.protocols)
	assert.Equal(t, []string{"8080"}, rules.fwclosed.ports)
	assert.Equal(t, []string{"tcp"}, rules.fwclosed.protocols)
	assert.Equal(t, []string{"public"}, rules.fwzones)
}

func TestParseMustHaveMounts(t *testing.T) {
	b, _ := parseBaselineGroup(groupSteps(t, fullFixture))
	mounts := b.musthave.mounts.mountname

	require.Contains(t, mounts, "Mount 1")
	m := mounts["Mount 1"]
	assert.Equal(t, "nfs", m.mounttype)
	assert.Equal(t, "1.2.3.4", m.address)
	assert.Equal(t, "nfsuser", m.username)
	assert.Equal(t, "nfspassword", m.pwd)
	assert.Equal(t, "/mnt/from", m.src)
	assert.Equal(t, "/mnt/to", m.dest)
}

func TestParseMustNotHave(t *testing.T) {
	b, _ := parseBaselineGroup(groupSteps(t, fullFixture))
	mnh := b.mustnothave

	assert.Equal(t, []string{"nmap"}, mnh.installed)
	// Guard the recently-fixed copy-paste bug on the must-not-have side too.
	assert.Equal(t, []string{"a-service"}, mnh.enabled)
	assert.Equal(t, []string{"httpd"}, mnh.disabled)
	assert.Equal(t, []string{"bob", "jane"}, mnh.users)
	assert.Equal(t, []string{"/path/to/mount1", "/path/to/mount2"}, mnh.mounts)

	assert.Equal(t, []string{"8080", "8443"}, mnh.rules.fwopen.ports)
	assert.Equal(t, []string{"tcp udp", "tcp"}, mnh.rules.fwopen.protocols)
	assert.Equal(t, []string{"80", "443"}, mnh.rules.fwclosed.ports)
	assert.Equal(t, []string{"tcp", "tcp udp"}, mnh.rules.fwclosed.protocols)
	assert.Equal(t, []string{"public"}, mnh.rules.fwzones)
}

func TestParseFinal(t *testing.T) {
	b, _ := parseBaselineGroup(groupSteps(t, fullFixture))
	f := b.final

	assert.Equal(t, []string{"/path/to/script", "/path/to/another"}, f.scripts)
	assert.Equal(t, []string{"command 1", "command 2"}, f.commands)

	assert.Equal(t, []string{"httpd", "sshd"}, f.collect.logs)
	assert.Equal(t, []string{"cpu", "memory"}, f.collect.stats)
	assert.Equal(t, []string{"/collect/file"}, f.collect.files)
	assert.True(t, f.collect.users)

	assert.True(t, f.restart.services)
	assert.False(t, f.restart.servers)
}

// TestParsePortsFromInts is a focused table check that YAML integer ports are
// stringified (the strconv.Itoa path), including the example's 8080 -> "8080".
func TestParsePortsFromInts(t *testing.T) {
	const src = `
Ports baseline:
  All:
    Must-Have:
      Rules:
        Open:
          Ports:
            - 22
            - 8080
            - 443
          Protocols:
            - tcp
`
	b, _ := parseBaselineGroup(groupSteps(t, src))
	assert.Equal(t, []string{"22", "8080", "443"}, b.musthave.rules.fwopen.ports)
}

// TestParseEmptyGroupNoPanic: a group whose sub-sections are all present but
// empty/nil must not panic and must leave content fields empty. This exercises
// the many dataWarnings nil branches.
func TestParseEmptyGroupNoPanic(t *testing.T) {
	const src = `
Stripped baseline:
  All:
    Exclude:
      OS:
      Servers:
    Prerequisites:
      Tools:
      Files:
        URLs:
        Local:
          Source:
          Destination:
        Remote:
          Type:
          Address:
          Username:
          Password:
          Source:
          Destination:
          Files:
      VCS:
        URLs:
        Execute:
      Script:
      Commands:
      Clean-up: false
    Must-Have:
      Installed:
      Enabled:
      Disabled:
      Policies:
        Status:
        Import:
        Reboot:
      Rules:
        Open:
          Ports:
          Protocols:
        Closed:
          Ports:
          Protocols:
        Zones:
    Must-Not-Have:
      Installed:
      Enabled:
      Disabled:
      Users:
      Mounts:
    Final:
      Scripts:
      Commands:
      Collect:
        Logs:
        Stats:
        Files:
        Users:
      Restart:
        Services:
        Servers:
`
	var b ParsedBaseline
	require.NotPanics(t, func() {
		b, _ = parseBaselineGroup(groupSteps(t, src))
	})

	// Nil-branch sentinels: empty-but-not-nil string slices / zero values.
	assert.Empty(t, firstNonSentinel(b.exclude.osExcl))
	assert.Empty(t, firstNonSentinel(b.prereq.tools))
	assert.Equal(t, "", b.prereq.files.local.src)
	assert.Equal(t, "", b.prereq.files.remote.address)
	assert.False(t, b.prereq.cleanup)
	assert.Empty(t, firstNonSentinel(b.musthave.installed))
	assert.Empty(t, firstNonSentinel(b.musthave.enabled))
	assert.Empty(t, firstNonSentinel(b.musthave.disabled))
	assert.Equal(t, "", b.musthave.policies.polstatus)
	assert.False(t, b.musthave.policies.polreboot)
	assert.False(t, b.final.collect.users)
	assert.False(t, b.final.restart.services)
	assert.False(t, b.final.restart.servers)
}

// firstNonSentinel returns the first meaningful (non-empty) element of a parsed
// string slice, treating absent/[]string{""} as "not set". Used to assert a
// section carries no real content without coupling to the empty-sentinel shape.
func firstNonSentinel(s []string) string {
	for _, v := range s {
		if v != "" {
			return v
		}
	}
	return ""
}
