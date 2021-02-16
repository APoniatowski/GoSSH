package sshlib

import (
	"fmt"
)

// TODO remove fmt.printf's here... took a while to find them all
func (blstruct *ParsedBaseline) applyFinals(sshList *map[string]string, rebootBool *bool) {
	commandset := make(map[string]string)
	// Final steps list
	fmt.Println("Applying final instructions:")
	if len(blstruct.final.scripts) == 0 &&
		len(blstruct.final.commands) == 0 &&
		len(blstruct.final.collect.logs) == 0 &&
		len(blstruct.final.collect.stats) == 0 &&
		len(blstruct.final.collect.files) == 0 &&
		!blstruct.final.collect.users &&
		!blstruct.final.restart.services &&
		!blstruct.final.restart.servers {
		// fmt.Println("No final steps have been specified  -- Please check your baseline, if you believe this to be incorrect")
	} else {
		// final scripts
		fmt.Println("  Execute: ")
		fmt.Printf("    Scripts: ")
		if len(blstruct.final.scripts) > 0 {
			fmt.Printf("\n")
			for _, ve := range blstruct.final.scripts {
				// TODO Final scripts that need to be transferred
				// transfer file to /tmp
				// execute script
				fmt.Println(ve)
			}
		} else {
			fmt.Printf("Skipping...\n")
		}
		// final commands
		fmt.Printf("    Commands: ")
		if len(blstruct.final.commands) > 0 {
			fmt.Printf("\n")
			for _, ve := range blstruct.final.commands {
				for key, val := range *sshList {
					if commandset[val] == "" {
						// TODO Final commands
						commandset[key] = finalCommandBuilder(&ve,"command")
					}
				}
			}
		} else {
			fmt.Printf("Skipping...\n")
		}
		// final collections
		if len(blstruct.final.collect.logs) == 0 &&
			len(blstruct.final.collect.stats) == 0 &&
			len(blstruct.final.collect.files) == 0 &&
			!blstruct.final.collect.users {
			// fmt.Println("No collections specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			fmt.Println("  Collect: ")
			fmt.Printf("    Logs:")
			if len(blstruct.final.collect.logs) > 0 {
				fmt.Printf("\n")
				for _, ve := range blstruct.final.collect.logs {
					// TODO transfer to ./collections/[servername]/logs
					finalCommandBuilder(&ve,"logs")
				}
			} else {
				fmt.Printf("Skipping...\n")
			}
			fmt.Printf("    Stats:")
			if len(blstruct.final.collect.stats) > 0 {
				fmt.Printf("\n")
				for _, ve := range blstruct.final.collect.stats {
					// TODO transfer to ./collections/[servername]/stats
					finalCommandBuilder(&ve,"stats")
				}
			} else {
				fmt.Printf("Skipping...\n")
			}
			fmt.Printf("    Files: ")
			if len(blstruct.final.collect.files) > 0 {
				fmt.Printf("\n")
				for _, ve := range blstruct.final.collect.files {
					// TODO transfer to ./collections/[servername]/files
					finalCommandBuilder(&ve,"files")
				}
			} else {
				fmt.Printf("Skipping...\n")
			}
			fmt.Printf("    Users: ")
			if blstruct.final.collect.users {
				fmt.Printf("\n")
				for key, val := range *sshList {
					if commandset[val] == "" {
						// TODO write to log in collections?
						commandset[key] = "w"
					}
				}
			} else {
				fmt.Printf("Skipping...\n")
			}
		}
		// final restarts
		fmt.Println("  Restart: ")
		if !blstruct.final.restart.services &&
			!blstruct.final.restart.servers &&
			!*rebootBool {
			fmt.Printf("Skipping...\n")
		} else {
			fmt.Printf("    Services:")
			// fmt.Println("Reboot:")
			if blstruct.final.restart.services {
				fmt.Printf("\n")
				for key, val := range *sshList {
					if commandset[val] == "" {
						// TODO Final service restart
						commandset[key] = "systemctl --daemon-reload"
					}
				}
			} else {
				fmt.Printf("Skipping...\n")
			}
			fmt.Printf("    Servers:")
			if *rebootBool || blstruct.final.restart.servers {
				fmt.Printf("\n")
				for key, val := range *sshList {
					if commandset[val] == "" {
						// TODO Final reboot
						commandset[key] = "reboot"
					}
				}
			} else {
				fmt.Printf("Skipping...\n")
			}
		}
	}
	return
}
