package sshlib

import (
	"fmt"
	"strings"

	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"gopkg.in/yaml.v2"
)

func (blstruct *ParsedBaseline) applyOSExcludes(servergroupname string, configs *yaml.MapSlice) (serverlist []string, oslist []string) {
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
				serverlist = append(serverlist, serverValue[0].Value.(string))
				oslist = append(oslist, serverValue[5].Value.(string))
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
							serverlist = append(serverlist, serverValue[0].Value.(string))
							oslist = append(oslist, serverValue[5].Value.(string))
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
						fmt.Println(serverValue[0].Value)
						if !ok {
							panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
						}
						serverlist = append(serverlist, serverValue[0].Value.(string))
						oslist = append(oslist, serverValue[5].Value.(string))
					}
				}
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
							serverlist = append(serverlist, serverValue[0].Value.(string))
							oslist = append(oslist, serverValue[5].Value.(string))
						}
					}
				}
			}
		}
	}
	return
}

func (blstruct *ParsedBaseline) applyPrereq(servers *[]string, oslist *[]string) (commandset []string) {
	fmt.Println(servers)
	fmt.Println("Verifying server group's prerequisites list:")
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
		// fmt.Println("No prerequisites have been specified  -- Please check your baseline, if you believe this to be incorrect")
	} else {
		// prerequisite tools
		if len(blstruct.prereq.tools) == 0 {
			// fmt.Println("No prerequisite tools specified")
		} else {
			// fmt.Println("The following prerequisite tools will be installed via the package manager:")
			for _, ve := range blstruct.prereq.tools {
				fmt.Println(ve)
			}
		}
		// prerequisite files URLs
		if len(blstruct.prereq.files.urls) == 0 {
			// fmt.Println("No prerequisite files URLs' specified")
		} else {
			// fmt.Println("The following prerequisite files URLs' will be downloaded via curl/wget:")
			for _, ve := range blstruct.prereq.files.urls {
				fmt.Println(ve)
			}
		}
		// prerequisite files local
		if blstruct.prereq.files.local.dest != "" &&
			blstruct.prereq.files.local.src != "" {
			// fmt.Println("No prerequisite files (local) specified")
		} else {
			// fmt.Println("The following files will be transferred locally via scp")
			if blstruct.prereq.files.local.src != "" {
				// fmt.Println("Source (locally):")
				// fmt.Println(blstruct.prereq.files.local.src)
			} else {
				// fmt.Println("No source/local file or directory specified")
			}
			if blstruct.prereq.files.local.dest != "" {
				// fmt.Println("Destination (remote):")
				// fmt.Println(blstruct.prereq.files.local.dest)
			} else {
				// fmt.Println("No destination/remote paths specified")
			}
			// fmt.Println("Please review your baseline if either of these are empty")
		}
		// prerequisite files remote
		if blstruct.prereq.files.remote.address == "" &&
			blstruct.prereq.files.remote.dest == "" &&
			blstruct.prereq.files.remote.mounttype == "" &&
			blstruct.prereq.files.remote.pwd == "" &&
			blstruct.prereq.files.remote.src == "" &&
			blstruct.prereq.files.remote.username == "" &&
			len(blstruct.prereq.files.remote.files) == 0 {
			// fmt.Println("No prerequisite files (remote) specified")
		} else {
			// fmt.Println("The following files will be transferred via the mount details specified:")
			if blstruct.prereq.files.remote.mounttype == "" {
				// fmt.Println("Mount type:")
				// fmt.Println(blstruct.prereq.files.remote.mounttype)
			} else {
				// fmt.Println("No mount type specified")
			}
			if blstruct.prereq.files.remote.address != "" {
				// fmt.Println("Address:")
				// fmt.Println(blstruct.prereq.files.remote.address)
			} else {
				// fmt.Println("No address specified")
			}
			if blstruct.prereq.files.remote.username != "" {
				// fmt.Println("Username:")
				// fmt.Println(blstruct.prereq.files.remote.username)

			} else {
				// fmt.Println("No username specified")
			}
			if blstruct.prereq.files.remote.pwd != "" {
				// fmt.Println("Password:")
				// fmt.Println(blstruct.prereq.files.remote.pwd)
			} else {
				// fmt.Println("No password specified")
			}
			if blstruct.prereq.files.remote.src != "" {
				// fmt.Println("Source (remote):")
				// fmt.Println(blstruct.prereq.files.remote.src)
			} else {
				// fmt.Println("No mount source specified")
			}
			if blstruct.prereq.files.remote.dest != "" {
				// fmt.Println("Destination (remote):")
				// fmt.Println(blstruct.prereq.files.remote.dest)
			} else {
				// fmt.Println("No mount destination specified")
			}
			if len(blstruct.prereq.files.remote.files) == 0 {
				// fmt.Println("No prerequisite tools specified")
			} else {
				// fmt.Println("Files to be transferred:")
				for _, ve := range blstruct.prereq.files.remote.files {
					fmt.Println(ve)
				}
			}
		}
		// prerequisite VCS instructions
		if len(blstruct.prereq.vcs.execute) == 0 &&
			len(blstruct.prereq.vcs.urls) == 0 {
			// fmt.Println("No VCS information specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			if len(blstruct.prereq.vcs.urls) > 0 {
				// fmt.Println("VCS URL's to be cloned to the home directory:")
				for _, ve := range blstruct.prereq.vcs.urls {
					fmt.Println(ve)
				}
			} else {
				// fmt.Println("No VCS URL's specified")
			}
			if len(blstruct.prereq.vcs.execute) > 0 {
				// fmt.Println("VCS related commands to be executed:")
				for _, ve := range blstruct.prereq.vcs.execute {
					fmt.Println(ve)
				}
			} else {
				// fmt.Println("No VCS commands to execute")
			}
			// fmt.Println("Please review your baseline if either of these are empty")
		}
		// prerequisite script
		if blstruct.prereq.script == "" {
			// fmt.Println("No prerequisite script information specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			// fmt.Println("Script:")
			// fmt.Println(blstruct.prereq.script)
		}
		// prerequisite cleanup
		if !blstruct.prereq.cleanup {
			// fmt.Println("Prerequisite cleanup is set to false.")
		} else {
			// fmt.Println("Prerequisite cleanup is set to true.")
		}

	}
	return
}

func (blstruct *ParsedBaseline) applyMustHaves(servers *[]string, oslist *[]string) (commandset []string) {
	fmt.Println(pkgmanlib.PkgRefresh["arch"]) //  just to prevent go from removing the import
	// MH list
	fmt.Printf("Must Have Checklist: ")
	if len(blstruct.musthave.installed) == 0 && // done
		len(blstruct.musthave.enabled) == 0 && // done
		len(blstruct.musthave.disabled) == 0 && // done
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
		fmt.Printf("Skipping...\n")
		commandset = append(commandset, "")
	} else {
		// MH installed
		fmt.Printf(" Installed: ")
		if len(blstruct.musthave.installed) > 0 {
			fmt.Printf("\n")
			for _, ve := range blstruct.musthave.installed {
				fmt.Println(ve)
			}
		} else {
			fmt.Printf("Skipping...\n")
		}

		// MH enabled
		fmt.Printf(" Enabled: ")
		if len(blstruct.musthave.enabled) > 0 {
			fmt.Printf("\n")
			for _, ve := range blstruct.musthave.enabled {
				fmt.Println(ve)
			}
		} else {
			fmt.Printf("Skipping...\n")
		}

		// MH disabled
		fmt.Printf(" Disabled: ")
		if len(blstruct.musthave.disabled) > 0 {
			fmt.Printf("\n")
			for _, ve := range blstruct.musthave.disabled {
				fmt.Println(ve)
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
				fmt.Printf("\n      %s:\n", ke)
				for _, val := range ve.source {
					fmt.Printf("Baseline File (Source): %s\n", val)
				}
				for _, val := range ve.destination {
					fmt.Printf("Current File (Destination): %s\n", val)
				}
			}
		}
		// MH Users
		fmt.Printf(" Users Checklist: ")
		for ke, ve := range blstruct.musthave.users.users {
			if ke == "" {
				fmt.Printf("Skipping...\n")
			} else {
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
			}
		}
		// MH Policies
		fmt.Printf(" Policies Checklist: ")
		if blstruct.musthave.policies.polstatus == "" &&
			blstruct.musthave.policies.polimport == "" &&
			!blstruct.musthave.policies.polreboot {
			fmt.Printf("Skipping...\n")
		} else {
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
			}
		}
	}
	return commandset
}

func (blstruct *ParsedBaseline) applyMustNotHaves(servers *[]string, oslist *[]string) (commandset []string) {
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
		}
		// MNH mounts
		fmt.Printf(" Mounts Checklist: ")
		if len(blstruct.mustnothave.mounts) > 0 {
			fmt.Printf("\n")
			for _, ve := range blstruct.mustnothave.mounts {
				fmt.Println(ve)
			}
		} else {
			fmt.Printf("Skipping...\n")
		}
	}
	return commandset
}

func (blstruct *ParsedBaseline) applyFinals(servers *[]string, oslist *[]string) (commandset []string) {
	// Final steps list
	// fmt.Println("Verifying server group's final steps list:")
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
		if len(blstruct.final.scripts) > 0 {
			// fmt.Println("Scripts:")
			for _, ve := range blstruct.final.scripts {
				fmt.Println(ve)
			}
		} else {
			// fmt.Println("No scripts specified  -- Please check your baseline, if you believe this to be incorrect")
		}
		// final commands
		if len(blstruct.final.commands) > 0 {
			// fmt.Println("Commands:")
			for _, ve := range blstruct.final.commands {
				fmt.Println(ve)
			}
		} else {
			// fmt.Println("No commands specified  -- Please check your baseline, if you believe this to be incorrect")
		}
		// final collections
		if len(blstruct.final.collect.logs) == 0 &&
			len(blstruct.final.collect.stats) == 0 &&
			len(blstruct.final.collect.files) == 0 &&
			!blstruct.final.collect.users {
			// fmt.Println("No collections specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			// fmt.Println("Collect:")
			if len(blstruct.final.collect.logs) > 0 {
				// fmt.Println("Logs:")
				for _, ve := range blstruct.final.collect.logs {
					fmt.Println(ve)
				}
			} else {
				// fmt.Println("No logs specified  -- Please check your baseline, if you believe this to be incorrect")
			}
			if len(blstruct.final.collect.stats) > 0 {
				// fmt.Println("Stats:")
				for _, ve := range blstruct.final.collect.stats {
					fmt.Println(ve)
				}
			} else {
				// fmt.Println("No stats specified  -- Please check your baseline, if you believe this to be incorrect")
			}
			if len(blstruct.final.collect.files) > 0 {
				// fmt.Println("Files:")
				for _, ve := range blstruct.final.collect.files {
					fmt.Println(ve)
				}
			} else {
				// fmt.Println("No files specified  -- Please check your baseline, if you believe this to be incorrect")
			}
		}
		// final restarts
		if !blstruct.final.restart.services &&
			!blstruct.final.restart.servers {
			// fmt.Println("No restart options specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			// fmt.Println("Reboot:")
			if blstruct.final.restart.services {
				// fmt.Printf("Services: %v\n", blstruct.final.restart.services)
			} else {
				// fmt.Println("Services reboot set to false, or not in baseline  -- Please check your baseline, if you believe this to be incorrect")
			}
			if blstruct.final.restart.servers {
				// fmt.Printf("Servers: %v\n", blstruct.final.restart.servers)
			} else {
				// fmt.Println("Servers reboot set to false, or not in baseline  -- Please check your baseline, if you believe this to be incorrect")
			}
		}
	}
	return
}