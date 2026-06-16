package sshlib

import (
	"fmt"

	"github.com/APoniatowski/GoSSH/pkgmanlib"
)

// applyMustNotHaves builds the ordered "must not have" steps (uninstall, disable,
// remove users, remove firewall rules, unmount). Pure builder.
func (baselineStruct *ParsedBaseline) applyMustNotHaves(sshList map[string]string, isRoot map[string]bool) []baselineStep {
	var steps []baselineStep
	mnh := &baselineStruct.mustnothave
	fmt.Printf("Must Not Have Checklist: ")
	if len(mnh.installed) == 0 &&
		len(mnh.enabled) == 0 &&
		len(mnh.disabled) == 0 &&
		len(mnh.users) == 0 &&
		len(mnh.rules.fwopen.ports) == 0 &&
		len(mnh.rules.fwopen.protocols) == 0 &&
		len(mnh.rules.fwclosed.ports) == 0 &&
		len(mnh.rules.fwclosed.protocols) == 0 &&
		len(mnh.rules.fwzones) == 0 &&
		len(mnh.mounts) == 0 {
		fmt.Printf("Skipping...\n")
		return steps
	}
	fmt.Printf("\n")

	// Installed -> uninstall
	fmt.Printf(" Installed Checklist: ")
	if len(mnh.installed) > 0 {
		fmt.Printf("\n")
		for _, ve := range mnh.installed {
			ve := ve
			cmds := buildCmds(sshList, isRoot, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "uninstall")
			})
			steps = append(steps, baselineStep{label: "Uninstall " + ve, cmds: cmds})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Enabled -> disable
	fmt.Printf(" Enabled Checklist: ")
	if len(mnh.enabled) > 0 {
		fmt.Printf("\n")
		for _, ve := range mnh.enabled {
			ve := ve
			cmds := buildCmds(sshList, isRoot, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "disable")
			})
			steps = append(steps, baselineStep{label: "Disable " + ve, cmds: cmds})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Disabled -> enable
	fmt.Printf(" Disabled Checklist: ")
	if len(mnh.disabled) > 0 {
		fmt.Printf("\n")
		for _, ve := range mnh.disabled {
			if ve == "" {
				continue
			}
			ve := ve
			cmds := buildCmds(sshList, isRoot, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "enable")
			})
			steps = append(steps, baselineStep{label: "Enable " + ve, cmds: cmds})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Users -> delete
	fmt.Printf(" Users Checklist: ")
	if len(mnh.users) > 0 {
		fmt.Printf("\n")
		for _, ve := range mnh.users {
			if ve == "" {
				continue
			}
			ve := ve
			cmds := buildCmds(sshList, isRoot, func(string) string {
				return pkgmanlib.OmniTools["userdel"] + ve
			})
			steps = append(steps, baselineStep{label: "Delete user " + ve, cmds: cmds})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	// Firewall rules -> remove
	fmt.Printf(" Firewall Checklist: ")
	if len(mnh.rules.fwopen.ports) == 0 &&
		len(mnh.rules.fwopen.protocols) == 0 &&
		len(mnh.rules.fwclosed.ports) == 0 &&
		len(mnh.rules.fwclosed.protocols) == 0 &&
		len(mnh.rules.fwzones) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		steps = append(steps, firewallSteps(sshList, isRoot, mnh.rules.fwopen.ports, mnh.rules.fwopen.protocols, mnh.rules.fwzones, "remove-open")...)
		steps = append(steps, firewallSteps(sshList, isRoot, mnh.rules.fwclosed.ports, mnh.rules.fwclosed.protocols, mnh.rules.fwzones, "remove-closed")...)
	}

	// Mounts -> check fstab presence (removal pending implementation)
	fmt.Printf(" Mounts Checklist: ")
	if len(mnh.mounts) > 0 {
		fmt.Printf("\n")
		for _, ve := range mnh.mounts {
			ve := ve
			cmds := buildCmds(sshList, isRoot, func(string) string {
				return "grep '" + ve + "' /etc/fstab"
			})
			steps = append(steps, baselineStep{label: "Mount check " + ve, cmds: cmds})
		}
	} else {
		fmt.Printf("Skipping...\n")
	}

	return steps
}
