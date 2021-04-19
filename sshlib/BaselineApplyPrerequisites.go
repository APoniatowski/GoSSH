package sshlib

import (
	"fmt"
	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"strings"
)

func (baselineStruct *ParsedBaseline) applyPrereq(sshList *map[string]string, commandChannel chan<- map[string]string) {
	commandSet := make(map[string]string)
	// skipping because of this check... need to make a few changes
	fmt.Printf("Prerequisites Checklist: ")
	if len(baselineStruct.prereq.vcs.execute) == 0 &&
		len(baselineStruct.prereq.vcs.urls) == 0 &&
		baselineStruct.prereq.files.local.dest == "" &&
		baselineStruct.prereq.files.local.src == "" &&
		baselineStruct.prereq.files.remote.address == "" &&
		baselineStruct.prereq.files.remote.dest == "" &&
		baselineStruct.prereq.files.remote.mounttype == "" &&
		baselineStruct.prereq.files.remote.pwd == "" &&
		baselineStruct.prereq.files.remote.src == "" &&
		baselineStruct.prereq.files.remote.username == "" &&
		len(baselineStruct.prereq.files.remote.files) == 0 &&
		len(baselineStruct.prereq.files.urls) == 0 &&
		len(baselineStruct.prereq.tools) == 0 &&
		baselineStruct.prereq.script == "" &&
		!baselineStruct.prereq.cleanup {
		commandSet[""] = ""
		fmt.Printf("Skipping...\n")
	} else {
		fmt.Printf("\n")
		// prerequisite tools
		fmt.Printf(" Prerequisite Tools: ")
		if len(baselineStruct.prereq.tools) == 0 {
			fmt.Printf("Skipping...\n")
		} else {
			fmt.Printf("\n")
			for _, ve := range baselineStruct.prereq.tools {
				for key, val := range *sshList {
					if commandSet[val] == "" {
						commandSet[key] = serviceCommandBuilder(&ve, &val, "install")
					}
				}
				commandChannel <- commandSet
				//TODO Prereq Tools apply
				//for k, v := range commandSet {
				//	fmt.Printf("%v   %v\n", k, v)
				//}
				//send to channel
				//wait for response and display compliancy
			}
		}
		// prerequisite files URLs
		fmt.Printf(" Prerequisite URL's: ")
		if len(baselineStruct.prereq.files.urls) == 0 {
			fmt.Printf("Skipping...\n")
		} else {
			fmt.Printf("\n")
			for _, ve := range baselineStruct.prereq.files.urls {
				for key, val := range *sshList {
					if commandSet[val] == "" {
						commandSet[key] = prereqURLFetch(&ve)
					}
				}
				commandChannel <- commandSet
				// TODO URL Files
				//for k, v := range commandSet {
				//	fmt.Printf("%v   %v\n", k, v)
				//}
				// send to channel
				// wait for response and display compliancy
			}
		}
		// prerequisite files local
		fmt.Printf(" Prerequisite Files (network transfer): ")
		if baselineStruct.prereq.files.local.dest != "" &&
			baselineStruct.prereq.files.local.src != "" {
			fmt.Printf("Skipping...\n")
		} else {
			fmt.Printf("\n")
			var srcFile string
			if baselineStruct.prereq.files.local.src != "" {
				srcFile = baselineStruct.prereq.files.local.src
			}
			if baselineStruct.prereq.files.local.dest != "" {
				for key, val := range *sshList {
					if commandSet[val] == "" {
						// TODO Prereq SCP Files apply make some changes and move to cmdbuilders
						commandSet[key] = pkgmanlib.OmniTools["suminfo"] + baselineStruct.prereq.files.local.dest + srcFile
						/*
							-will need to find a better way to compare files and directories-
							cat would kill memory, if its a large file or binary
							sum only does files, not dirs
							need to create for loop command if its a directory with md5sum
						*/
					}
				}
				commandChannel <- commandSet
				//for k, v := range commandSet {
				//	fmt.Printf("%v   %v\n", k, v)
				//}
				// send to channel
				// wait for response
				// diff the file/dir with the source
				// display compliancy
			}
		}
		// prerequisite files remote
		fmt.Printf(" Prerequisite Files (via mount): ")
		if baselineStruct.prereq.files.remote.address == "" &&
			baselineStruct.prereq.files.remote.dest == "" &&
			baselineStruct.prereq.files.remote.mounttype == "" &&
			baselineStruct.prereq.files.remote.pwd == "" &&
			baselineStruct.prereq.files.remote.src == "" &&
			baselineStruct.prereq.files.remote.username == "" &&
			len(baselineStruct.prereq.files.remote.files) == 0 {
			fmt.Printf("Skipping...\n")
		} else {
			if baselineStruct.prereq.files.remote.src != "" {
				if len(baselineStruct.prereq.files.remote.files) != 0 {
					for _, ve := range baselineStruct.prereq.files.remote.files {
						for key, val := range *sshList {
							if commandSet[val] == "" {
								commandSet[key] = baselineStruct.prereq.files.remote.remoteFilesCommandBuilder(&ve, "apply")
							}
						}
						commandChannel <- commandSet
						// TODO Prereq Mount Files apply
						//for k, v := range commandSet {
						//	fmt.Printf("%v   %v\n", k, v)
						//}
						// send to channel
						// wait for response
						// diff the file/dir with the source
						// display compliancy
					}
				}
			}
			fmt.Printf("\n")
		}
		// prerequisite VCS instructions
		fmt.Printf(" Prerequisite Files (via VCS): ")
		if len(baselineStruct.prereq.vcs.execute) == 0 &&
			len(baselineStruct.prereq.vcs.urls) == 0 {
			fmt.Printf("Skipping...\n")
		} else {
			fmt.Printf("\n")
			var parsedFile []string
			if len(baselineStruct.prereq.vcs.urls) > 0 {
				for _, ve := range baselineStruct.prereq.vcs.urls {
					parseFile := strings.Split(ve, "/")
					parsedFile = append(parsedFile, parseFile[len(parseFile)-1])
					for key, val := range *sshList {
						if commandSet[val] == "" {
							// TODO Prereq VCS Files apply
							commandSet[key] = prereqURLFetch(&ve)
						}
					}
					commandChannel <- commandSet
					//for k, v := range commandSet {
					//	fmt.Printf("%v   %v\n", k, v)
					//}
					// send to channel
					// wait for response
				}
			}
			if len(baselineStruct.prereq.vcs.execute) > 0 {
				for _, ve := range baselineStruct.prereq.vcs.execute {
					for key, val := range *sshList {
						if commandSet[val] == "" {
							// TODO Prereq VCS execution
							commandSet[key] = ve
						}
					}
					commandChannel <- commandSet
					//for k, v := range commandSet {
					//	fmt.Printf("%v   %v\n", k, v)
					//}
					// send to channel
					// wait for response
				}
			}
		}
	}
	return
}
