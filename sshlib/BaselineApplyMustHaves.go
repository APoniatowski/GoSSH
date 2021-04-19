package sshlib

import "fmt"

func (baselineStruct *ParsedBaseline) applyMustHaves(sshList *map[string]string, rebootBool *bool, commandChannel chan<- map[string]string) {
	commandSet := make(map[string]string)
	// MH list
	fmt.Printf("Must Have Checklist: ")
	if len(baselineStruct.musthave.installed) == 0 &&
		len(baselineStruct.musthave.enabled) == 0 &&
		len(baselineStruct.musthave.disabled) == 0 &&
		len(baselineStruct.musthave.configured.services) == 0 &&
		len(baselineStruct.musthave.users.users) == 0 &&
		baselineStruct.musthave.policies.polimport == "" &&
		!baselineStruct.musthave.policies.polreboot &&
		baselineStruct.musthave.policies.polstatus == "" &&
		len(baselineStruct.musthave.rules.fwopen.ports) == 0 &&
		len(baselineStruct.musthave.rules.fwopen.protocols) == 0 &&
		len(baselineStruct.musthave.rules.fwclosed.ports) == 0 &&
		len(baselineStruct.musthave.rules.fwclosed.protocols) == 0 &&
		len(baselineStruct.musthave.rules.fwzones) == 0 &&
		len(baselineStruct.musthave.mounts.mountname) == 0 {
		commandSet[""] = ""
		fmt.Printf("Skipping...\n")
	} else {
		// MH installed
		fmt.Printf("\n")
		fmt.Printf(" Installed: ")
		if len(baselineStruct.musthave.installed) > 0 {
			fmt.Printf("\n")
			for _, ve := range baselineStruct.musthave.installed {
				for key, val := range *sshList {
					if commandSet[val] == "" {
						// TODO Must Have Installed apply make some changes and move to commandBuilders
						commandSet[key] = serviceCommandBuilder(&ve, &val, "install")
					}
				}
				commandChannel <- commandSet
				//for k, v := range commandSet {
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
		if len(baselineStruct.musthave.enabled) > 0 {
			commandSet = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range baselineStruct.musthave.enabled {
				for key, val := range *sshList {
					if commandSet[val] == "" {
						// TODO Must Have Enabled apply
						commandSet[key] = serviceCommandBuilder(&ve, &val, "enable")
					}
				}
				commandChannel <- commandSet
				//for k, v := range commandSet {
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
		if len(baselineStruct.musthave.disabled) > 0 {
			commandSet = make(map[string]string)
			for _, ve := range baselineStruct.musthave.disabled {
				if ve != "" {
					fmt.Printf("\n")
					for key, val := range *sshList {
						if commandSet[val] == "" {
							commandSet[key] = serviceCommandBuilder(&ve, &val, "disable")
						}
					}
					commandChannel <- commandSet
					// TODO Must Have Disabled apply
					//for k, v := range commandSet {
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
		for ke, ve := range baselineStruct.musthave.configured.services {
			if ke == "" {
				fmt.Printf("Skipping...\n")
			} else {
				commandSet = make(map[string]string)
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
		for ke, ve := range baselineStruct.musthave.users.users {
			if ke == "" {
				fmt.Printf("Skipping...\n")
			} else {
				commandSet = make(map[string]string)
				fmt.Printf("\n      %s:\n", ke)
				for key, val := range *sshList {
					if commandSet[val] == "" {
						commandSet[key] = ve.userManagementCommandBuilder(&ke, "add")
					}
				}
				commandChannel <- commandSet
				// TODO User apply
				//for k, v := range commandSet {
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
		if baselineStruct.musthave.policies.polstatus == "" &&
			baselineStruct.musthave.policies.polimport == "" &&
			!baselineStruct.musthave.policies.polreboot {
			fmt.Printf("Skipping...\n")
		} else {
			commandSet = make(map[string]string)
			for key, val := range *sshList {
				if commandSet[val] == "" {
					// TODO Must Have Policies apply make some changes and move to commandBuilders
					commandSet[key] = baselineStruct.musthave.policies.policyCommandBuilder("apply")
				}
			}
			if baselineStruct.musthave.policies.polreboot {
				*rebootBool = true
			}
			commandChannel <- commandSet
			//for k, v := range commandSet {
			//	fmt.Printf("%v   %v\n", k, v)
			//}
			// Send command to channel
			fmt.Printf("\n")
		}
		// MH Firewall rules
		fmt.Printf(" Firewall Checklist: ")
		if len(baselineStruct.musthave.rules.fwopen.ports) == 0 &&
			len(baselineStruct.musthave.rules.fwopen.protocols) == 0 &&
			len(baselineStruct.musthave.rules.fwclosed.ports) == 0 &&
			len(baselineStruct.musthave.rules.fwclosed.protocols) == 0 &&
			len(baselineStruct.musthave.rules.fwzones) == 0 {
			fmt.Printf("Skipping...\n")
		} else {
			commandSet = make(map[string]string)
			fmt.Printf("\n")
			if len(baselineStruct.musthave.rules.fwopen.ports) == len(baselineStruct.musthave.rules.fwopen.protocols) {
				if len(baselineStruct.musthave.rules.fwzones) > 0 {
					fmt.Println("   Firewall zones:")
					for _, ve := range baselineStruct.musthave.rules.fwzones {
						fmt.Printf("      %v\n", ve)
						for i := range baselineStruct.musthave.rules.fwopen.ports {
							for key, val := range *sshList {
								if commandSet[val] == "" {
									commandSet[key] = firewallCommandBuilder(&baselineStruct.musthave.rules.fwopen.ports[i],
										&baselineStruct.musthave.rules.fwopen.protocols[i],
										&ve,
										"apply-open")
								}
							}
							commandChannel <- commandSet
							//for k, v := range commandSet {
							//	// TODO Open Firewall ports & protocols check per firewall zone apply
							//	fmt.Printf("%v   %v\n", k, v)
							//}
							// firewall check creation per zone
							// channel to ssh session and wait for a reply
						}
					}
				} else {
					for i := range baselineStruct.musthave.rules.fwopen.ports {
						fmt.Printf("%s  %s\n", baselineStruct.musthave.rules.fwopen.ports[i],
							baselineStruct.musthave.rules.fwopen.protocols[i])
						for key, val := range *sshList {
							if commandSet[val] == "" {
								emptyZone := ""
								commandSet[key] = firewallCommandBuilder(&baselineStruct.musthave.rules.fwopen.ports[i],
									&baselineStruct.musthave.rules.fwopen.protocols[i],
									&emptyZone,
									"apply-open")
							}
						}
						commandChannel <- commandSet
						//for k, v := range commandSet {
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
			if len(baselineStruct.musthave.rules.fwclosed.ports) == len(baselineStruct.musthave.rules.fwclosed.protocols) {
				if len(baselineStruct.musthave.rules.fwzones) > 0 {
					fmt.Println("   Firewall zones:")
					for _, ve := range baselineStruct.musthave.rules.fwzones {
						fmt.Printf("      %v\n", ve)
						for i := range baselineStruct.musthave.rules.fwclosed.ports {
							for key, val := range *sshList {
								if commandSet[val] == "" {
									commandSet[key] = firewallCommandBuilder(&baselineStruct.musthave.rules.fwclosed.ports[i],
										&baselineStruct.musthave.rules.fwclosed.protocols[i],
										&ve,
										"apply-closed")
								}
							}
							commandChannel <- commandSet
							//for k, v := range commandSet {
							//	// TODO Closed Firewall ports & protocols check per firewall zone apply
							//	fmt.Printf("%v   %v\n", k, v)
							//}
							// firewall check creation per zone
							// channel to ssh session and wait for a reply
						}
					}
				} else {
					for i := range baselineStruct.musthave.rules.fwclosed.ports {
						for key, val := range *sshList {
							if commandSet[val] == "" {
								emptyZone := ""
								commandSet[key] = firewallCommandBuilder(&baselineStruct.musthave.rules.fwclosed.ports[i],
									&baselineStruct.musthave.rules.fwclosed.protocols[i],
									&emptyZone,
									"apply-closed")
							}
						}
						commandChannel <- commandSet
						//for k, v := range commandSet {
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
		for ke, ve := range baselineStruct.musthave.mounts.mountname {
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
					commandSet = make(map[string]string)
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
							if commandSet[val] == "" {
								// TODO Must Have Mounts apply
								commandSet[key] = ve.mountCommandBuilder("apply")
							}
						}
						// iterate through sshList and create command for each server
						// pass info to ssh session and waiting for a response
						commandChannel <- commandSet
					}
					//for k, v := range commandSet {
					//	fmt.Printf("%v   %v\n", k, v)
					//}
				}
			}
		}
	}
	return
}
