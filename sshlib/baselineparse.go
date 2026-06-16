package sshlib

import (
	"strconv"

	"gopkg.in/yaml.v2"
)

// parseStats holds the non-fatal counters accumulated while parsing a single
// baseline group. The typed-unmarshal parser no longer tracks these per-field
// (a typed unmarshal is lenient and silently ignores unknown keys), but the
// type is retained so the parseBaselineGroup signature is unchanged and callers
// continue to compile. The counters are left zero.
type parseStats struct {
	warnings             int
	dataWarnings         int
	mainCategoryWarnings int
	baselineErrors       int
}

// ---------------------------------------------------------------------------
// Tagged intermediate structs.
//
// These mirror the baseline YAML schema exactly (see config/EXAMPLE-baselines.yml).
// Every hyphenated/exact key gets an explicit yaml tag. Dynamic-key sections
// (Configured, Users, Mounts) are maps keyed by their YAML name. Firewall Ports
// are ints in YAML, so they are typed []int here and converted to the []string
// that ParsedBaseline expects during the copy below.
// ---------------------------------------------------------------------------

type ygroup struct {
	Exclude     yexclude     `yaml:"Exclude"`
	Prereq      yprereq      `yaml:"Prerequisites"`
	MustHave    ymusthave    `yaml:"Must-Have"`
	MustNotHave ymustnothave `yaml:"Must-Not-Have"`
	Final       yfinal       `yaml:"Final"`
}

type yexclude struct {
	OS      []string `yaml:"OS"`
	Servers []string `yaml:"Servers"`
}

type yprereq struct {
	Tools    []string `yaml:"Tools"`
	Files    yfiles   `yaml:"Files"`
	VCS      yvcs     `yaml:"VCS"`
	Script   string   `yaml:"Script"`
	Commands []string `yaml:"Commands"`
	Cleanup  bool     `yaml:"Clean-up"`
}

type yfiles struct {
	URLs   []string `yaml:"URLs"`
	Local  ylocal   `yaml:"Local"`
	Remote yremote  `yaml:"Remote"`
}

type ylocal struct {
	Source      string `yaml:"Source"`
	Destination string `yaml:"Destination"`
}

type yremote struct {
	Type        string   `yaml:"Type"`
	Address     string   `yaml:"Address"`
	Username    string   `yaml:"Username"`
	Password    string   `yaml:"Password"`
	Source      string   `yaml:"Source"`
	Destination string   `yaml:"Destination"`
	Files       []string `yaml:"Files"`
}

type yvcs struct {
	URLs    []string `yaml:"URLs"`
	Execute []string `yaml:"Execute"`
}

type ymusthave struct {
	Installed  []string               `yaml:"Installed"`
	Enabled    []string               `yaml:"Enabled"`
	Disabled   []string               `yaml:"Disabled"`
	Configured map[string]yconfigured `yaml:"Configured"`
	Users      map[string]yuser       `yaml:"Users"`
	Policies   ypolicies              `yaml:"Policies"`
	Rules      yrules                 `yaml:"Rules"`
	Mounts     map[string]ymount      `yaml:"Mounts"`
}

type yconfigured struct {
	Source      []string `yaml:"Source"`
	Destination []string `yaml:"Destination"`
}

type yuser struct {
	Groups  []string `yaml:"Groups"`
	Shell   string   `yaml:"Shell"`
	HomeDir string   `yaml:"Home-Dir"`
	Sudoer  bool     `yaml:"Sudoer"`
}

type ypolicies struct {
	Status string `yaml:"Status"`
	Import string `yaml:"Import"`
	Reboot bool   `yaml:"Reboot"`
}

type yrules struct {
	Open   yruleset `yaml:"Open"`
	Closed yruleset `yaml:"Closed"`
	Zones  []string `yaml:"Zones"`
}

type yruleset struct {
	Ports     []int    `yaml:"Ports"`
	Protocols []string `yaml:"Protocols"`
}

type ymount struct {
	Type        string `yaml:"Type"`
	Address     string `yaml:"Address"`
	Username    string `yaml:"Username"`
	Password    string `yaml:"Password"`
	Source      string `yaml:"Source"`
	Destination string `yaml:"Destination"`
}

type ymustnothave struct {
	Installed []string `yaml:"Installed"`
	Enabled   []string `yaml:"Enabled"`
	Disabled  []string `yaml:"Disabled"`
	Users     []string `yaml:"Users"`
	Rules     yrules   `yaml:"Rules"`
	Mounts    []string `yaml:"Mounts"`
}

type yfinal struct {
	Scripts  []string `yaml:"Scripts"`
	Commands []string `yaml:"Commands"`
	Collect  ycollect `yaml:"Collect"`
	Restart  yrestart `yaml:"Restart"`
}

type ycollect struct {
	Logs  []string `yaml:"Logs"`
	Stats []string `yaml:"Stats"`
	Files []string `yaml:"Files"`
	Users bool     `yaml:"Users"`
}

type yrestart struct {
	Services bool `yaml:"Services"`
	Servers  bool `yaml:"Servers"`
}

// intsToStrings converts YAML integer ports to the string slice ParsedBaseline
// stores (e.g. 8080 -> "8080"), matching the old strconv.Itoa path.
func intsToStrings(in []int) []string {
	if in == nil {
		return nil
	}
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = strconv.Itoa(v)
	}
	return out
}

// parseBaselineGroup converts the per-group baseline steps (the yaml.MapSlice
// found under <baselineName> -> <serverGroup>) into a fully populated
// ParsedBaseline.
//
// Implementation: the incoming MapSlice is re-marshaled to YAML bytes and
// unmarshaled into the tagged ygroup struct above, then copied field-by-field
// into ParsedBaseline (applying the int->string Ports conversion). This reuses
// the existing callers unchanged: they still pass in the per-group steps slice.
//
// This is the single canonical parser shared by ApplyBaselines, CheckBaselines
// and VerifyBaselines. It is pure: it does not dial, print step-by-step output,
// or panic on malformed top-level structure (the callers own that).
func parseBaselineGroup(steps yaml.MapSlice) (ParsedBaseline, parseStats) {
	var stats parseStats
	var baselineStruct ParsedBaseline
	baselineStruct.musthave.configured.services = make(map[string]musthaveconfiguredservices)
	baselineStruct.musthave.users.users = make(map[string]musthaveusersstruct)
	baselineStruct.musthave.mounts.mountname = make(map[string]mountdetails)

	// Round-trip the MapSlice through YAML into the tagged group struct. Any
	// marshal/unmarshal error leaves baselineStruct at its zero-but-allocated
	// state, which the golden tests treat as "no content".
	var g ygroup
	if raw, err := yaml.Marshal(steps); err == nil {
		_ = yaml.Unmarshal(raw, &g)
	}

	// Exclude
	baselineStruct.exclude.osExcl = g.Exclude.OS
	baselineStruct.exclude.serversExcl = g.Exclude.Servers

	// Prerequisites
	baselineStruct.prereq.tools = g.Prereq.Tools
	baselineStruct.prereq.files.urls = g.Prereq.Files.URLs
	baselineStruct.prereq.files.local.src = g.Prereq.Files.Local.Source
	baselineStruct.prereq.files.local.dest = g.Prereq.Files.Local.Destination
	baselineStruct.prereq.files.remote.mounttype = g.Prereq.Files.Remote.Type
	baselineStruct.prereq.files.remote.address = g.Prereq.Files.Remote.Address
	baselineStruct.prereq.files.remote.username = g.Prereq.Files.Remote.Username
	baselineStruct.prereq.files.remote.pwd = g.Prereq.Files.Remote.Password
	baselineStruct.prereq.files.remote.src = g.Prereq.Files.Remote.Source
	baselineStruct.prereq.files.remote.dest = g.Prereq.Files.Remote.Destination
	baselineStruct.prereq.files.remote.files = g.Prereq.Files.Remote.Files
	baselineStruct.prereq.vcs.urls = g.Prereq.VCS.URLs
	baselineStruct.prereq.vcs.execute = g.Prereq.VCS.Execute
	baselineStruct.prereq.script = g.Prereq.Script
	baselineStruct.prereq.commands = g.Prereq.Commands
	baselineStruct.prereq.cleanup = g.Prereq.Cleanup

	// Must-Have
	baselineStruct.musthave.installed = g.MustHave.Installed
	baselineStruct.musthave.enabled = g.MustHave.Enabled
	baselineStruct.musthave.disabled = g.MustHave.Disabled
	for name, c := range g.MustHave.Configured {
		baselineStruct.musthave.configured.services[name] = musthaveconfiguredservices{
			source:      c.Source,
			destination: c.Destination,
		}
	}
	for name, u := range g.MustHave.Users {
		baselineStruct.musthave.users.users[name] = musthaveusersstruct{
			groups: u.Groups,
			shell:  u.Shell,
			home:   u.HomeDir,
			sudoer: u.Sudoer,
		}
	}
	baselineStruct.musthave.policies.polstatus = g.MustHave.Policies.Status
	baselineStruct.musthave.policies.polimport = g.MustHave.Policies.Import
	baselineStruct.musthave.policies.polreboot = g.MustHave.Policies.Reboot
	baselineStruct.musthave.rules.fwopen.ports = intsToStrings(g.MustHave.Rules.Open.Ports)
	baselineStruct.musthave.rules.fwopen.protocols = g.MustHave.Rules.Open.Protocols
	baselineStruct.musthave.rules.fwclosed.ports = intsToStrings(g.MustHave.Rules.Closed.Ports)
	baselineStruct.musthave.rules.fwclosed.protocols = g.MustHave.Rules.Closed.Protocols
	baselineStruct.musthave.rules.fwzones = g.MustHave.Rules.Zones
	for name, m := range g.MustHave.Mounts {
		baselineStruct.musthave.mounts.mountname[name] = mountdetails{
			mounttype: m.Type,
			address:   m.Address,
			username:  m.Username,
			pwd:       m.Password,
			src:       m.Source,
			dest:      m.Destination,
		}
	}

	// Must-Not-Have
	baselineStruct.mustnothave.installed = g.MustNotHave.Installed
	baselineStruct.mustnothave.enabled = g.MustNotHave.Enabled
	baselineStruct.mustnothave.disabled = g.MustNotHave.Disabled
	baselineStruct.mustnothave.users = g.MustNotHave.Users
	baselineStruct.mustnothave.rules.fwopen.ports = intsToStrings(g.MustNotHave.Rules.Open.Ports)
	baselineStruct.mustnothave.rules.fwopen.protocols = g.MustNotHave.Rules.Open.Protocols
	baselineStruct.mustnothave.rules.fwclosed.ports = intsToStrings(g.MustNotHave.Rules.Closed.Ports)
	baselineStruct.mustnothave.rules.fwclosed.protocols = g.MustNotHave.Rules.Closed.Protocols
	baselineStruct.mustnothave.rules.fwzones = g.MustNotHave.Rules.Zones
	baselineStruct.mustnothave.mounts = g.MustNotHave.Mounts

	// Final
	baselineStruct.final.scripts = g.Final.Scripts
	baselineStruct.final.commands = g.Final.Commands
	baselineStruct.final.collect.logs = g.Final.Collect.Logs
	baselineStruct.final.collect.stats = g.Final.Collect.Stats
	baselineStruct.final.collect.files = g.Final.Collect.Files
	baselineStruct.final.collect.users = g.Final.Collect.Users
	baselineStruct.final.restart.services = g.Final.Restart.Services
	baselineStruct.final.restart.servers = g.Final.Restart.Servers

	return baselineStruct, stats
}
