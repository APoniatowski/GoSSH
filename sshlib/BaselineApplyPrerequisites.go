package sshlib

import (
	"fmt"
	"path"
	"strings"
)

// applyPrereq builds the ordered prerequisite steps (tools, downloaded files,
// transfers, mounts, VCS, custom commands) for every live host. It is a pure
// builder: it returns the steps to dispatch and performs no I/O of its own.
func (baselineStruct *ParsedBaseline) applyPrereq(sshList map[string]string, isRoot map[string]bool) []baselineStep {
	var steps []baselineStep
	fmt.Printf("Prerequisites Checklist: ")
	if len(baselineStruct.prereq.vcs.execute) == 0 &&
		len(baselineStruct.prereq.vcs.urls) == 0 &&
		baselineStruct.prereq.files.local.dest == "" &&
		baselineStruct.prereq.files.local.src == "" &&
		baselineStruct.prereq.files.remote.address == "" &&
		baselineStruct.prereq.files.remote.dest == "" &&
		baselineStruct.prereq.files.remote.mounttype == "" &&
		baselineStruct.prereq.files.remote.pwd == "" &&
		baselineStruct.prereq.files.remote.src == "" &&
		baselineStruct.prereq.files.remote.username == "" &&
		len(baselineStruct.prereq.files.remote.files) == 0 &&
		len(baselineStruct.prereq.files.urls) == 0 &&
		len(baselineStruct.prereq.tools) == 0 &&
		baselineStruct.prereq.script == "" &&
		len(baselineStruct.prereq.commands) == 0 &&
		!baselineStruct.prereq.cleanup {
		fmt.Printf("Skipping...\n")
		return steps
	}
	fmt.Printf("\n")

	// prerequisite tools
	fmt.Printf(" Prerequisite Tools: ")
	if len(baselineStruct.prereq.tools) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		for _, ve := range baselineStruct.prereq.tools {
			ve := ve
			cmds := buildCmds(sshList, isRoot, func(os string) string {
				return serviceCommandBuilder(&ve, &os, "install")
			})
			steps = append(steps, baselineStep{label: "Tool " + ve, cmds: cmds})
		}
	}

	// prerequisite files URLs
	fmt.Printf(" Prerequisite URL's: ")
	if len(baselineStruct.prereq.files.urls) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		for _, ve := range baselineStruct.prereq.files.urls {
			ve := ve
			cmds := buildCmds(sshList, nil, func(string) string {
				return prereqURLFetch(&ve, &baselineStruct.prereq.cleanup)
			})
			steps = append(steps, baselineStep{label: "URL " + ve, cmds: cmds})
		}
	}

	// prerequisite files local (network transfer)
	// Pushes the local source to the destination by streaming its bytes as
	// base64 inside the command channel (see fileToCommand). The old marker's
	// cleanup/symlink behaviour is not yet supported on this path.
	fmt.Printf(" Prerequisite Files (network transfer): ")
	if baselineStruct.prereq.files.local.src == "" {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		src := baselineStruct.prereq.files.local.src
		dest := baselineStruct.prereq.files.local.dest
		cmds := make(map[string]string)
		buildErr := false
		for host := range sshList {
			cmd, err := fileToCommand(src, dest, !isRoot[host])
			if err != nil {
				fmt.Printf("      %s: cannot read local source %q: %v, skipping\n", host, src, err)
				buildErr = true
				break
			}
			cmds[host] = cmd
		}
		if !buildErr {
			steps = append(steps, baselineStep{label: "Transfer " + src, cmds: cmds})
		}
	}

	// prerequisite files remote (via mount)
	fmt.Printf(" Prerequisite Files (via mount): ")
	if baselineStruct.prereq.files.remote.src == "" || len(baselineStruct.prereq.files.remote.files) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		for _, ve := range baselineStruct.prereq.files.remote.files {
			ve := ve
			command := baselineStruct.prereq.files.remote.remoteFilesCommandBuilder(&ve, "apply")
			symlinker := prereqRemoteCleanup(&baselineStruct.prereq.files.remote.dest, &ve, &baselineStruct.prereq.cleanup)
			cmds := buildCmds(sshList, isRoot, func(string) string {
				return command + symlinker
			})
			steps = append(steps, baselineStep{label: "Mount file " + ve, cmds: cmds})
		}
	}

	// prerequisite VCS instructions
	fmt.Printf(" Prerequisite Files (via VCS): ")
	if len(baselineStruct.prereq.vcs.execute) == 0 && len(baselineStruct.prereq.vcs.urls) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		for _, ve := range baselineStruct.prereq.vcs.urls {
			ve := ve
			parseFile := strings.Split(ve, "/")
			_ = parseFile[len(parseFile)-1]
			cmds := buildCmds(sshList, nil, func(string) string {
				return prereqURLFetch(&ve, &baselineStruct.prereq.cleanup)
			})
			steps = append(steps, baselineStep{label: "VCS clone " + ve, cmds: cmds})
		}
		for _, ve := range baselineStruct.prereq.vcs.execute {
			ve := ve
			cmds := buildCmds(sshList, nil, func(string) string {
				return ve
			})
			steps = append(steps, baselineStep{label: "VCS exec " + ve, cmds: cmds})
		}
	}

	// prerequisite script: push a local script to /tmp, make it executable, run
	// it (sudo when the host is non-root). Mirrors Final Scripts.
	if baselineStruct.prereq.script != "" {
		fmt.Printf(" Prerequisite Script: \n")
		p := baselineStruct.prereq.script
		base := path.Base(p)
		remote := "/tmp/" + base
		cmds := make(map[string]string)
		skip := false
		for host := range sshList {
			push, err := fileToCommand(p, remote, false)
			if err != nil {
				fmt.Printf("      %s: cannot read script %q: %v, skipping\n", host, p, err)
				skip = true
				break
			}
			cmds[host] = push + " && chmod +x '" + remote + "' && " + withSudo(isRoot[host], "'"+remote+"'")
		}
		if !skip {
			steps = append(steps, baselineStep{label: "Prereq script " + base, cmds: cmds})
		}
	}

	// prerequisite custom commands
	if len(baselineStruct.prereq.commands) > 0 {
		fmt.Printf(" Prerequisite Commands: \n")
		for _, ve := range baselineStruct.prereq.commands {
			ve := ve
			cmds := buildCmds(sshList, isRoot, func(string) string {
				return ve
			})
			steps = append(steps, baselineStep{label: "Command " + ve, cmds: cmds})
		}
	}

	return steps
}
