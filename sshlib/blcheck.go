package sshlib

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"

	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"gopkg.in/yaml.v2"
)

func (blstruct *ParsedBaseline) checkOSExcludes(servergroupname string, configs *yaml.MapSlice) map[string]string {
	sshList := make(map[string]string)
	if strings.ToLower(servergroupname) == "all" {
		if len(blstruct.exclude.osExcl) == 0 &&
			len(blstruct.exclude.serversExcl) == 0 {
			var allServers yaml.MapSlice
			// Concatenates the groups to create a single group
			for _, groupItem := range *configs {
				groupValue, ok := groupItem.Value.(yaml.MapSlice)
				if !ok {
					panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
				}
				allServers = append(allServers, groupValue...)
			}
			for _, serverItem := range allServers {
				serverValue, ok := serverItem.Value.(yaml.MapSlice)
				if !ok {
					panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
				}
				pp := parseServer(serverValue)
				sshList[pp.FQDN] = pp.OS
			}
		} else {
			for _, groupItem := range *configs {
				groupValue, ok := groupItem.Value.(yaml.MapSlice)
				if !ok {
					panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
				}
				if groupItem.Key == servergroupname {
					for _, serverItem := range groupValue {
						var osnamecheck bool
						var servernamecheck bool
						serverValue, ok := serverItem.Value.(yaml.MapSlice)
						if !ok {
							panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
						}
						pp := parseServer(serverValue)
						if len(blstruct.exclude.osExcl) > 0 {
							for _, ve := range blstruct.exclude.osExcl {
								if strings.EqualFold(pp.OS, ve) {
									osnamecheck = true
								}
							}
						}
						if len(blstruct.exclude.serversExcl) > 0 {
							for _, ve := range blstruct.exclude.serversExcl {
								if strings.EqualFold(pp.FQDN, ve) {
									servernamecheck = true
								}
							}
						}
						if !servernamecheck && !osnamecheck {
							sshList[pp.FQDN] = pp.OS
						}
					}
				}
			}
		}
	} else {
		if len(blstruct.exclude.osExcl) == 0 &&
			len(blstruct.exclude.serversExcl) == 0 {
			for _, groupItem := range *configs {
				groupValue, ok := groupItem.Value.(yaml.MapSlice)
				if !ok {
					panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
				}
				if strings.EqualFold(groupItem.Key.(string), servergroupname) {
					for _, serverItem := range groupValue {
						serverValue, ok := serverItem.Value.(yaml.MapSlice)
						if !ok {
							panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
						}
						pp := parseServer(serverValue)
						sshList[pp.FQDN] = pp.OS
					}
				}
			}
		} else {
			for _, groupItem := range *configs {
				groupValue, ok := groupItem.Value.(yaml.MapSlice)
				if !ok {
					panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
				}
				if strings.EqualFold(groupItem.Key.(string), servergroupname) {
					for _, serverItem := range groupValue {
						var osnamecheck bool
						var servernamecheck bool
						serverValue, ok := serverItem.Value.(yaml.MapSlice)
						if !ok {
							panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
						}
						pp := parseServer(serverValue)
						if len(blstruct.exclude.osExcl) > 0 {
							for _, ve := range blstruct.exclude.osExcl {
								if strings.EqualFold(pp.OS, ve) {
									osnamecheck = true
								}
							}
						}
						if len(blstruct.exclude.serversExcl) > 0 {
							for _, ve := range blstruct.exclude.serversExcl {
								if strings.EqualFold(pp.FQDN, ve) {
									servernamecheck = true
								}
							}
						}
						if !servernamecheck && !osnamecheck {
							sshList[pp.FQDN] = pp.OS
						}
					}
				}
			}
		}
	}
	return sshList
}

// checkPrereqs builds read-only prerequisite probe steps (tool presence,
// downloaded files, VCS clones). Pure builder; no changes are made on hosts.
func (blstruct *ParsedBaseline) checkPrereqs(sshList map[string]string) []baselineStep {
	var steps []baselineStep
	if blstruct.prereq.cleanup {
		return steps
	}
	fmt.Printf("Prerequisites Checklist: ")
	if len(blstruct.prereq.vcs.urls) == 0 &&
		len(blstruct.prereq.files.urls) == 0 &&
		len(blstruct.prereq.tools) == 0 {
		fmt.Printf("Skipping...\n")
		return steps
	}
	fmt.Printf("\n")

	// tools present?
	fmt.Printf(" Prerequisite Tools: ")
	if len(blstruct.prereq.tools) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		for _, ve := range blstruct.prereq.tools {
			ve := ve
			cmds := buildCmds(sshList, nil, func(os string) string {
				return pkgmanlib.PkgSearch[os] + ve
			})
			steps = append(steps, baselineStep{label: "Tool present " + ve, cmds: cmds})
		}
	}

	// downloaded URL files present?
	fmt.Printf(" Prerequisite URL's: ")
	if len(blstruct.prereq.files.urls) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		for _, ve := range blstruct.prereq.files.urls {
			parseFile := strings.Split(ve, "/")
			parsedFile := parseFile[len(parseFile)-1]
			cmds := buildCmds(sshList, nil, func(string) string {
				return pkgmanlib.OmniTools["statinfo"] + parsedFile
			})
			steps = append(steps, baselineStep{label: "URL present " + parsedFile, cmds: cmds})
		}
	}

	// VCS clones present?
	fmt.Printf(" Prerequisite Files (via VCS): ")
	if len(blstruct.prereq.vcs.urls) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		for _, ve := range blstruct.prereq.vcs.urls {
			parseFile := strings.Split(ve, "/")
			parsedFile := parseFile[len(parseFile)-1]
			cmds := buildCmds(sshList, nil, func(string) string {
				return pkgmanlib.OmniTools["statinfo"] + parsedFile
			})
			steps = append(steps, baselineStep{label: "VCS present " + parsedFile, cmds: cmds})
		}
	}

	return steps
}

// checkMustHaves builds read-only "must have" compliance probes. Pure builder.
func (blstruct *ParsedBaseline) checkMustHaves(sshList map[string]string) []baselineStep {
	var steps []baselineStep
	mh := &blstruct.musthave
	fmt.Printf("Must Have Checklist: ")
	if len(mh.installed) == 0 &&
		len(mh.enabled) == 0 &&
		len(mh.disabled) == 0 &&
		len(mh.configured.services) == 0 &&
		len(mh.users.users) == 0 &&
		mh.policies.polimport == "" &&
		mh.policies.polstatus == "" &&
		len(mh.rules.fwopen.ports) == 0 &&
		len(mh.rules.fwclosed.ports) == 0 &&
		len(mh.rules.fwzones) == 0 &&
		len(mh.mounts.mountname) == 0 {
		fmt.Printf("Skipping...\n")
		return steps
	}
	fmt.Printf("\n")

	// Installed?
	fmt.Printf(" Installed: ")
	if len(mh.installed) > 0 {
		fmt.Printf("\n")
		for _, ve := range mh.installed {
			ve := ve
			cmds := buildCmds(sshList, nil, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "search")
			})
			steps = append(steps, baselineStep{label: "Installed? " + ve, cmds: cmds})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Enabled / active?
	fmt.Printf(" Enabled: ")
	if len(mh.enabled) > 0 {
		fmt.Printf("\n")
		for _, ve := range mh.enabled {
			ve := ve
			cmds := buildCmds(sshList, nil, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "isactive")
			})
			steps = append(steps, baselineStep{label: "Active? " + ve, cmds: cmds})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Disabled (probe active state)?
	fmt.Printf(" Disabled: ")
	if len(mh.disabled) > 0 {
		fmt.Printf("\n")
		for _, ve := range mh.disabled {
			if ve == "" {
				continue
			}
			ve := ve
			cmds := buildCmds(sshList, nil, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "isactive")
			})
			// Must-Have Disabled wants the service INACTIVE, so the is-active
			// probe is inverted: NOK (inactive) == compliant, OK (active) == not.
			steps = append(steps, baselineStep{label: "Inactive? " + ve, cmds: cmds, invert: true})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Configured files match? Compute each config file's md5 locally (at build
	// time) and verify the remote copy matches with `md5sum -c`. Exit 0 (match)
	// reports OK; any mismatch or missing remote file reports NOK.
	// TODO: check has no per-host root-ness, so a non-root-readable destination
	// will fail; add sudo support to the check path once it carries isRoot.
	fmt.Printf(" Configured Checklist: ")
	if len(mh.configured.services) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		for ke, ve := range mh.configured.services {
			if ke == "" {
				continue
			}
			if len(ve.source) != len(ve.destination) {
				fmt.Printf("      %s: config source/destination mismatch, skipping\n", ke)
				continue
			}
			ke, ve := ke, ve
			for i := range ve.source {
				data, err := os.ReadFile(ve.source[i])
				if err != nil {
					fmt.Printf("      %s: cannot read config source %q: %v, skipping\n", ke, ve.source[i], err)
					continue
				}
				localMD5 := fmt.Sprintf("%x", md5.Sum(data))
				probe := "echo '" + localMD5 + "  " + ve.destination[i] + "' | md5sum -c --status -"
				cmds := make(map[string]string)
				for host := range sshList {
					cmds[host] = probe
				}
				steps = append(steps, baselineStep{label: "Config? " + ke, cmds: cmds})
			}
		}
	}

	// Users exist?
	fmt.Printf(" Users Checklist: ")
	if len(mh.users.users) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		for ke, ve := range mh.users.users {
			if ke == "" {
				continue
			}
			ke, ve := ke, ve
			cmds := buildCmds(sshList, nil, func(string) string {
				return ve.userManagementCommandBuilder(&ke, "check")
			})
			steps = append(steps, baselineStep{label: "User? " + ke, cmds: cmds})
		}
	}

	// Policies status?
	fmt.Printf(" Policies Checklist: ")
	if mh.policies.polstatus == "" && mh.policies.polimport == "" {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		if mh.policies.polstatus != "" {
			cmds := buildCmds(sshList, nil, func(string) string {
				return pkgmanlib.OmniTools["policystatus"]
			})
			steps = append(steps, baselineStep{label: "Policy status", cmds: cmds})
		}
		if mh.policies.polimport != "" {
			cmds := buildCmds(sshList, nil, func(string) string {
				return pkgmanlib.OmniTools["policycheck"]
			})
			steps = append(steps, baselineStep{label: "Policy import", cmds: cmds})
		}
	}

	// Firewall present?
	fmt.Printf(" Firewall Checklist: ")
	if len(mh.rules.fwopen.ports) == 0 && len(mh.rules.fwclosed.ports) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		steps = append(steps, firewallSteps(sshList, nil, mh.rules.fwopen.ports, mh.rules.fwopen.protocols, mh.rules.fwzones, "check")...)
		steps = append(steps, firewallSteps(sshList, nil, mh.rules.fwclosed.ports, mh.rules.fwclosed.protocols, mh.rules.fwzones, "check")...)
	}

	// Mounts present in fstab?
	fmt.Printf(" Mounts Checklist: ")
	if len(mh.mounts.mountname) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		for ke, ve := range mh.mounts.mountname {
			if ke == "" || ve.address == "" {
				continue
			}
			ve := ve
			cmds := buildCmds(sshList, nil, func(string) string {
				return "grep '" + ve.address + "' /etc/fstab"
			})
			steps = append(steps, baselineStep{label: "Mount? " + ke, cmds: cmds})
		}
	}

	return steps
}

// checkMustNotHaves builds read-only "must not have" compliance probes. Pure
// builder. A statusOK from these probes means the unwanted item IS present
// (i.e. non-compliant) — interpretation is left to reporting.
func (blstruct *ParsedBaseline) checkMustNotHaves(sshList map[string]string) []baselineStep {
	var steps []baselineStep
	mnh := &blstruct.mustnothave
	fmt.Printf("Must Not Have Checklist: ")
	if len(mnh.installed) == 0 &&
		len(mnh.enabled) == 0 &&
		len(mnh.disabled) == 0 &&
		len(mnh.users) == 0 &&
		len(mnh.rules.fwopen.ports) == 0 &&
		len(mnh.rules.fwclosed.ports) == 0 &&
		len(mnh.rules.fwzones) == 0 &&
		len(mnh.mounts) == 0 {
		fmt.Printf("Skipping...\n")
		return steps
	}
	fmt.Printf("\n")

	// Installed?
	fmt.Printf(" Installed Checklist: ")
	if len(mnh.installed) > 0 {
		fmt.Printf("\n")
		for _, ve := range mnh.installed {
			ve := ve
			cmds := buildCmds(sshList, nil, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "search")
			})
			steps = append(steps, baselineStep{label: "Installed? " + ve, cmds: cmds, invert: true})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Enabled / active?
	fmt.Printf(" Enabled Checklist: ")
	if len(mnh.enabled) > 0 {
		fmt.Printf("\n")
		for _, ve := range mnh.enabled {
			ve := ve
			cmds := buildCmds(sshList, nil, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "isactive")
			})
			steps = append(steps, baselineStep{label: "Active? " + ve, cmds: cmds, invert: true})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Disabled probe
	fmt.Printf(" Disabled Checklist: ")
	if len(mnh.disabled) > 0 {
		fmt.Printf("\n")
		for _, ve := range mnh.disabled {
			if ve == "" {
				continue
			}
			ve := ve
			cmds := buildCmds(sshList, nil, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "isactive")
			})
			steps = append(steps, baselineStep{label: "Active? " + ve, cmds: cmds})
		}
	}

	// Users?
	fmt.Printf(" Users Checklist: ")
	if len(mnh.users) > 0 {
		fmt.Printf("\n")
		for _, ve := range mnh.users {
			if ve == "" {
				continue
			}
			ve := ve
			cmds := buildCmds(sshList, nil, func(string) string {
				return pkgmanlib.OmniTools["userinfo"] + ve
			})
			steps = append(steps, baselineStep{label: "User? " + ve, cmds: cmds, invert: true})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Firewall?
	fmt.Printf(" Firewall Checklist: ")
	if len(mnh.rules.fwopen.ports) == 0 && len(mnh.rules.fwclosed.ports) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		fwSteps := firewallSteps(sshList, nil, mnh.rules.fwopen.ports, mnh.rules.fwopen.protocols, mnh.rules.fwzones, "check")
		fwSteps = append(fwSteps, firewallSteps(sshList, nil, mnh.rules.fwclosed.ports, mnh.rules.fwclosed.protocols, mnh.rules.fwzones, "check")...)
		for i := range fwSteps {
			fwSteps[i].invert = true
		}
		steps = append(steps, fwSteps...)
	}

	// Mounts?
	fmt.Printf(" Mounts Checklist: ")
	if len(mnh.mounts) > 0 {
		fmt.Printf("\n")
		for _, ve := range mnh.mounts {
			ve := ve
			cmds := buildCmds(sshList, nil, func(string) string {
				return "grep '" + ve + "' /etc/fstab"
			})
			steps = append(steps, baselineStep{label: "Mount? " + ve, cmds: cmds, invert: true})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	return steps
}
