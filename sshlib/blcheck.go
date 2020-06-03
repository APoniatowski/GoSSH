package sshlib

import (
	"fmt"
	"strings"

	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"gopkg.in/yaml.v2"
)

func (blstruct *ParsedBaseline) checkOSExcludes(servergroupname string, configs *yaml.MapSlice) map[string]string {
	sshList := make(map[string]string)
	if strings.ToLower(servergroupname) == "all" {
		if len(blstruct.exclude.osExcl) == 0 &&
			len(blstruct.exclude.serversExcl) == 0 {
			var allServers yaml.MapSlice
			// Concatenates the groups to create a single group
			for _, groupItem := range *configs {
				groupValue, ok := groupItem.Value.(yaml.MapSlice)
				if !ok {
					panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
				}
				allServers = append(allServers, groupValue...)
			}
			for _, serverItem := range allServers {
				serverValue, ok := serverItem.Value.(yaml.MapSlice)
				if !ok {
					panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
				}
				sshList[serverValue[0].Value.(string)] = serverValue[5].Value.(string)
			}
		} else {
			for _, groupItem := range *configs {
				groupValue, ok := groupItem.Value.(yaml.MapSlice)
				if !ok {
					panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
				}
				if groupItem.Key == servergroupname {
					for _, serverItem := range groupValue {
						var osnamecheck bool
						var servernamecheck bool
						serverValue, ok := serverItem.Value.(yaml.MapSlice)
						if !ok {
							panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
						}
						if len(blstruct.exclude.osExcl) > 0 {
							for _, ve := range blstruct.exclude.osExcl {
								if strings.EqualFold(serverValue[5].Value.(string), ve) {
									osnamecheck = true
								}
							}
						}
						if len(blstruct.exclude.serversExcl) > 0 {
							for _, ve := range blstruct.exclude.serversExcl {
								if strings.EqualFold(serverValue[0].Value.(string), ve) {
									servernamecheck = true
								}
							}
						}
						if !servernamecheck && !osnamecheck {
							sshList[serverValue[0].Value.(string)] = serverValue[5].Value.(string)
						}
					}
				}
			}
		}
	} else {
		if len(blstruct.exclude.osExcl) == 0 &&
			len(blstruct.exclude.serversExcl) == 0 {
			for _, groupItem := range *configs {
				groupValue, ok := groupItem.Value.(yaml.MapSlice)
				if !ok {
					panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
				}
				if strings.EqualFold(groupItem.Key.(string), servergroupname) {
					for _, serverItem := range groupValue {
						serverValue, ok := serverItem.Value.(yaml.MapSlice)
						if !ok {
							panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
						}
						sshList[serverValue[0].Value.(string)] = serverValue[5].Value.(string)
					}
				}
			}
		} else {
			for _, groupItem := range *configs {
				groupValue, ok := groupItem.Value.(yaml.MapSlice)
				if !ok {
					panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
				}
				if strings.EqualFold(groupItem.Key.(string), servergroupname) {
					for _, serverItem := range groupValue {
						var osnamecheck bool
						var servernamecheck bool
						serverValue, ok := serverItem.Value.(yaml.MapSlice)
						if !ok {
							panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
						}
						if len(blstruct.exclude.osExcl) > 0 {
							for _, ve := range blstruct.exclude.osExcl {
								if strings.EqualFold(serverValue[5].Value.(string), ve) {
									osnamecheck = true
								}
							}
						}
						if len(blstruct.exclude.serversExcl) > 0 {
							for _, ve := range blstruct.exclude.serversExcl {
								if strings.EqualFold(serverValue[0].Value.(string), ve) {
									servernamecheck = true
								}
							}
						}
						if !servernamecheck && !osnamecheck {
							sshList[serverValue[0].Value.(string)] = serverValue[5].Value.(string)
						}
					}
				}
			}
		}
	}
	return sshList
}

func (blstruct *ParsedBaseline) checkMustHaves(sshList *map[string]string) map[string]string {
	commandset := make(map[string]string)
	fmt.Println(pkgmanlib.PkgSearch["arch"]) //  just to prevent go from removing the import
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
				fmt.Println(ve)
				for key, val := range *sshList {
					if commandset[val] == "" {
						commandset[key] = pkgmanlib.PkgSearch[val] + ve
					}
				}
				for k, v := range commandset {
					fmt.Printf("%v   %v\n", k, v)
				}
				// send to channel
				// wait for response and display compliancy
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
				fmt.Println(ve)
				for key, val := range *sshList {
					if commandset[val] == "" {
						commandset[key] = pkgmanlib.PkgSearch[val] + ve
					}
				}
				for k, v := range commandset {
					fmt.Printf("%v   %v\n", k, v)
				}
				// send to channel
				// wait for response and display compliancy
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
					fmt.Println(ve)
					for key, val := range *sshList {
						if commandset[val] == "" {
							commandset[key] = pkgmanlib.PkgSearch[val] + ve
						}
					}
					for k, v := range commandset {
						fmt.Printf("%v   %v\n", k, v)
					}
					// send to channel
					// wait for response and display compliancy
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
				for _, val := range ve.source {
					fmt.Printf("Baseline File (Source): %s\n", val)
				}
				for _, val := range ve.destination {
					fmt.Printf("Current File (Destination): %s\n", val)
				}
				// iterate through sshList and create command for each server
				// pass info to ssh session and waiting for a response
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
				if len(ve.groups) == 0 &&
					ve.home == "" &&
					ve.shell == "" &&
					!ve.sudoer {
					fmt.Printf("\n") // Here it will only check if the user exists
				} else {
					fmt.Printf("   Groups: ")
					if len(ve.groups) > 0 {
						for _, val := range ve.groups {
							fmt.Printf("%s\n", val)
						}
					} else {
						fmt.Printf("\n")
					}
					fmt.Printf("   Shell: %v\n", ve.shell)
					fmt.Printf("   Home: %v\n", ve.home)
					fmt.Printf("   Sudoer: %v\n", ve.sudoer)
				}
				// iterate through sshList and create command for each server
				// pass info to ssh session and waiting for a response
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
			fmt.Printf("\n")
			if blstruct.musthave.policies.polstatus != "" {
				fmt.Printf("   Status: %s\n", blstruct.musthave.policies.polstatus)
			}
			if blstruct.musthave.policies.polimport != "" {
				fmt.Printf("   Import: %s\n", blstruct.musthave.policies.polimport)
			}
			if blstruct.musthave.policies.polreboot {
				fmt.Printf("   Reboot: %v\n", blstruct.musthave.policies.polreboot)
			}
			// iterate through sshList and create command for each server
			// pass info to ssh session and waiting for a response
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
			if len(blstruct.musthave.rules.fwopen.ports) > 0 {
				fmt.Println("   Open ports:")
				for _, ve := range blstruct.musthave.rules.fwopen.ports {
					fmt.Println(ve)
				}
			}
			if len(blstruct.musthave.rules.fwopen.protocols) > 0 {
				fmt.Println("   Open protocols:")
				for _, ve := range blstruct.musthave.rules.fwopen.protocols {
					fmt.Println(ve)
				}
			}
			if len(blstruct.musthave.rules.fwclosed.ports) > 0 {
				fmt.Println("   Closed ports:")
				for _, ve := range blstruct.musthave.rules.fwclosed.ports {
					fmt.Println(ve)
				}
			}
			if len(blstruct.musthave.rules.fwclosed.protocols) > 0 {
				fmt.Println("   Closed protocols:")
				for _, ve := range blstruct.musthave.rules.fwclosed.protocols {
					fmt.Println(ve)
				}
			}
			if len(blstruct.musthave.rules.fwzones) > 0 {
				fmt.Println("   Firewall zones:")
				for _, ve := range blstruct.musthave.rules.fwzones {
					fmt.Println(ve)
				}
			}
			// iterate through sshList and create command for each server
			// pass info to ssh session and waiting for a response
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
					commandset = make(map[string]string)
					fmt.Printf("\n      %s:\n", ke)
					fmt.Printf("   Mount Type: %v\n", ve.mounttype)
					fmt.Printf("   Address: %v\n", ve.address)
					fmt.Printf("   Source: %v\n", ve.src)
					if ve.dest == "" {
						fmt.Printf("Mount directory info not found for %s. Skipping...\n", ke)
					} else {
						fmt.Printf("   Destination: %v\n", ve.dest)
					}
				}
				// iterate through sshList and create command for each server
				// pass info to ssh session and waiting for a response
			}
		}
	}
	return commandset
}

func (blstruct *ParsedBaseline) checkMustNotHaves(sshList *map[string]string) (commandset map[string]string) {
	//MNH list
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
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.installed {
				fmt.Println(ve)
			}
			// iterate through sshList and create command for each server
			// pass info to ssh session and waiting for a response
		} else {
			fmt.Printf("Skipping...\n")
		}
		// MNH enabled
		fmt.Printf(" Enabled Checklist: ")
		if len(blstruct.mustnothave.enabled) > 0 {
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.enabled {
				fmt.Println(ve)
			}
			// iterate through sshList and create command for each server
			// pass info to ssh session and waiting for a response
		} else {
			fmt.Printf("Skipping...\n")
		}
		// MNH disabled
		fmt.Printf(" Disabled Checklist: ")
		if len(blstruct.mustnothave.disabled) > 0 {
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.disabled {
				fmt.Println(ve)
			}
			// iterate through sshList and create command for each server
			// pass info to ssh session and waiting for a response
		} else {
			fmt.Printf("Skipping...\n")
		}
		// MNH Users
		fmt.Printf(" Users Checklist: ")
		if len(blstruct.mustnothave.users) > 0 {
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.users {
				fmt.Println(ve)
			}
			// iterate through sshList and create command for each server
			// pass info to ssh session and waiting for a response
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
			fmt.Printf("\n")
			if len(blstruct.mustnothave.rules.fwopen.ports) > 0 {
				fmt.Println("   Open ports:")
				for _, ve := range blstruct.musthave.rules.fwopen.ports {
					fmt.Println(ve)
				}
			}
			if len(blstruct.mustnothave.rules.fwopen.protocols) > 0 {
				fmt.Println("   Open protocols:")
				for _, ve := range blstruct.mustnothave.rules.fwopen.protocols {
					fmt.Println(ve)
				}
			}
			if len(blstruct.mustnothave.rules.fwclosed.ports) > 0 {
				fmt.Println("   Closed ports:")
				for _, ve := range blstruct.mustnothave.rules.fwclosed.ports {
					fmt.Println(ve)
				}
			}
			if len(blstruct.mustnothave.rules.fwclosed.protocols) > 0 {
				fmt.Println("   Closed protocols:")
				for _, ve := range blstruct.mustnothave.rules.fwclosed.protocols {
					fmt.Println(ve)
				}
			}
			if len(blstruct.mustnothave.rules.fwzones) > 0 {
				fmt.Println("   Firewall zones:")
				for _, ve := range blstruct.mustnothave.rules.fwzones {
					fmt.Println(ve)
				}
			}
			// iterate through sshList and create command for each server
			// pass info to ssh session and waiting for a response
		}
		// MNH mounts
		fmt.Printf(" Mounts Checklist: ")
		if len(blstruct.mustnothave.mounts) > 0 {
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.mounts {
				fmt.Println(ve)
				// commandset[] =  "" + ve + commandset[]
			}
		} else {
			fmt.Printf("Skipping...\n")
		}
	}
	return
}