package sshlib

import "fmt"

func (blstruct *ParsedBaseline) verification(servergroupname string) {
	fmt.Printf("%s has been parsed.\n", servergroupname)
	fmt.Println("Proceeding with baseline verification:")
	fmt.Println("Verifying server group's exclusion list:")
	if len(blstruct.exclude.osExcl) == 0 &&
		len(blstruct.exclude.serversExcl) == 0 {
		fmt.Println("No exclusions have been specified or found.")
		fmt.Println("If you have specified an exclusion list, please check your baseline for errors.")
	} else {
		if len(blstruct.exclude.osExcl) > 0 {
			fmt.Println("Operating Systems to be excluded:")
			for _, ve := range blstruct.exclude.osExcl {
				fmt.Println(ve)
			}
		} else {
			fmt.Println("No operating systems were found in the exlusion list.")
			fmt.Println("If there were any specified, it might be due to a syntax error")
			fmt.Println("  Exclude:")
			fmt.Println("    OS:")
			fmt.Println("      -'OS to be exluded here, eg debian'")
		}
		if len(blstruct.exclude.serversExcl) > 0 {
			fmt.Println("Servers to be exluded:")
			for _, ve := range blstruct.exclude.serversExcl {
				fmt.Println(ve)
			}
		} else {
			fmt.Println("No server were found in the exlusion list.")
			fmt.Println("If there were any specified, it might be due to a syntax error")
			fmt.Println("  Exclude:")
			fmt.Println("    Servers:")
			fmt.Println("      -'server hostname to be excluded here'")
			fmt.Println("Please make sure that the hostname is the same as in your pool")
		}
	}
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
		blstruct.prereq.cleanup == false {
		fmt.Println("No prerequisites have been specified  -- Please check your baseline, if you believe this to be incorrect")
	} else {
		// prerequisite tools
		if len(blstruct.prereq.tools) == 0 {
			fmt.Println("No prerequisite tools specified")
		} else {
			fmt.Println("The following prerequisite tools will be installed via the package manager:")
			for _, ve := range blstruct.prereq.tools {
				fmt.Println(ve)
			}
		}
		// prerequisite files URLs
		if len(blstruct.prereq.files.urls) == 0 {
			fmt.Println("No prerequisite files URLs' specified")
		} else {
			fmt.Println("The following prerequisite files URLs' will be downloaded via curl/wget:")
			for _, ve := range blstruct.prereq.files.urls {
				fmt.Println(ve)
			}
		}
		// prerequisite files local
		if blstruct.prereq.files.local.dest != "" &&
			blstruct.prereq.files.local.src != "" {
			fmt.Println("No prerequisite files (local) specified")
		} else {
			fmt.Println("The following files will be transferred locally via scp")
			if blstruct.prereq.files.local.src != "" {
				fmt.Println("Source (locally):")
				fmt.Println(blstruct.prereq.files.local.src)
			} else {
				fmt.Println("No source/local file or directory specified")
			}
			if blstruct.prereq.files.local.dest != "" {
				fmt.Println("Destination (remote):")
				fmt.Println(blstruct.prereq.files.local.dest)
			} else {
				fmt.Println("No destination/remote paths specified")
			}
			fmt.Println("Please review your baseline if either of these are empty")
		}
		// prerequisite files remote
		if blstruct.prereq.files.remote.address == "" &&
			blstruct.prereq.files.remote.dest == "" &&
			blstruct.prereq.files.remote.mounttype == "" &&
			blstruct.prereq.files.remote.pwd == "" &&
			blstruct.prereq.files.remote.src == "" &&
			blstruct.prereq.files.remote.username == "" &&
			len(blstruct.prereq.files.remote.files) == 0 {
			fmt.Println("No prerequisite files (remote) specified")
		} else {
			fmt.Println("The following files will be transferred via the mount details specified:")
			if blstruct.prereq.files.remote.mounttype == "" {
				fmt.Println("Mount type:")
				fmt.Println(blstruct.prereq.files.remote.mounttype)
			} else {
				fmt.Println("No mount type specified")
			}
			if blstruct.prereq.files.remote.address != "" {
				fmt.Println("Address:")
				fmt.Println(blstruct.prereq.files.remote.address)
			} else {
				fmt.Println("No address specified")
			}
			if blstruct.prereq.files.remote.username != "" {
				fmt.Println("Username:")
				fmt.Println(blstruct.prereq.files.remote.username)

			} else {
				fmt.Println("No username specified")
			}
			if blstruct.prereq.files.remote.pwd != "" {
				fmt.Println("Password:")
				fmt.Println(blstruct.prereq.files.remote.pwd)
			} else {
				fmt.Println("No password specified")
			}
			if blstruct.prereq.files.remote.src != "" {
				fmt.Println("Source (remote):")
				fmt.Println(blstruct.prereq.files.remote.src)
			} else {
				fmt.Println("No mount source specified")
			}
			if blstruct.prereq.files.remote.dest != "" {
				fmt.Println("Destination (remote):")
				fmt.Println(blstruct.prereq.files.remote.dest)
			} else {
				fmt.Println("No mount destination specified")
			}
			if len(blstruct.prereq.files.remote.files) == 0 {
				fmt.Println("No prerequisite tools specified")
			} else {
				fmt.Println("Files to be transferred:")
				for _, ve := range blstruct.prereq.files.remote.files {
					fmt.Println(ve)
				}
			}
		}
		// prerequisite VCS instructions
		if len(blstruct.prereq.vcs.execute) == 0 &&
			len(blstruct.prereq.vcs.urls) == 0 {
			fmt.Println("No VCS information specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			if len(blstruct.prereq.vcs.urls) > 0 {
				fmt.Println("VCS URL's to be cloned to the home directory:")
				for _, ve := range blstruct.prereq.vcs.urls {
					fmt.Println(ve)
				}
			} else {
				fmt.Println("No VCS URL's specified")
			}
			if len(blstruct.prereq.vcs.execute) > 0 {
				fmt.Println("VCS related commands to be executed:")
				for _, ve := range blstruct.prereq.vcs.execute {
					fmt.Println(ve)
				}
			} else {
				fmt.Println("No VCS commands to execute")
			}
			fmt.Println("Please review your baseline if either of these are empty")
		}
		// prerequisite script
		if blstruct.prereq.script == "" {
			fmt.Println("No prerequisite script information specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			fmt.Println("Script:")
			fmt.Println(blstruct.prereq.script)
		}
		// prerequisite cleanup
		if blstruct.prereq.cleanup == false {
			fmt.Println("Prerequisite cleanup is set to false.")
		} else {
			fmt.Println("Prerequisite cleanup is set to true.")
		}

	}
	fmt.Println("Verifying server group's must-have list:")
	if len(blstruct.musthave.installed) == 0 && // done
		len(blstruct.musthave.enabled) == 0 && // done
		len(blstruct.musthave.disabled) == 0 && // done
		len(blstruct.musthave.configured.services) == 0 &&
		len(blstruct.musthave.users.users) == 0 &&
		blstruct.musthave.policies.polimport == "" &&
		blstruct.musthave.policies.polreboot == false &&
		blstruct.musthave.policies.polstatus == "" &&
		len(blstruct.musthave.rules.fwopen.ports) == 0 &&
		len(blstruct.musthave.rules.fwopen.protocols) == 0 &&
		len(blstruct.musthave.rules.fwclosed.ports) == 0 &&
		len(blstruct.musthave.rules.fwclosed.protocols) == 0 &&
		len(blstruct.musthave.rules.fwzones) == 0 &&
		len(blstruct.musthave.mounts.mountname) == 0 {
		fmt.Println("No must-haves have been specified  -- Please check your baseline, if you believe this to be incorrect")
	} else {
		// MH installed
		if len(blstruct.musthave.installed) == 0 {
			fmt.Println("No must-have installed information specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			if len(blstruct.musthave.installed) > 0 {
				fmt.Println("The following must-have tools/software will be installed, if they haven't been installed previously:")
				for _, ve := range blstruct.musthave.installed {
					fmt.Println(ve)
				}
			} else {
				fmt.Println("No must-have installed specified")
			}
		}
		// MH enabled
		if len(blstruct.musthave.enabled) == 0 {
			fmt.Println("No must-have enabled information specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			if len(blstruct.musthave.enabled) > 0 {
				fmt.Println("The following must-have tools/software will be enabled, if they haven't been enabled previously:")
				for _, ve := range blstruct.musthave.enabled {
					fmt.Println(ve)
				}
			} else {
				fmt.Println("No must-have enabled specified")
			}
		}
		// MH disabled
		if len(blstruct.musthave.disabled) == 0 {
			fmt.Println("No must-have disabled information specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			if len(blstruct.musthave.disabled) > 0 {
				fmt.Println("The following must-have tools/software will be disabled, if they haven't been disabled previously:")
				for _, ve := range blstruct.musthave.disabled {
					fmt.Println(ve)
				}
			} else {
				fmt.Println("No must-have disabled specified")
			}
		}
		// MH configured
		for ke, ve := range blstruct.musthave.configured.services {
			if ke == "" {
				fmt.Println("No must-have configurations specified  -- Please check your baseline, if you believe this to be incorrect")
				break
			} else {
				fmt.Printf("%s:\n", ke)
				for _, val := range ve.source {
					fmt.Printf("Source: %s\n", val)
				}
				for _, val := range ve.destination {
					fmt.Printf("Destination: %s\n", val)
				}
			}
		}
		// MH Users
		for ke, ve := range blstruct.musthave.users.users {
			if ke == "" {
				fmt.Println("No must-have users specified  -- Please check your baseline, if you believe this to be incorrect")
				break
			} else {
				fmt.Printf("%s:\n", ke)
				for _, val := range ve.groups {
					fmt.Printf("Groups: %s\n", val)
				}
				fmt.Printf("Shell: %v\n", ve.shell)
				fmt.Printf("Home: %v\n", ve.home)
				fmt.Printf("Sudoer: %v\n", ve.sudoer)
			}
		}
		// MH Policies
		if blstruct.musthave.policies.polstatus == "" &&
			blstruct.musthave.policies.polimport == "" &&
			blstruct.musthave.policies.polreboot == false {
			fmt.Println("No must-have policies specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			if blstruct.musthave.policies.polstatus != "" {
				fmt.Printf("Status: %s\n", blstruct.musthave.policies.polstatus)
			} else {
				fmt.Println("No must-have policy status specified  -- Please check your baseline, if you believe this to be incorrect")
			}
			if blstruct.musthave.policies.polimport != "" {
				fmt.Printf("Import: %s\n", blstruct.musthave.policies.polimport)
			} else {
				fmt.Println("No must-have policy to import  -- Please check your baseline, if you believe this to be incorrect")
			}
			if blstruct.musthave.policies.polreboot != false {
				fmt.Printf("Reboot: %v\n", blstruct.musthave.policies.polreboot)
			} else {
				fmt.Println("Reboot set to false, or not in baseline  -- Please check your baseline, if you believe this to be incorrect")
			}
		}
		// MH Firewall rules
		if len(blstruct.musthave.rules.fwopen.ports) == 0 &&
			len(blstruct.musthave.rules.fwopen.protocols) == 0 &&
			len(blstruct.musthave.rules.fwclosed.ports) == 0 &&
			len(blstruct.musthave.rules.fwclosed.protocols) == 0 &&
			len(blstruct.musthave.rules.fwzones) == 0 {
			fmt.Println("No must-have firewall rules specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			fmt.Println("Firewall rules:")
			if len(blstruct.musthave.rules.fwopen.ports) > 0 {
				fmt.Println("Open ports:")
				for _, ve := range blstruct.musthave.rules.fwopen.ports {
					fmt.Println(ve)
				}
			} else {
				fmt.Println("No open ports specified  -- Please check your baseline, if you believe this to be incorrect")
			}
			if len(blstruct.musthave.rules.fwopen.protocols) > 0 {
				fmt.Println("Open protocols:")
				for _, ve := range blstruct.musthave.rules.fwopen.protocols {
					fmt.Println(ve)
				}
			} else {
				fmt.Println("No open protocols specified  -- Please check your baseline, if you believe this to be incorrect")
			}
			if len(blstruct.musthave.rules.fwclosed.ports) > 0 {
				fmt.Println("Closed ports:")
				for _, ve := range blstruct.musthave.rules.fwclosed.ports {
					fmt.Println(ve)
				}
			} else {
				fmt.Println("No closed ports specified  -- Please check your baseline, if you believe this to be incorrect")
			}
			if len(blstruct.musthave.rules.fwclosed.protocols) > 0 {
				fmt.Println("Closed protocols:")
				for _, ve := range blstruct.musthave.rules.fwclosed.protocols {
					fmt.Println(ve)
				}
			} else {
				fmt.Println("No closed protocols specified  -- Please check your baseline, if you believe this to be incorrect")
			}
			if len(blstruct.musthave.rules.fwzones) > 0 {
				fmt.Println("Firewall zones:")
				for _, ve := range blstruct.musthave.rules.fwzones {
					fmt.Println(ve)
				}
			} else {
				fmt.Println("No firewall zones specified  -- Please check your baseline, if you believe this to be incorrect")
			}
		}
		// MH mounts
		if len(blstruct.musthave.mounts.mountname) == 0 {
			fmt.Println("No must-have mount details specified  -- Please check your baseline, if you believe this to be incorrect")
		} else {
			// 							handle ifs here
		}
	}
	//MNH list
	fmt.Println("Verifying server group's must-not-have list:")
	if len(blstruct.mustnothave.installed) == 0 && // done
		len(blstruct.mustnothave.enabled) == 0 && // done
		len(blstruct.mustnothave.disabled) == 0 && // done
		len(blstruct.mustnothave.users) == 0 &&
		len(blstruct.mustnothave.rules.fwopen.ports) == 0 &&
		len(blstruct.mustnothave.rules.fwopen.protocols) == 0 &&
		len(blstruct.mustnothave.rules.fwclosed.ports) == 0 &&
		len(blstruct.mustnothave.rules.fwclosed.protocols) == 0 &&
		len(blstruct.mustnothave.rules.fwzones) == 0 &&
		len(blstruct.mustnothave.mounts) == 0 {
		fmt.Println("No must-not-haves have been specified  -- Please check your baseline, if you believe this to be incorrect")
	} else {
		// 							handle ifs here
	}
	// Final steps list
	fmt.Println("Verifying server group's final steps list:")
	if len(blstruct.final.scripts) == 0 &&
		len(blstruct.final.commands) == 0 &&
		len(blstruct.final.collect.logs) == 0 &&
		len(blstruct.final.collect.stats) == 0 &&
		len(blstruct.final.collect.files) == 0 &&
		blstruct.final.collect.users == false &&
		blstruct.final.restart.services == false &&
		blstruct.final.restart.servers == false {
		fmt.Println("No final steps have been specified  -- Please check your baseline, if you believe this to be incorrect")
	} else {
		// 							handle ifs here
	}
}
