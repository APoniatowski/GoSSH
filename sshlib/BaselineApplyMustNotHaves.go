package sshlib

import (
	"fmt"
	"github.com/APoniatowski/GoSSH/pkgmanlib"
)

func (baselineStruct *ParsedBaseline) applyMustNotHaves(sshList *map[string]string, commandChannel chan<- map[string]string) {
	//MNH list
	commandSet := make(map[string]string)
	fmt.Printf("Must Not Have Checklist: ")
	if len(baselineStruct.mustnothave.installed) == 0 &&
		len(baselineStruct.mustnothave.enabled) == 0 &&
		len(baselineStruct.mustnothave.disabled) == 0 &&
		len(baselineStruct.mustnothave.users) == 0 &&
		len(baselineStruct.mustnothave.rules.fwopen.ports) == 0 &&
		len(baselineStruct.mustnothave.rules.fwopen.protocols) == 0 &&
		len(baselineStruct.mustnothave.rules.fwclosed.ports) == 0 &&
		len(baselineStruct.mustnothave.rules.fwclosed.protocols) == 0 &&
		len(baselineStruct.mustnothave.rules.fwzones) == 0 &&
		len(baselineStruct.mustnothave.mounts) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		// MNH installed
		fmt.Printf("\n")
		fmt.Printf(" Installed Checklist: ")
		if len(baselineStruct.mustnothave.installed) > 0 {
			commandSet = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range baselineStruct.mustnothave.installed {
				for key, val := range *sshList {
					if commandSet[val] == "" {
						// TODO Must Not Have Installed apply make some changes and move to commandBuilders
						commandSet[key] = serviceCommandBuilder(&ve, &val, "uninstall")
					}
				}
				commandChannel <- commandSet
				//for k, v := range commandSet {
				//	fmt.Printf("%v   %v\n", k, v)
				//}
				// send to channel
				// wait for response and display compliancy
			}
		} else {
			fmt.Printf("Skipping...\n")
		}
		// MNH enabled
		fmt.Printf(" Enabled Checklist: ")
		if len(baselineStruct.mustnothave.enabled) > 0 {
			commandSet = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range baselineStruct.mustnothave.enabled {
				for key, val := range *sshList {
					if commandSet[val] == "" {
						// TODO Must Not Have Enabled apply make some changes and move to commandBuilders
						commandSet[key] = serviceCommandBuilder(&ve, &val, "disable")
					}
				}
				commandChannel <- commandSet
				//for k, v := range commandSet {
				//	fmt.Printf("%v   %v\n", k, v)
				//}
				// send to channel
				// wait for response and display compliancy
				// check if service is active
			}
		} else {
			fmt.Printf("Skipping...\n")
		}
		// MNH disabled
		fmt.Printf(" Disabled Checklist: ")
		if len(baselineStruct.mustnothave.disabled) > 0 {
			commandSet = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range baselineStruct.mustnothave.disabled {
				if ve != "" {
					for key, val := range *sshList {
						if commandSet[val] == "" {
							// TODO Must Not Have Disabled apply make some changes and move to commandBuilders
							commandSet[key] = serviceCommandBuilder(&ve, &val, "enable")
						}
					}
					commandChannel <- commandSet
					//for k, v := range commandSet {
					//	fmt.Printf("%v   %v\n", k, v)
					//}
					// send to channel
					// wait for response and display compliancy
					// check if service is inactive
				} else {
					fmt.Printf("Skipping...\n")
				}
			}
		}
		// MNH Users
		fmt.Printf(" Users Checklist: ")
		if len(baselineStruct.mustnothave.users) > 0 {
			commandSet = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range baselineStruct.mustnothave.users {
				if ve != "" {
					for key, val := range *sshList {
						if commandSet[val] == "" {
							// TODO Must Not Have Users apply make some changes and move to commandBuilders
							commandSet[key] = pkgmanlib.OmniTools["userdel"] + ve
						}
					}
					commandChannel <- commandSet
					//for k, v := range commandSet {
					//	fmt.Printf("%v   %v\n", k, v)
					//}
					// iterate through sshList and create command for each server
					// pass info to ssh session and waiting for a response
				}
			}
		} else {
			fmt.Printf("Skipping...\n")
		}
		// MNH Firewall rules
		fmt.Printf(" Firewall Checklist: ")
		if len(baselineStruct.mustnothave.rules.fwopen.ports) == 0 &&
			len(baselineStruct.mustnothave.rules.fwopen.protocols) == 0 &&
			len(baselineStruct.mustnothave.rules.fwclosed.ports) == 0 &&
			len(baselineStruct.mustnothave.rules.fwclosed.protocols) == 0 &&
			len(baselineStruct.mustnothave.rules.fwzones) == 0 {
			fmt.Printf("Skipping...\n")
		} else {
			commandSet = make(map[string]string)
			fmt.Printf("\n")
			if len(baselineStruct.mustnothave.rules.fwopen.ports) == len(baselineStruct.mustnothave.rules.fwopen.protocols) {
				if len(baselineStruct.mustnothave.rules.fwzones) > 0 {
					fmt.Println("   Firewall zones:")
					for _, ve := range baselineStruct.mustnothave.rules.fwzones {
						fmt.Printf("      %v\n", ve)
						for i := range baselineStruct.mustnothave.rules.fwopen.ports {
							for key, val := range *sshList {
								if commandSet[val] == "" {
									commandSet[key] = firewallCommandBuilder(&baselineStruct.mustnothave.rules.fwopen.ports[i],
										&baselineStruct.mustnothave.rules.fwopen.protocols[i],
										&ve,
										"remove-open")
								}
							}
							commandChannel <- commandSet
							//for k, v := range commandSet {
							//	// TODO No Open Firewall ports & protocols check per firewall zone apply
							//	fmt.Printf("%v   %v\n", k, v)
							//}
							// firewall check creation per zone
							// channel to ssh session and wait for a reply
						}
					}
				} else {
					for i := range baselineStruct.mustnothave.rules.fwopen.ports {
						fmt.Printf("%s  %s\n", baselineStruct.mustnothave.rules.fwopen.ports[i],
							baselineStruct.mustnothave.rules.fwopen.protocols[i])
						for key, val := range *sshList {
							if commandSet[val] == "" {
								emptyZone := ""
								commandSet[key] = firewallCommandBuilder(&baselineStruct.mustnothave.rules.fwopen.ports[i],
									&baselineStruct.mustnothave.rules.fwopen.protocols[i],
									&emptyZone,
									"remove-open")
							}
						}
						commandChannel <- commandSet
						//for k, v := range commandSet {
						//	// TODO No Open Firewall ports & protocols apply
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
			if len(baselineStruct.mustnothave.rules.fwclosed.ports) == len(baselineStruct.mustnothave.rules.fwclosed.protocols) {
				if len(baselineStruct.mustnothave.rules.fwzones) > 0 {
					fmt.Println("   Firewall zones:")
					for _, ve := range baselineStruct.mustnothave.rules.fwzones {
						fmt.Printf("      %v\n", ve)
						for i := range baselineStruct.mustnothave.rules.fwclosed.ports {
							for key, val := range *sshList {
								if commandSet[val] == "" {
									commandSet[key] = firewallCommandBuilder(&baselineStruct.mustnothave.rules.fwclosed.ports[i],
										&baselineStruct.mustnothave.rules.fwclosed.protocols[i],
										&ve,
										"remove-closed")
								}
							}
							commandChannel <- commandSet
							//for k, v := range commandSet {
							//	// TODO No Closed Firewall ports & protocols check per firewall zone apply
							//	fmt.Printf("%v   %v\n", k, v)
							//}
							// firewall check creation per zone
							// channel to ssh session and wait for a reply
						}
					}
				} else {
					for i := range baselineStruct.mustnothave.rules.fwclosed.ports {
						for key, val := range *sshList {
							if commandSet[val] == "" {
								emptyZone := ""
								commandSet[key] = firewallCommandBuilder(&baselineStruct.mustnothave.rules.fwclosed.ports[i],
									&baselineStruct.mustnothave.rules.fwclosed.protocols[i],
									&emptyZone,
									"remove-closed")
							}
						}
						commandChannel <- commandSet
						//for k, v := range commandSet {
						//	// TODO No Open Firewall ports & protocols apply
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
		// MNH mounts
		fmt.Printf(" Mounts Checklist: ")
		if len(baselineStruct.mustnothave.mounts) > 0 {
			commandSet = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range baselineStruct.mustnothave.mounts {
				for key, val := range *sshList {
					if commandSet[val] == "" {
						commandSet[key] = "grep '" + ve + "' /etc/fstab"
					}
				}
				commandChannel <- commandSet
				//    check if the mount address is in fstab
				// iterate through sshList and create command for each server
				// pass info to ssh session and waiting for a response
				for key, val := range *sshList {
					if commandSet[val] == "" {
						// TODO Must Not Have Mounts apply
						commandSet[key] = " grep '" + ve + "' /etc/fstab | awk -F: '{ print $1 }' | " // need to delete the line
					}
				}
				commandChannel <- commandSet
				//   grep the mount address
				// iterate through sshList and create command for each server
				// pass info to ssh session and waiting for a response
			}
		} else {
			fmt.Printf("Skipping...\n")
		}
	}
	return
}
