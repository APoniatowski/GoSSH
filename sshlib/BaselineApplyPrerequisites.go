package sshlib

import (
	"fmt"
	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"strings"
)

func (blstruct *ParsedBaseline) applyPrereq(sshList *map[string]string) {
	commandset := make(map[string]string)
	// skipping because of this check... need to make a few changes
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
				for key, val := range *sshList {
					if commandset[val] == "" {
						commandset[key] = serviceCommandBuilder(&ve, &val, "install")
					}
				}
				// TODO Prereq Tools apply
				//for k, v := range commandset {
				//	fmt.Printf("%v   %v\n", k, v)
				//}
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
				for key, val := range *sshList {
					if commandset[val] == "" {
						commandset[key] = prereqURLFetch(&ve)
					}
				}
				// TODO URL Files
				//for k, v := range commandset {
				//	fmt.Printf("%v   %v\n", k, v)
				//}
				// send to channel
				// wait for response and display compliancy
			}
		}
		// prerequisite files local
		fmt.Printf(" Prerequisite Files (network transfer): ")
		if blstruct.prereq.files.local.dest != "" &&
			blstruct.prereq.files.local.src != "" {
			fmt.Printf("Skipping...\n")
		} else {
			fmt.Printf("\n")
			var srcFile string
			if blstruct.prereq.files.local.src != "" {
				srcFile = blstruct.prereq.files.local.src
			}
			if blstruct.prereq.files.local.dest != "" {
				for key, val := range *sshList {
					if commandset[val] == "" {
						// TODO Prereq SCP Files apply make some changes and move to cmdbuilders
						commandset[key] = pkgmanlib.OmniTools["suminfo"] + blstruct.prereq.files.local.dest + srcFile
						/*
							-will need to find a better way to compare files and directories-
							cat would kill memory, if its a large file or binary
							sum only does files, not dirs
							need to create for loop command if its a directory with md5sum
						*/
					}
				}
				//for k, v := range commandset {
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
		if blstruct.prereq.files.remote.address == "" &&
			blstruct.prereq.files.remote.dest == "" &&
			blstruct.prereq.files.remote.mounttype == "" &&
			blstruct.prereq.files.remote.pwd == "" &&
			blstruct.prereq.files.remote.src == "" &&
			blstruct.prereq.files.remote.username == "" &&
			len(blstruct.prereq.files.remote.files) == 0 {
			fmt.Printf("Skipping...\n")
		} else {
			if blstruct.prereq.files.remote.src != "" {
				if len(blstruct.prereq.files.remote.files) != 0 {
					for _, ve := range blstruct.prereq.files.remote.files {
						for key, val := range *sshList {
							if commandset[val] == "" {
								commandset[key] = blstruct.prereq.files.remote.remoteFilesCommandBuilder(&ve, "apply")
							}
						}
						// TODO Prereq Mount Files apply
						//for k, v := range commandset {
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
		if len(blstruct.prereq.vcs.execute) == 0 &&
			len(blstruct.prereq.vcs.urls) == 0 {
			fmt.Printf("Skipping...\n")
		} else {
			fmt.Printf("\n")
			var parsedFile []string
			if len(blstruct.prereq.vcs.urls) > 0 {
				for _, ve := range blstruct.prereq.vcs.urls {
					parseFile := strings.Split(ve, "/")
					parsedFile = append(parsedFile, parseFile[len(parseFile)-1])
					for key, val := range *sshList {
						if commandset[val] == "" {
							// TODO Prereq VCS Files apply
							commandset[key] = prereqURLFetch(&ve)
						}
					}
					//for k, v := range commandset {
					//	fmt.Printf("%v   %v\n", k, v)
					//}
					// send to channel
					// wait for response
				}
			}
			if len(blstruct.prereq.vcs.execute) > 0 {
				for _, ve := range blstruct.prereq.vcs.execute {
					for key, val := range *sshList {
						if commandset[val] == "" {
							// TODO Prereq VCS execution
							commandset[key] = ve
						}
					}
					//for k, v := range commandset {
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
