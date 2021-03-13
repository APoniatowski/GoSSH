package sshlib

import (
	"fmt"
	"github.com/APoniatowski/GoSSH/pkgmanlib"
)

func (blstruct *ParsedBaseline) applyMustNotHaves(sshList *map[string]string, commandChannel chan<- map[string]string) {
	//MNH list
	commandset := make(map[string]string)
	fmt.Printf("Must Not Have Checklist: ")
	if len(blstruct.mustnothave.installed) == 0 &&
		len(blstruct.mustnothave.enabled) == 0 &&
		len(blstruct.mustnothave.disabled) == 0 &&
		len(blstruct.mustnothave.users) == 0 &&
		len(blstruct.mustnothave.rules.fwopen.ports) == 0 &&
		len(blstruct.mustnothave.rules.fwopen.protocols) == 0 &&
		len(blstruct.mustnothave.rules.fwclosed.ports) == 0 &&
		len(blstruct.mustnothave.rules.fwclosed.protocols) == 0 &&
		len(blstruct.mustnothave.rules.fwzones) == 0 &&
		len(blstruct.mustnothave.mounts) == 0 {
		fmt.Printf("Skipping...\n")
	} else {
		// MNH installed
		fmt.Printf("\n")
		fmt.Printf(" Installed Checklist: ")
		if len(blstruct.mustnothave.installed) > 0 {
			commandset = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.installed {
				for key, val := range *sshList {
					if commandset[val] == "" {
						// TODO Must Not Have Installed apply make some changes and move to cmdbuilders
						commandset[key] = serviceCommandBuilder(&ve, &val, "uninstall")
					}
				}
				commandChannel <- commandset
				//for k, v := range commandset {
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
		if len(blstruct.mustnothave.enabled) > 0 {
			commandset = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.enabled {
				for key, val := range *sshList {
					if commandset[val] == "" {
						// TODO Must Not Have Enabled apply make some changes and move to cmdbuilders
						commandset[key] = serviceCommandBuilder(&ve, &val, "disable")
					}
				}
				commandChannel <- commandset
				//for k, v := range commandset {
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
		if len(blstruct.mustnothave.disabled) > 0 {
			commandset = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.disabled {
				if ve != "" {
					for key, val := range *sshList {
						if commandset[val] == "" {
							// TODO Must Not Have Disabled apply make some changes and move to cmdbuilders
							commandset[key] = serviceCommandBuilder(&ve, &val, "enable")
						}
					}
					commandChannel <- commandset
					//for k, v := range commandset {
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
		if len(blstruct.mustnothave.users) > 0 {
			commandset = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.users {
				if ve != "" {
					for key, val := range *sshList {
						if commandset[val] == "" {
							// TODO Must Not Have Users apply make some changes and move to cmdbuilders
							commandset[key] = pkgmanlib.OmniTools["userdel"] + ve
						}
					}
					commandChannel <- commandset
					//for k, v := range commandset {
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
		if len(blstruct.mustnothave.rules.fwopen.ports) == 0 &&
			len(blstruct.mustnothave.rules.fwopen.protocols) == 0 &&
			len(blstruct.mustnothave.rules.fwclosed.ports) == 0 &&
			len(blstruct.mustnothave.rules.fwclosed.protocols) == 0 &&
			len(blstruct.mustnothave.rules.fwzones) == 0 {
			fmt.Printf("Skipping...\n")
		} else {
			commandset = make(map[string]string)
			fmt.Printf("\n")
			if len(blstruct.mustnothave.rules.fwopen.ports) == len(blstruct.mustnothave.rules.fwopen.protocols) {
				if len(blstruct.mustnothave.rules.fwzones) > 0 {
					fmt.Println("   Firewall zones:")
					for _, ve := range blstruct.mustnothave.rules.fwzones {
						fmt.Printf("      %v\n", ve)
						for i := range blstruct.mustnothave.rules.fwopen.ports {
							for key, val := range *sshList {
								if commandset[val] == "" {
									commandset[key] = firewallCommandBuilder(&blstruct.mustnothave.rules.fwopen.ports[i],
										&blstruct.mustnothave.rules.fwopen.protocols[i],
										&ve,
										"remove-open")
								}
							}
							commandChannel <- commandset
							//for k, v := range commandset {
							//	// TODO No Open Firewall ports & protocols check per firewall zone apply
							//	fmt.Printf("%v   %v\n", k, v)
							//}
							// firewall check creation per zone
							// channel to ssh session and wait for a reply
						}
					}
				} else {
					for i := range blstruct.mustnothave.rules.fwopen.ports {
						fmt.Printf("%s  %s\n", blstruct.mustnothave.rules.fwopen.ports[i],
							blstruct.mustnothave.rules.fwopen.protocols[i])
						for key, val := range *sshList {
							if commandset[val] == "" {
								emptyZone := ""
								commandset[key] = firewallCommandBuilder(&blstruct.mustnothave.rules.fwopen.ports[i],
									&blstruct.mustnothave.rules.fwopen.protocols[i],
									&emptyZone,
									"remove-open")
							}
						}
						commandChannel <- commandset
						//for k, v := range commandset {
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
			if len(blstruct.mustnothave.rules.fwclosed.ports) == len(blstruct.mustnothave.rules.fwclosed.protocols) {
				if len(blstruct.mustnothave.rules.fwzones) > 0 {
					fmt.Println("   Firewall zones:")
					for _, ve := range blstruct.mustnothave.rules.fwzones {
						fmt.Printf("      %v\n", ve)
						for i := range blstruct.mustnothave.rules.fwclosed.ports {
							for key, val := range *sshList {
								if commandset[val] == "" {
									commandset[key] = firewallCommandBuilder(&blstruct.mustnothave.rules.fwclosed.ports[i],
										&blstruct.mustnothave.rules.fwclosed.protocols[i],
										&ve,
										"remove-closed")
								}
							}
							commandChannel <- commandset
							//for k, v := range commandset {
							//	// TODO No Closed Firewall ports & protocols check per firewall zone apply
							//	fmt.Printf("%v   %v\n", k, v)
							//}
							// firewall check creation per zone
							// channel to ssh session and wait for a reply
						}
					}
				} else {
					for i := range blstruct.mustnothave.rules.fwclosed.ports {
						for key, val := range *sshList {
							if commandset[val] == "" {
								emptyZone := ""
								commandset[key] = firewallCommandBuilder(&blstruct.mustnothave.rules.fwclosed.ports[i],
									&blstruct.mustnothave.rules.fwclosed.protocols[i],
									&emptyZone,
									"remove-closed")
							}
						}
						commandChannel <- commandset
						//for k, v := range commandset {
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
		if len(blstruct.mustnothave.mounts) > 0 {
			commandset = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.mounts {

				for key, val := range *sshList {
					if commandset[val] == "" {
						commandset[key] = "grep '" + ve + "' /etc/fstab"
					}
				}
				commandChannel <- commandset
				//    check if the mount address is in fstab
				// iterate through sshList and create command for each server
				// pass info to ssh session and waiting for a response
				for key, val := range *sshList {
					if commandset[val] == "" {
						// TODO Must Not Have Mounts apply
						commandset[key] = " grep '" + ve + "' /etc/fstab | awk -F: '{ print $1 }' | " // need to delete the line
					}
				}
				commandChannel <- commandset
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
