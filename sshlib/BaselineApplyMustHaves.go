package sshlib

import "fmt"

// applyMustHaves builds the ordered "must have" steps (installs, enables,
// disables, users, policies, firewall rules, mounts). Pure builder. It also
// flips rebootBool when a policy requires a reboot, consumed later by finals.
func (baselineStruct *ParsedBaseline) applyMustHaves(sshList map[string]string, isRoot map[string]bool, rebootBool *bool) []baselineStep {
	var steps []baselineStep
	mh := &baselineStruct.musthave
	fmt.Printf("Must Have Checklist: ")
	if len(mh.installed) == 0 &&
		len(mh.enabled) == 0 &&
		len(mh.disabled) == 0 &&
		len(mh.configured.services) == 0 &&
		len(mh.users.users) == 0 &&
		mh.policies.polimport == "" &&
		!mh.policies.polreboot &&
		mh.policies.polstatus == "" &&
		len(mh.rules.fwopen.ports) == 0 &&
		len(mh.rules.fwopen.protocols) == 0 &&
		len(mh.rules.fwclosed.ports) == 0 &&
		len(mh.rules.fwclosed.protocols) == 0 &&
		len(mh.rules.fwzones) == 0 &&
		len(mh.mounts.mountname) == 0 {
		fmt.Printf("Skipping...\n")
		return steps
	}
	fmt.Printf("\n")

	// Installed
	fmt.Printf(" Installed: ")
	if len(mh.installed) > 0 {
		fmt.Printf("\n")
		for _, ve := range mh.installed {
			ve := ve
			cmds := buildCmds(sshList, isRoot, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "install")
			})
			steps = append(steps, baselineStep{label: "Install " + ve, cmds: cmds})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Configured services (config-file transfer). Each config file is streamed
	// to its destination as base64 over the command channel (see fileToCommand).
	// NOTE: Configured (and Users) run BEFORE Enable/Disable so a service is
	// started with its config already in place and required users present. (For a
	// config change on an ALREADY-running service, systemd still needs an explicit
	// restart/reload — use Final Restart; ordering alone does not reload it.)
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
				src, dest := ve.source[i], ve.destination[i]
				cmds := make(map[string]string)
				skip := false
				for host := range sshList {
					cmd, err := fileToCommand(src, dest, !isRoot[host])
					if err != nil {
						fmt.Printf("      %s: cannot read config source %q: %v, skipping\n", ke, src, err)
						skip = true
						break
					}
					cmds[host] = cmd
				}
				if skip {
					continue
				}
				steps = append(steps, baselineStep{label: "Config " + ke, cmds: cmds})
			}
		}
	}

	// Users
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
			cmds := buildCmds(sshList, isRoot, func(string) string {
				return ve.userManagementCommandBuilder(&ke, "add")
			})
			steps = append(steps, baselineStep{label: "User " + ke, cmds: cmds})
		}
	}

	// Enabled (after Configured/Users so the service starts with config in place)
	fmt.Printf(" Enabled: ")
	if len(mh.enabled) > 0 {
		fmt.Printf("\n")
		for _, ve := range mh.enabled {
			ve := ve
			cmds := buildCmds(sshList, isRoot, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "enable")
			})
			steps = append(steps, baselineStep{label: "Enable " + ve, cmds: cmds})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Disabled
	fmt.Printf(" Disabled: ")
	if len(mh.disabled) > 0 {
		fmt.Printf("\n")
		for _, ve := range mh.disabled {
			if ve == "" {
				continue
			}
			ve := ve
			cmds := buildCmds(sshList, isRoot, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "disable")
			})
			steps = append(steps, baselineStep{label: "Disable " + ve, cmds: cmds})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Policies
	fmt.Printf(" Policies Checklist: ")
	if mh.policies.polstatus == "" && mh.policies.polimport == "" && !mh.policies.polreboot {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		if mh.policies.polreboot {
			*rebootBool = true
		}
		cmds := buildCmds(sshList, isRoot, func(string) string {
			return mh.policies.policyCommandBuilder("apply")
		})
		steps = append(steps, baselineStep{label: "Policies", cmds: cmds})
	}

	// Firewall rules
	fmt.Printf(" Firewall Checklist: ")
	if len(mh.rules.fwopen.ports) == 0 &&
		len(mh.rules.fwopen.protocols) == 0 &&
		len(mh.rules.fwclosed.ports) == 0 &&
		len(mh.rules.fwclosed.protocols) == 0 &&
		len(mh.rules.fwzones) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		steps = append(steps, firewallSteps(sshList, isRoot, mh.rules.fwopen.ports, mh.rules.fwopen.protocols, mh.rules.fwzones, "apply-open")...)
		steps = append(steps, firewallSteps(sshList, isRoot, mh.rules.fwclosed.ports, mh.rules.fwclosed.protocols, mh.rules.fwzones, "apply-closed")...)
	}

	// Mounts
	fmt.Printf(" Mounts Checklist: ")
	if len(mh.mounts.mountname) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		for ke, ve := range mh.mounts.mountname {
			if ke == "" {
				continue
			}
			if ve.mounttype == "" || ve.address == "" || ve.src == "" || ve.dest == "" {
				fmt.Printf("      %s: critical mount info missing, skipping\n", ke)
				continue
			}
			ve := ve
			cmds := buildCmds(sshList, isRoot, func(string) string {
				return ve.mountCommandBuilder("apply")
			})
			steps = append(steps, baselineStep{label: "Mount " + ke, cmds: cmds})
		}
	}

	return steps
}

// firewallSteps builds one step per port/protocol pair (optionally per zone).
// Shared by must-have (apply-*), must-not-have (remove-*), and the read-only
// check path. When action == "check" the commands are built without sudo
// (isRoot is ignored) and labelled "fw check <port>/<proto>"; the apply actions
// keep per-host sudo and the "<action> <port>/<proto>" label.
func firewallSteps(sshList map[string]string, isRoot map[string]bool, ports, protocols, zones []string, action string) []baselineStep {
	var steps []baselineStep
	if len(ports) != len(protocols) {
		fmt.Println("There seems to be inconsistencies between your firewall ports and protocols.")
		fmt.Println("Please review your baseline and rectify it.")
		return steps
	}
	check := action == "check"
	build := func(port, protocol, zone string) {
		cmds := make(map[string]string)
		for host := range sshList {
			cmd := firewallCommandBuilder(&port, &protocol, &zone, action)
			if !check {
				cmd = withSudo(isRoot[host], cmd)
			}
			cmds[host] = cmd
		}
		labelAction := action
		if check {
			labelAction = "fw check"
		}
		label := labelAction + " " + port + "/" + protocol
		if zone != "" {
			label += " zone " + zone
		}
		steps = append(steps, baselineStep{label: label, cmds: cmds})
	}
	if len(zones) > 0 {
		for _, zone := range zones {
			for i := range ports {
				build(ports[i], protocols[i], zone)
			}
		}
	} else {
		for i := range ports {
			build(ports[i], protocols[i], "")
		}
	}
	return steps
}
