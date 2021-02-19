package sshlib

import "fmt"

func (blstruct *ParsedBaseline) applyMustHaves(sshList *map[string]string, rebootBool *bool) {
	commandset := make(map[string]string)
	// MH list
	fmt.Printf("Must Have Checklist: ")
	if len(blstruct.musthave.installed) == 0 &&
		len(blstruct.musthave.enabled) == 0 &&
		len(blstruct.musthave.disabled) == 0 &&
		len(blstruct.musthave.configured.services) == 0 &&
		len(blstruct.musthave.users.users) == 0 &&
		blstruct.musthave.policies.polimport == "" &&
		!blstruct.musthave.policies.polreboot &&
		blstruct.musthave.policies.polstatus == "" &&
		len(blstruct.musthave.rules.fwopen.ports) == 0 &&
		len(blstruct.musthave.rules.fwopen.protocols) == 0 &&
		len(blstruct.musthave.rules.fwclosed.ports) == 0 &&
		len(blstruct.musthave.rules.fwclosed.protocols) == 0 &&
		len(blstruct.musthave.rules.fwzones) == 0 &&
		len(blstruct.musthave.mounts.mountname) == 0 {
		commandset[""] = ""
		fmt.Printf("Skipping...\n")
	} else {
		// MH installed
		fmt.Printf("\n")
		fmt.Printf(" Installed: ")
		if len(blstruct.musthave.installed) > 0 {
			fmt.Printf("\n")
			for _, ve := range blstruct.musthave.installed {
				for key, val := range *sshList {
					if commandset[val] == "" {
						// TODO Must Have Installed apply make some changes and move to cmdbuilders
						commandset[key] = serviceCommandBuilder(&ve, &val, "install")
					}
				}
				//for k, v := range commandset {
				//	fmt.Printf("%v   %v\n", k, v)
				//}
				// send to channel
				// wait for response
			}
		} else {
			fmt.Printf("Skipping...\n")
		}

		// MH enabled
		fmt.Printf(" Enabled: ")
		if len(blstruct.musthave.enabled) > 0 {
			commandset = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range blstruct.musthave.enabled {
				for key, val := range *sshList {
					if commandset[val] == "" {
						// TODO Must Have Enabled apply
						commandset[key] = serviceCommandBuilder(&ve, &val, "enable")
					}
				}
				//for k, v := range commandset {
				//	fmt.Printf("%v   %v\n", k, v)
				//}
				// send to channel
				// wait for response
			}
		} else {
			fmt.Printf("Skipping...\n")
		}

		// MH disabled
		fmt.Printf(" Disabled: ")
		if len(blstruct.musthave.disabled) > 0 {
			commandset = make(map[string]string)
			for _, ve := range blstruct.musthave.disabled {
				if ve != "" {
					fmt.Printf("\n")
					for key, val := range *sshList {
						if commandset[val] == "" {
							commandset[key] = serviceCommandBuilder(&ve, &val, "disable")
						}
					}
					// TODO Must Have Disabled apply
					//for k, v := range commandset {
					//	fmt.Printf("%v   %v\n", k, v)
					//}
					// send to channel
					// wait for response
				} else {
					fmt.Printf("Skipping...\n")
				}
			}
		} else {
			fmt.Printf("Skipping...\n")
		}
		// MH configured
		fmt.Printf(" Configured Checklist: ")
		for ke, ve := range blstruct.musthave.configured.services {
			if ke == "" {
				fmt.Printf("Skipping...\n")
			} else {
				commandset = make(map[string]string)
				fmt.Printf("\n      %s:\n", ke)
				if len(ve.source) == len(ve.destination) {
					for i := range ve.source {
						fmt.Println(ve.source[i])
						fmt.Println(ve.destination[i])
						// TODO Transfer config files via ssh
						// pass source and destination to channel, to transfer the file
					}
				} else {
					fmt.Printf("Skipping... Config mismatch in baseline file\n")
				}
			}
		}
		// MH Users
		fmt.Printf(" Users Checklist: ")
		for ke, ve := range blstruct.musthave.users.users {
			if ke == "" {
				fmt.Printf("Skipping...\n")
			} else {
				commandset = make(map[string]string)
				fmt.Printf("\n      %s:\n", ke)
				for key, val := range *sshList {
					if commandset[val] == "" {
						commandset[key] = ve.userManagementCommandBuilder(&ke, "add")
					}
				}
				// TODO User apply
				//for k, v := range commandset {
				//	fmt.Printf("%v   %v\n", k, v)
				//}

				// fmt.Printf("   Groups: ")
				// if len(ve.groups) > 0 {
				// 	for _, val := range ve.groups {
				// 		fmt.Printf("%s\n", val)
				// 	}
				// } else {
				// 	fmt.Printf("\n")
				// }
				// fmt.Printf("   Shell: %v\n", ve.shell)
				// fmt.Printf("   Home: %v\n", ve.home)
				// fmt.Printf("   Sudoer: %v\n", ve.sudoer)
			}
		}
		// MH Policies
		fmt.Printf(" Policies Checklist: ")
		if blstruct.musthave.policies.polstatus == "" &&
			blstruct.musthave.policies.polimport == "" &&
			!blstruct.musthave.policies.polreboot {
			fmt.Printf("Skipping...\n")
		} else {
			commandset = make(map[string]string)
			for key, val := range *sshList {
				if commandset[val] == "" {
					// TODO Must Have Policies apply make some changes and move to cmdbuilders
					commandset[key] = blstruct.musthave.policies.policyCommandBuilder("apply")
				}
			}
			if blstruct.musthave.policies.polreboot {
				*rebootBool = true
			}
			//for k, v := range commandset {
			//	fmt.Printf("%v   %v\n", k, v)
			//}
			// Send command to channel

			fmt.Printf("\n")
		}
		// MH Firewall rules
		fmt.Printf(" Firewall Checklist: ")
		if len(blstruct.musthave.rules.fwopen.ports) == 0 &&
			len(blstruct.musthave.rules.fwopen.protocols) == 0 &&
			len(blstruct.musthave.rules.fwclosed.ports) == 0 &&
			len(blstruct.musthave.rules.fwclosed.protocols) == 0 &&
			len(blstruct.musthave.rules.fwzones) == 0 {
			fmt.Printf("Skipping...\n")
		} else {
			commandset = make(map[string]string)
			fmt.Printf("\n")
			if len(blstruct.musthave.rules.fwopen.ports) == len(blstruct.musthave.rules.fwopen.protocols) {
				if len(blstruct.musthave.rules.fwzones) > 0 {
					fmt.Println("   Firewall zones:")
					for _, ve := range blstruct.musthave.rules.fwzones {
						fmt.Printf("      %v\n", ve)
						for i := range blstruct.musthave.rules.fwopen.ports {
							for key, val := range *sshList {
								if commandset[val] == "" {
									commandset[key] = firewallCommandBuilder(&blstruct.musthave.rules.fwopen.ports[i],
										&blstruct.musthave.rules.fwopen.protocols[i],
										&ve,
										"apply-open")
								}
							}
							//for k, v := range commandset {
							//	// TODO Open Firewall ports & protocols check per firewall zone apply
							//	fmt.Printf("%v   %v\n", k, v)
							//}
							// firewall check creation per zone
							// channel to ssh session and wait for a reply
						}
					}
				} else {
					for i := range blstruct.musthave.rules.fwopen.ports {
						fmt.Printf("%s  %s\n", blstruct.musthave.rules.fwopen.ports[i],
							blstruct.musthave.rules.fwopen.protocols[i])
						for key, val := range *sshList {
							if commandset[val] == "" {
								emptyZone := ""
								commandset[key] = firewallCommandBuilder(&blstruct.musthave.rules.fwopen.ports[i],
									&blstruct.musthave.rules.fwopen.protocols[i],
									&emptyZone,
									"apply-open")
							}
						}
						//for k, v := range commandset {
						//	// TODO Open Firewall ports & protocols apply
						//	fmt.Printf("%v   %v\n", k, v)
						//}
						// firewall check creation with no zone specified
						// channel to ssh session and wait for a reply
					}
				}
			} else {
				fmt.Println("There seems to be inconsistencies between your firewall ports and protocols.")
				fmt.Println("Please review your baseline and rectify it.")
			}
			if len(blstruct.musthave.rules.fwclosed.ports) == len(blstruct.musthave.rules.fwclosed.protocols) {
				if len(blstruct.musthave.rules.fwzones) > 0 {
					fmt.Println("   Firewall zones:")
					for _, ve := range blstruct.musthave.rules.fwzones {
						fmt.Printf("      %v\n", ve)
						for i := range blstruct.musthave.rules.fwclosed.ports {
							for key, val := range *sshList {
								if commandset[val] == "" {
									commandset[key] = firewallCommandBuilder(&blstruct.musthave.rules.fwclosed.ports[i],
										&blstruct.musthave.rules.fwclosed.protocols[i],
										&ve,
										"apply-closed")
								}
							}
							//for k, v := range commandset {
							//	// TODO Closed Firewall ports & protocols check per firewall zone apply
							//	fmt.Printf("%v   %v\n", k, v)
							//}
							// firewall check creation per zone
							// channel to ssh session and wait for a reply
						}
					}
				} else {
					for i := range blstruct.musthave.rules.fwclosed.ports {
						for key, val := range *sshList {
							if commandset[val] == "" {
								emptyZone := ""
								commandset[key] = firewallCommandBuilder(&blstruct.musthave.rules.fwclosed.ports[i],
									&blstruct.musthave.rules.fwclosed.protocols[i],
									&emptyZone,
									"apply-closed")
							}
						}
						//for k, v := range commandset {
						//	// TODO Open Firewall ports & protocols apply
						//	fmt.Printf("%v   %v\n", k, v)
						//}
						// firewall check creation with no zone specified
						// channel to ssh session and wait for a reply
					}
				}
			} else {
				fmt.Println("There seems to be inconsistencies between your firewall ports and protocols.")
				fmt.Println("Please review your baseline and rectify it.")
			}
		}
		// MH mounts
		fmt.Printf(" Mounts Checklist: ")
		for ke, ve := range blstruct.musthave.mounts.mountname {
			if ke == "" {
				fmt.Printf("Skipping...\n")
			} else {
				if ve.mounttype == "" &&
					ve.address == "" &&
					ve.src == "" &&
					ve.dest == "" {
					fmt.Printf("\nNo info found for %s. Skipping...\n", ke)
				} else {
					fmt.Printf("\n")
					commandset = make(map[string]string)
					notEnoughInfo := false
					fmt.Printf("      %s:\n", ke)
					if ve.mounttype == "" {
						notEnoughInfo = true
					}
					if ve.address == "" {
						notEnoughInfo = true
					}
					if ve.src == "" {
						notEnoughInfo = true
					}
					if ve.dest == "" {
						notEnoughInfo = true
					}
					if notEnoughInfo {
						fmt.Printf("Critical mounting info missing for %s. Please review your baseline's mounting information. Skipping...\n", ke)
					} else {
						for key, val := range *sshList {
							if commandset[val] == "" {
								// TODO Must Have Mounts apply
								commandset[key] = ve.mountCommandBuilder("apply")
							}
						}
						// iterate through sshList and create command for each server
						// pass info to ssh session and waiting for a response
					}
					//for k, v := range commandset {
					//	fmt.Printf("%v   %v\n", k, v)
					//}
				}
			}
		}
	}
	return
}