package sshlib

import (
	"fmt"
)

// TODO remove fmt.printf's here... took a while to find them all
func (baselineStruct *ParsedBaseline) applyFinals(sshList *map[string]string, rebootBool *bool, commandChannel chan<- map[string]string) {
	commandSet := make(map[string]string)
	// Final steps list
	fmt.Println("Applying final instructions:")
	if len(baselineStruct.final.scripts) == 0 &&
		len(baselineStruct.final.commands) == 0 &&
		len(baselineStruct.final.collect.logs) == 0 &&
		len(baselineStruct.final.collect.stats) == 0 &&
		len(baselineStruct.final.collect.files) == 0 &&
		!baselineStruct.final.collect.users &&
		!baselineStruct.final.restart.services &&
		!baselineStruct.final.restart.servers {
		// fmt.Println("No final steps have been specified  -- Please check your baseline, if you believe this to be incorrect")
	} else {
		// final scripts
		fmt.Println("  Execute: ")
		fmt.Printf("    Scripts: ")
		if len(baselineStruct.final.scripts) > 0 {
			fmt.Printf("\n")
			for _, ve := range baselineStruct.final.scripts {
				// TODO Final scripts that need to be transferred
				// transfer file to /tmp
				// execute script
				fmt.Println(ve)
			}
			commandChannel <- commandSet
		} else {
			fmt.Printf("Skipping...\n")
		}
		// final commands
		fmt.Printf("    Commands: ")
		if len(baselineStruct.final.commands) > 0 {
			fmt.Printf("\n")
			for _, ve := range baselineStruct.final.commands {
				for key, val := range *sshList {
					if commandSet[val] == "" {
						// TODO Final commands
						commandSet[key] = finalCommandBuilder(&ve,"command")
					}
				}
			}
			commandChannel <- commandSet
		} else {
			fmt.Printf("Skipping...\n")
		}
		// final collections
		if len(baselineStruct.final.collect.logs) == 0 &&
			len(baselineStruct.final.collect.stats) == 0 &&
			len(baselineStruct.final.collect.files) == 0 &&
			!baselineStruct.final.collect.users {
			// fmt.Println("No collections specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			fmt.Println("  Collect: ")
			fmt.Printf("    Logs:")
			if len(baselineStruct.final.collect.logs) > 0 {
				fmt.Printf("\n")
				for _, ve := range baselineStruct.final.collect.logs {
					// TODO transfer to ./collections/[servername]/logs
					finalCommandBuilder(&ve,"logs")
				}
				commandChannel <- commandSet
			} else {
				fmt.Printf("Skipping...\n")
			}
			fmt.Printf("    Stats:")
			if len(baselineStruct.final.collect.stats) > 0 {
				fmt.Printf("\n")
				for _, ve := range baselineStruct.final.collect.stats {
					// TODO transfer to ./collections/[servername]/stats
					finalCommandBuilder(&ve,"stats")
				}
				commandChannel <- commandSet
			} else {
				fmt.Printf("Skipping...\n")
			}
			fmt.Printf("    Files: ")
			if len(baselineStruct.final.collect.files) > 0 {
				fmt.Printf("\n")
				for _, ve := range baselineStruct.final.collect.files {
					// TODO transfer to ./collections/[servername]/files
					finalCommandBuilder(&ve,"files")
				}
				commandChannel <- commandSet
			} else {
				fmt.Printf("Skipping...\n")
			}
			fmt.Printf("    Users: ")
			if baselineStruct.final.collect.users {
				fmt.Printf("\n")
				for key, val := range *sshList {
					if commandSet[val] == "" {
						// TODO write to log in collections?
						commandSet[key] = "w"
					}
				}
				commandChannel <- commandSet
			} else {
				fmt.Printf("Skipping...\n")
			}
		}
		// final restarts
		fmt.Println("  Restart: ")
		if !baselineStruct.final.restart.services &&
			!baselineStruct.final.restart.servers &&
			!*rebootBool {
			fmt.Printf("Skipping...\n")
		} else {
			fmt.Printf("    Services:")
			// fmt.Println("Reboot:")
			if baselineStruct.final.restart.services {
				fmt.Printf("\n")
				for key, val := range *sshList {
					if commandSet[val] == "" {
						// TODO Final service restart
						commandSet[key] = "systemctl --daemon-reload"
					}
				}
				commandChannel <- commandSet
			} else {
				fmt.Printf("Skipping...\n")
			}
			fmt.Printf("    Servers:")
			if *rebootBool || baselineStruct.final.restart.servers {
				fmt.Printf("\n")
				for key, val := range *sshList {
					if commandSet[val] == "" {
						// TODO Final reboot
						commandSet[key] = "reboot"
					}
				}
				commandChannel <- commandSet
			} else {
				fmt.Printf("Skipping...\n")
			}
		}
	}
	return
}
