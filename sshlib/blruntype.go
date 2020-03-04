package sshlib

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

// ApplyBaselines testing
func ApplyBaselines(baselineyaml *yaml.MapSlice) {
	for _, groupItem := range *baselineyaml {
		fmt.Printf("Processing %s:\n", groupItem.Key)
		groupValue, ok := groupItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
		}
		for _, serverItem := range groupValue {
			servername := serverItem.Key
			fmt.Println(servername)
		}
	}
}

// CheckBaselines testing
func CheckBaselines(baselineyaml *yaml.MapSlice) {
	var blstruct ParsedBaseline
	var servergroupnames []string
	// first - BL names
	for _, blItem := range *baselineyaml {
		fmt.Printf("%s:\n", blItem.Key)
		groupValues, ok := blItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Check your baseline for issues\nAlternatively generate a template to see what is missing/wrong\n"))
		}
		// second - Server group names
		for _, groupItem := range groupValues {
			servergroupnames = append(servergroupnames, groupItem.Key.(string)) // done
			fmt.Printf("%s:\n", groupItem.Key)
			blstepsValue, ok := groupItem.Value.(yaml.MapSlice)
			if !ok {
				panic(fmt.Sprintf("Error reading server groups.\nAborting...\n"))
			}
			// third - BL steps or phases (Excludes, Prerequisites, Must-Haves, Must-Not-Haves,etc)
			for _, blstepItem := range blstepsValue {
				nextValues, ok := blstepItem.Value.(yaml.MapSlice)
				if !ok {
					// If excludes/prereqs/etc are missing or empty, create empty/blank data for structs
					// skipping those steps. An extra error will be created if too many fields are missing
				}
				blstepcheck := blstepItem.Key
				if blstepItem.Key == nil {
					fmt.Println("blank this step")
				}
				// fourth - OS, Servers, Tools, Files, VCS, etc
				for _, thirdStep := range nextValues {
					nnnextValue, ok := thirdStep.Value.(yaml.MapSlice)
					if !ok {
						// If excludes/prereqs/etc are missing or empty, create empty/blank data dor structs
						// skipping those steps. An extra error will be created if too many fields are missing
					}
					//check OS
					switch blstepcheck {
					case "Exclude":
						switch thirdStep.Key {
						case "OS":
							exclOS := make([]string, len(thirdStep.Value.([]interface{})))
							OSslice := thirdStep.Value.([]interface{})
							for i, v := range OSslice {
								exclOS[i] = v.(string)
							}
							blstruct.exclude.osExcl = exclOS
							fmt.Println("OS stored")
							time.Sleep(1 * time.Second)
						case "Servers":
							exclServers := make([]string, len(thirdStep.Value.([]interface{})))
							serverSlice := thirdStep.Value.([]interface{})
							for i, v := range serverSlice {
								exclServers[i] = v.(string)
							}
							blstruct.exclude.serversExcl = exclServers
							fmt.Println("servers stored")
							time.Sleep(1 * time.Second)
						default:
							fmt.Println("Nothing to exclude")
						}
					case "Prerequisites":
						switch thirdStep.Key {
						case "Tools":
							prereqTools := make([]string, len(thirdStep.Value.([]interface{})))
							prereqToolsSlice := thirdStep.Value.([]interface{})
							for i, v := range prereqToolsSlice {
								prereqTools[i] = v.(string)
							}
							blstruct.prereq.tools = prereqTools
							fmt.Println("Prerequisite tools stored")
							time.Sleep(1 * time.Second)
						case "Files": // for loop
							fmt.Println("\t\tstill working on Files")
						case "VCS": // for loop
							fmt.Println("\t\tstill working on VCS")
						case "Script":
							var prereqScript string
							prereqScript = thirdStep.Value.(string)
							blstruct.prereq.script = prereqScript

						case "Clean-up":
							var prereqCU bool
							prereqCU = thirdStep.Value.(bool)
							blstruct.prereq.cleanup = prereqCU

						default:
						}
					case "Must-Have":
						switch thirdStep.Key {
						case "Installed":
							mhInst := make([]string, len(thirdStep.Value.([]interface{})))
							mhInstSlice := thirdStep.Value.([]interface{})
							for i, v := range mhInstSlice {
								mhInst[i] = v.(string)
							}
							blstruct.musthave.installed = mhInst
							fmt.Println("Must-Have installed stored")
							time.Sleep(1 * time.Second)
						case "Enabled":
							mhEnabled := make([]string, len(thirdStep.Value.([]interface{})))
							mhEnabledSlice := thirdStep.Value.([]interface{})
							for i, v := range mhEnabledSlice {
								mhEnabled[i] = v.(string)
							}
							blstruct.musthave.enabled = mhEnabled
							fmt.Println("Must-Have enabled stored")
							time.Sleep(1 * time.Second)
						case "Disabled":
							if thirdStep.Value == nil {
								continue // continue is a place holder for now
							} else {
								mhDisabled := make([]string, len(thirdStep.Value.([]interface{})))
								mhDisabledSlice := thirdStep.Value.([]interface{})
								for i, v := range mhDisabledSlice {
									mhDisabled[i] = v.(string)
								}
								blstruct.musthave.enabled = mhDisabled
								fmt.Println("Must-Have disabled stored")
								time.Sleep(1 * time.Second)
							}
						case "Configured": // for loop
							fmt.Println("\t\tstill working on Configured")
						case "Users": // for loop
							fmt.Println("\t\tstill working on users")
						case "Policies": // for loop
							fmt.Println("\t\tstill working on Policies")
						case "Rules": // for loop
							fmt.Println("\t\tstill working on Rules")
						case "Mounts": // for loop
							fmt.Println("\t\tstill working on Mounts")
						default:
						}
					case "Must-Not-Have":
						switch thirdStep.Key {
						case "Installed":
							mnhInst := make([]string, len(thirdStep.Value.([]interface{})))
							mnhInstSlice := thirdStep.Value.([]interface{})
							for i, v := range mnhInstSlice {
								mnhInst[i] = v.(string)
							}
							blstruct.mustnothave.installed = mnhInst
							fmt.Println("Must-Not-Have installed stored")
							time.Sleep(1 * time.Second)
						case "Enabled":
							mnhEnabled := make([]string, len(thirdStep.Value.([]interface{})))
							mnhEnabledSlice := thirdStep.Value.([]interface{})
							for i, v := range mnhEnabledSlice {
								mnhEnabled[i] = v.(string)
							}
							blstruct.mustnothave.installed = mnhEnabled
							fmt.Println("Must-Not-Have enabled stored")
							time.Sleep(1 * time.Second)
						case "Disabled":
							mnhDisabled := make([]string, len(thirdStep.Value.([]interface{})))
							mnhDisabledSlice := thirdStep.Value.([]interface{})
							for i, v := range mnhDisabledSlice {
								mnhDisabled[i] = v.(string)
							}
							blstruct.mustnothave.enabled = mnhDisabled
							fmt.Println("Must-Not-Have disabled stored")
							time.Sleep(1 * time.Second)
						case "Users":
							mnhUsers := make([]string, len(thirdStep.Value.([]interface{})))
							mnhUsersSlice := thirdStep.Value.([]interface{})
							for i, v := range mnhUsersSlice {
								mnhUsers[i] = v.(string)
							}
							blstruct.mustnothave.users = mnhUsers
							fmt.Println("Must-Not-Have users stored")
							time.Sleep(1 * time.Second)
						case "Rules": // for loop
							fmt.Println("\t\tstill working on MNH rules")
						case "Mounts":
							mnhMounts := make([]string, len(thirdStep.Value.([]interface{})))
							mnhMountsSlice := thirdStep.Value.([]interface{})
							for i, v := range mnhMountsSlice {
								mnhMounts[i] = v.(string)
							}
							blstruct.mustnothave.mounts = mnhMounts
							fmt.Println("Must-Not-Have users stored")
							time.Sleep(1 * time.Second)
						default:
						}
					case "Final":
						switch thirdStep.Key {
						case "Scripts":
							fnlScripts := make([]string, len(thirdStep.Value.([]interface{})))
							fnlScriptsSlice := thirdStep.Value.([]interface{})
							for i, v := range fnlScriptsSlice {
								fnlScripts[i] = v.(string)
							}
							blstruct.final.scripts = fnlScripts
							fmt.Println("Final scripts stored")
							time.Sleep(1 * time.Second)
						case "Commands":
							fnlCommands := make([]string, len(thirdStep.Value.([]interface{})))
							fnlCommandsSlice := thirdStep.Value.([]interface{})
							for i, v := range fnlCommandsSlice {
								fnlCommands[i] = v.(string)
							}
							blstruct.final.commands = fnlCommands
							fmt.Println("Final commands stored")
							time.Sleep(1 * time.Second)
						case "Collect": // for loop
							fmt.Println("\t\tstill working on collect")
						case "Restart": // for loop
							fmt.Println("\t\tstill working on restart")
						default:
						}
					default:
					}

					// test2 := strings.Split(test, " ")
					//  fifth
					for _, ffforItem := range nnnextValue {
						// fmt.Println("4 ", ffforItem.Value)
						// ffforname := ffforItem.Value
						// fffornamekey := ffforItem.Key
						// fmt.Println("\t", ffforname)
						// fmt.Println(fffornamekey)
						nnnnextValue, ok := ffforItem.Value.(yaml.MapSlice)
						if !ok {
							// panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
						}
						// sixth
						for _, fffforItem := range nnnnextValue {
							// fmt.Printf("%s:\n", fffforItem.Key)
							// fffforname := fffforItem.Value
							// ffffornamekey := fffforItem.Key
							// fmt.Println("\t", fffforname)
							// fmt.Println(ffffornamekey)
							nnnnnextValue, ok := fffforItem.Value.(yaml.MapSlice)
							if !ok {
								// panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
							}
							for _, ffffforItem := range nnnnnextValue {
								// fmt.Printf("%s:\n", ffffforItem.Key)
								// somevalue := ffffforItem.Value
								// fmt.Println("\t", somevalue)
								_, ok := ffffforItem.Value.(yaml.MapSlice)
								if !ok {
									// panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
								}
							}
						}
					}
				}
			}
		}
		// run pool parser here and bypass excluded servers and OS
		// trying to avoid too many unnecessary loops
		// should not forget about 'all'/'All'/'ALL' in servergroups
	}
	fmt.Println(servergroupnames)
	fmt.Println(blstruct.exclude.osExcl)
	fmt.Println(blstruct.exclude.serversExcl)
}
