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

func (blstruct *ParsedBaseline) checkPrereqs(sshList *map[string]string) {
	commandset := make(map[string]string)
	if !blstruct.prereq.cleanup {
		fmt.Printf("Prerequisites Checklist: ")
		if len(blstruct.prereq.vcs.execute) == 0 &&
			len(blstruct.prereq.vcs.urls) == 0 &&
			blstruct.prereq.files.local.dest == "" &&
			blstruct.prereq.files.local.src == "" &&
			blstruct.prereq.files.remote.address == "" &&
			blstruct.prereq.files.remote.dest == "" &&
			blstruct.prereq.files.remote.mounttype == "" &&
			blstruct.prereq.files.remote.pwd == "" &&
			blstruct.prereq.files.remote.src == "" &&
			blstruct.prereq.files.remote.username == "" &&
			len(blstruct.prereq.files.remote.files) == 0 &&
			len(blstruct.prereq.files.urls) == 0 &&
			len(blstruct.prereq.tools) == 0 &&
			blstruct.prereq.script == "" &&
			!blstruct.prereq.cleanup {
			commandset[""] = ""
			fmt.Printf("Skipping...\n")
		} else {
			fmt.Printf("\n")
			// prerequisite tools
			fmt.Printf(" Prerequisite Tools: ")
			if len(blstruct.prereq.tools) == 0 {
				fmt.Printf("Skipping...\n")
			} else {
				fmt.Printf("\n")
				for _, ve := range blstruct.prereq.tools {
					fmt.Printf(ve)
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
			}
			// prerequisite files URLs
			fmt.Printf(" Prerequisite URL's: ")
			if len(blstruct.prereq.files.urls) == 0 {
				fmt.Printf("Skipping...\n")
			} else {
				fmt.Printf("\n")
				for _, ve := range blstruct.prereq.files.urls {
					fmt.Printf(ve)
					parseFile := strings.Split(ve, "/")
					parsedFile := parseFile[len(parseFile)-1]
					for key, val := range *sshList {
						if commandset[val] == "" {
							commandset[key] = pkgmanlib.OmniTools["statinfo"] + parsedFile
						}
					}
					for k, v := range commandset {
						fmt.Printf("%v   %v\n", k, v)
					}
					// send to channel
					// wait for response and display compliancy
				}
			}
			// prerequisite files local
			fmt.Printf(" Prerequisite Files (via scp): ")
			if blstruct.prereq.files.local.dest != "" &&
				blstruct.prereq.files.local.src != "" {
				fmt.Printf("Skipping...\n")
			} else {
				fmt.Printf("\n")
				var srcFile string
				fmt.Println("The following files will be transferred locally via scp")
				if blstruct.prereq.files.local.src != "" {
					srcFile = blstruct.prereq.files.local.src
				}
				if blstruct.prereq.files.local.dest != "" {
					for key, val := range *sshList {
						if commandset[val] == "" {
							commandset[key] = pkgmanlib.OmniTools["suminfo"] + blstruct.prereq.files.local.dest + srcFile
							/*
								-will need to find a better way to compare files and directories-
								cat would kill memory, if its a large file or binary
								sum only does files, not dirs
								need to create for loop command if its a directory with md5sum
							*/
						}
					}
					for k, v := range commandset {
						fmt.Printf("%v   %v\n", k, v)
					}
					// send to channel
					// wait for response
					// diff the file/dir with the source
					// display compliancy
				}
			}
			// prerequisite files remote
			fmt.Printf(" Prerequisite Files (via mount): ")
			if blstruct.prereq.files.remote.address == "" &&
				blstruct.prereq.files.remote.dest == "" &&
				blstruct.prereq.files.remote.mounttype == "" &&
				blstruct.prereq.files.remote.pwd == "" &&
				blstruct.prereq.files.remote.src == "" &&
				blstruct.prereq.files.remote.username == "" &&
				len(blstruct.prereq.files.remote.files) == 0 {
				fmt.Printf("Skipping...\n")
			} else {
				fmt.Printf("\n")
				if blstruct.prereq.files.remote.dest != "" {
					fmt.Println(blstruct.prereq.files.remote.dest)
				}
				if len(blstruct.prereq.files.remote.files) != 0 {
					for _, ve := range blstruct.prereq.files.remote.files {
						fmt.Println(ve)
						for key, val := range *sshList {
							if commandset[val] == "" {
								commandset[key] = pkgmanlib.OmniTools["suminfo"] + blstruct.prereq.files.local.dest
								/*
									-will need to find a better way to compare files and directories-
									cat would kill memory, if its a large file or binary
									sum only does files, not dirs
									need to create for loop command if its a directory with md5sum
								*/
							}
						}
						for k, v := range commandset {
							fmt.Printf("%v   %v\n", k, v)
						}
						// send to channel
						// wait for response
						// diff the file/dir with the source
						// display compliancy
					}
				}
			}
			// prerequisite VCS instructions
			fmt.Printf(" Prerequisite Files (via VCS): ")
			if len(blstruct.prereq.vcs.execute) == 0 &&
				len(blstruct.prereq.vcs.urls) == 0 {
				fmt.Printf("Skipping...\n")
			} else {
				fmt.Printf("\n")
				if len(blstruct.prereq.vcs.urls) > 0 {
					fmt.Println("VCS URL's to be cloned to the home directory:")
					var vcsDirs string
					for _, ve := range blstruct.prereq.vcs.urls {
						fmt.Println(ve)
						parseFile := strings.Split(ve, "/")
						parsedFile := parseFile[len(parseFile)-1]
						vcsDirs = vcsDirs + ve
						for key, val := range *sshList {
							if commandset[val] == "" {
								commandset[key] = pkgmanlib.OmniTools["statinfo"] + parsedFile
								/*
									-will need to find a better way to compare files and directories-
									ls the dir and check if it exists?
									or use stat?
									add home dir path?
								*/
							}
						}
						for k, v := range commandset {
							fmt.Printf("%v   %v\n", k, v)
						}
						// send to channel
						// wait for response
						// diff the file/dir with the source
						// display compliancy
					}
				}
			}
		}
	} else {
		commandset[""] = ""
	}
	return
}

func (blstruct *ParsedBaseline) checkMustHaves(sshList *map[string]string)  {
	commandset := make(map[string]string)
	fmt.Println(pkgmanlib.PkgSearch["arch"]) //  just to prevent go from removing the import
	fmt.Println("line 307  !!! remember to remove later !!!")
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
						commandset[key] = pkgmanlib.OmniTools["serviceisactive"] + ve
					}
				}
				for k, v := range commandset {
					fmt.Printf("%v   %v\n", k, v)
				}
				// send to channel
				// wait for response and display compliancy
				// check if service is active
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
							commandset[key] = pkgmanlib.OmniTools["serviceisactive"] + ve
						}
					}
					for k, v := range commandset {
						fmt.Printf("%v   %v\n", k, v)
					}
					// send to channel
					// wait for response and display compliancy
					// check if service is inactive
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
				var infoAvailable bool
				fmt.Printf("\n      %s:\n", ke)
				if len(ve.groups) == 0 &&
					ve.home == "" &&
					ve.shell == "" &&
					!ve.sudoer {
					fmt.Printf("\n") // Here it will only check if the user exists
					for key, val := range *sshList {
						if commandset[val] == "" {
							commandset[key] = pkgmanlib.OmniTools["userinfo"] + ke
						}
					}
					for k, v := range commandset {
						fmt.Printf("%v   %v\n", k, v)
					}
					infoAvailable = false
				} else {
					for key, val := range *sshList {
						if commandset[val] == "" {
							commandset[key] = pkgmanlib.OmniTools["userinfo"] + ke
						}
					}
					infoAvailable = true
					for k, v := range commandset {
						fmt.Printf("%v   %v\n", k, v)
					}
				}
				// iterate through sshList and create command for each server
				// pass info to ssh session and waiting for a response
				// process the info received available info

				fmt.Println(infoAvailable)
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
			fmt.Printf("\n")
			if blstruct.musthave.policies.polstatus != "" {
				fmt.Printf("   Status: %s\n", blstruct.musthave.policies.polstatus)
				for key, val := range *sshList {
					if commandset[val] == "" {
						commandset[key] = pkgmanlib.OmniTools["userinfo"]
					}
				}
				// iterate through sshList and create command for each server
				// pass info to ssh session and waiting for a response
			}
			if blstruct.musthave.policies.polimport != "" {

				fmt.Printf("   Import: %s\n", blstruct.musthave.policies.polimport)
				for key, val := range *sshList {
					if commandset[val] == "" {
						commandset[key] = pkgmanlib.OmniTools["userinfo"]
					}
				}
				// iterate through sshList and create command for each server
				// pass info to ssh session and waiting for a response
			}
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
						fmt.Printf("      %v\n",ve)
						for i := range blstruct.musthave.rules.fwopen.ports {
							fmt.Printf("%s  %s\n",blstruct.musthave.rules.fwopen.ports[i],
								blstruct.musthave.rules.fwopen.protocols[i])
							// firewall check creation per zone
							// channel to ssh session and wait for a reply
						}
					}
				} else {
					for i := range blstruct.musthave.rules.fwopen.ports {
						fmt.Printf("%s  %s\n",blstruct.musthave.rules.fwopen.ports[i],
							blstruct.musthave.rules.fwopen.protocols[i])
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
						fmt.Printf("      %v\n",ve)
						for i := range blstruct.musthave.rules.fwclosed.ports {
							fmt.Printf("%s  %s\n",blstruct.musthave.rules.fwclosed.ports[i],
								blstruct.musthave.rules.fwclosed.protocols[i])
							// firewall check creation per zone
							// channel to ssh session and wait for a reply
						}
					}
				} else {
					for i := range blstruct.musthave.rules.fwclosed.ports {
						fmt.Printf("%s  %s\n",blstruct.musthave.rules.fwclosed.ports[i],
							blstruct.musthave.rules.fwclosed.protocols[i])
						// firewall check creation with no zone specified
						// channel to ssh session and wait for a reply
					}
				}
			} else {
				fmt.Println("There seems to be inconsistencies between your firewall ports and protocols.")
				fmt.Println("Please review your baseline and rectify it.")
			}
			for key, val := range *sshList {
				if commandset[val] == "" {
					commandset[key] = pkgmanlib.OmniTools["userinfo"]
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
	return
}

func (blstruct *ParsedBaseline) checkMustNotHaves(sshList *map[string]string) {
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
		// MNH enabled
		fmt.Printf(" Enabled Checklist: ")
		if len(blstruct.mustnothave.enabled) > 0 {
			commandset = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.enabled {
				fmt.Println(ve)
				for key, val := range *sshList {
					if commandset[val] == "" {
						commandset[key] = pkgmanlib.OmniTools["serviceisactive"] + ve
					}
				}
				for k, v := range commandset {
					fmt.Printf("%v   %v\n", k, v)
				}
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
					fmt.Printf("\n")
					fmt.Println(ve)
					for key, val := range *sshList {
						if commandset[val] == "" {
							commandset[key] = pkgmanlib.OmniTools["serviceisactive"] + ve
						}
					}
					for k, v := range commandset {
						fmt.Printf("%v   %v\n", k, v)
					}
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
				if len(ve) == 0 {
					fmt.Printf("\n") // Here it will only check if the user exists
					for key, val := range *sshList {
						if commandset[val] == "" {
							commandset[key] = pkgmanlib.OmniTools["userinfo"]
						}
					}
					for k, v := range commandset {
						fmt.Printf("%v   %v\n", k, v)
					}
				} else {
					for key, val := range *sshList {
						if commandset[val] == "" {
							commandset[key] = pkgmanlib.OmniTools["userinfo"]
						}
					}
					for k, v := range commandset {
						fmt.Printf("%v   %v\n", k, v)
					}
				}
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
			commandset = make(map[string]string)

			if len(blstruct.mustnothave.rules.fwopen.ports) == len(blstruct.mustnothave.rules.fwopen.protocols) {
				if len(blstruct.mustnothave.rules.fwzones) > 0 {
					fmt.Println("   Firewall zones:")
					for _, ve := range blstruct.mustnothave.rules.fwzones {
						fmt.Printf("      %v\n",ve)
						for i := range blstruct.mustnothave.rules.fwopen.ports {
							fmt.Printf("%s  %s\n",blstruct.mustnothave.rules.fwopen.ports[i],
								blstruct.mustnothave.rules.fwopen.protocols[i])
							// firewall check creation per zone
							// channel to ssh session and wait for a reply
						}
					}
				} else {
					for i := range blstruct.mustnothave.rules.fwopen.ports {
						fmt.Printf("%s  %s\n",blstruct.mustnothave.rules.fwopen.ports[i],
							blstruct.mustnothave.rules.fwopen.protocols[i])
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
						fmt.Printf("      %v\n",ve)
						for i := range blstruct.mustnothave.rules.fwclosed.ports {
							fmt.Printf("%s  %s\n",blstruct.mustnothave.rules.fwclosed.ports[i],
								blstruct.mustnothave.rules.fwclosed.protocols[i])
							// firewall check creation per zone
							// channel to ssh session and wait for a reply
						}
					}
				} else {
					for i := range blstruct.mustnothave.rules.fwclosed.ports {
						fmt.Printf("%s  %s\n",blstruct.mustnothave.rules.fwclosed.ports[i],
							blstruct.mustnothave.rules.fwclosed.protocols[i])
						// firewall check creation with no zone specified
						// channel to ssh session and wait for a reply
					}
				}
			} else {
				fmt.Println("There seems to be inconsistencies between your firewall ports and protocols.")
				fmt.Println("Please review your baseline and rectify it.")
			}
			for key, val := range *sshList {
				if commandset[val] == "" {
					commandset[key] = pkgmanlib.OmniTools["userinfo"]
				}
			}
			// iterate through sshList and create command for each server
			// pass info to ssh session and waiting for a response
		}
		// MNH mounts
		fmt.Printf(" Mounts Checklist: ")
		if len(blstruct.mustnothave.mounts) > 0 {
			commandset = make(map[string]string)
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.mounts {
				fmt.Println(ve)
				for key, val := range *sshList {
					if commandset[val] == "" {
						commandset[key] = pkgmanlib.OmniTools["userinfo"]
					}
				}
			}
		} else {
			fmt.Printf("Skipping...\n")
		}
	}
	return
}
