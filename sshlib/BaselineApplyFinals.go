package sshlib

import (
	"fmt"
	"path"
)

// applyFinals builds the ordered final steps (custom commands, collections,
// restarts/reboots). Pure builder. rebootBool (set by must-have policies) forces
// a server reboot step. Script transfer and artifact collection are not yet
// implemented end-to-end and are emitted as no-ops with a notice.
func (baselineStruct *ParsedBaseline) applyFinals(sshList map[string]string, isRoot map[string]bool, rebootBool *bool) []baselineStep {
	var steps []baselineStep
	fin := &baselineStruct.final
	fmt.Println("Applying final instructions:")
	if len(fin.scripts) == 0 &&
		len(fin.commands) == 0 &&
		len(fin.collect.logs) == 0 &&
		len(fin.collect.stats) == 0 &&
		len(fin.collect.files) == 0 &&
		!fin.collect.users &&
		!fin.restart.services &&
		!fin.restart.servers &&
		!*rebootBool {
		return steps
	}

	// Final commands
	if len(fin.commands) > 0 {
		fmt.Printf("  Commands:\n")
		for _, ve := range fin.commands {
			ve := ve
			cmds := buildCmds(sshList, isRoot, func(string) string {
				return finalCommandBuilder(&ve, "command")
			})
			steps = append(steps, baselineStep{label: "Final cmd " + ve, cmds: cmds})
		}
	}

	// Final scripts: push each script to /tmp as base64 over the command
	// channel, make it executable, then run it (sudo when the host is non-root).
	if len(fin.scripts) > 0 {
		fmt.Printf("  Scripts:\n")
		for _, p := range fin.scripts {
			p := p
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
			if skip {
				continue
			}
			steps = append(steps, baselineStep{label: "Script " + base, cmds: cmds})
		}
	}

	// Final collections: run a command per host, capture its output, and let the
	// orchestrator write each host's artifact under ./collections/<host>/.
	if len(fin.collect.logs) > 0 || len(fin.collect.stats) > 0 || len(fin.collect.files) > 0 {
		fmt.Printf("  Collect:\n")

		// Logs: capture the service journal for today. --no-pager is REQUIRED:
		// the executor runs commands over a PTY, so without it journalctl pipes to
		// a pager (less) and blocks forever (-> dispatch timeout -> host dead).
		for _, s := range fin.collect.logs {
			s := s
			cmds := buildCmds(sshList, nil, func(string) string {
				return "journalctl --no-pager -u " + s + " -S today --no-tail 2>&1"
			})
			steps = append(steps, baselineStep{label: "Collect log " + s, cmds: cmds, collectAs: s + ".log"})
		}

		// Stats: reuse the final command builder; skip entries with no command.
		for _, s := range fin.collect.stats {
			s := s
			cmd := finalCommandBuilder(&s, "stats")
			if cmd == "" {
				continue
			}
			cmds := buildCmds(sshList, nil, func(string) string {
				return cmd
			})
			steps = append(steps, baselineStep{label: "Collect stat " + s, cmds: cmds, collectAs: s + ".stat"})
		}

		// Files: cat each requested file.
		for _, fpath := range fin.collect.files {
			fpath := fpath
			base := path.Base(fpath)
			cmds := buildCmds(sshList, nil, func(string) string {
				return "cat '" + fpath + "' 2>&1"
			})
			steps = append(steps, baselineStep{label: "Collect file " + base, cmds: cmds, collectAs: base})
		}
	}

	// Collect logged-in users
	if fin.collect.users {
		cmds := buildCmds(sshList, nil, func(string) string {
			return "w"
		})
		steps = append(steps, baselineStep{label: "Collect users", cmds: cmds, collectAs: "users.txt"})
	}

	// Restart services
	if fin.restart.services {
		cmds := buildCmds(sshList, isRoot, func(string) string {
			return "systemctl daemon-reload"
		})
		steps = append(steps, baselineStep{label: "Reload services", cmds: cmds})
	}

	// Reboot servers
	if fin.restart.servers || *rebootBool {
		cmds := buildCmds(sshList, isRoot, func(string) string {
			return "shutdown -r +1"
		})
		steps = append(steps, baselineStep{label: "Reboot", cmds: cmds})
	}

	return steps
}
