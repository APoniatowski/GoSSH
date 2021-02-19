package sshlib

import (
	"fmt"
	//"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
)

// ApplyBaselines Apply defined baselines
func ApplyBaselines(baselineyaml *yaml.MapSlice, configs *yaml.MapSlice) {
	var warnings int
	var maincategorywarnings int
	var datawarnings int
	var blerrors int
	// Baseline
	for _, blItem := range *baselineyaml {
		fmt.Printf("%s:\n", blItem.Key)
		groupValues, ok := blItem.Value.(yaml.MapSlice)
		if !ok {
			panic("\nCheck your baseline for issues\nAlternatively generate a template to see what is missing/wrong\n")
		}
		// Server groups
		for _, groupItem := range groupValues {
			// initialize the data
			servergroupname := groupItem.Key.(string)
			var blstruct ParsedBaseline
			blstruct.musthave.configured.services = make(map[string]musthaveconfiguredservices)
			blstruct.musthave.users.users = make(map[string]musthaveusersstruct)
			blstruct.musthave.mounts.mountname = make(map[string]mountdetails)
			blstepsValue, ok := groupItem.Value.(yaml.MapSlice)
			if !ok {
				blerrors++
				panic("\nError parsing server groups.\nAborting...\n")
			}
			if strings.ToLower(servergroupname) == "all" {
				fmt.Println("Checking baseline on all servers:")
			} else {
				fmt.Println("Checking baseline on", servergroupname+":")
			}
			// Exclude, Prerequisites, Must-Have, Must-Not-Have, Final
			for _, blstepItem := range blstepsValue {
				nextValues, ok := blstepItem.Value.(yaml.MapSlice)
				if !ok {
					maincategorywarnings++
				}
				blstepcheck := blstepItem.Key
				if blstepItem.Key == nil {
					blerrors++
				}

				// OS, Servers, Tools, Files, VCS, etc
				for _, thirdStep := range nextValues {
					nextblValues, ok := thirdStep.Value.(yaml.MapSlice)
					if !ok {
						warnings++
					}
					if thirdStep.Key == nil {
						warnings++
						blerrors++
					} else {
						switch blstepcheck {
						case "Exclude":
							switch thirdStep.Key {
							case "OS":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.exclude.osExcl = []string{""}
								} else {
									exclOS := make([]string, len(thirdStep.Value.([]interface{})))
									OSslice := thirdStep.Value.([]interface{})
									for i, v := range OSslice {
										exclOS[i] = v.(string)
									}
									blstruct.exclude.osExcl = exclOS
								}
							case "Servers":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.exclude.serversExcl = []string{""}
								} else {
									exclServers := make([]string, len(thirdStep.Value.([]interface{})))
									serverSlice := thirdStep.Value.([]interface{})
									for i, v := range serverSlice {
										exclServers[i] = v.(string)
									}
									blstruct.exclude.serversExcl = exclServers
								}
							}
						case "Prerequisites":
							switch thirdStep.Key {
							case "Tools":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.prereq.tools = []string{""}
								} else {
									prereqTools := make([]string, len(thirdStep.Value.([]interface{})))
									prereqToolsSlice := thirdStep.Value.([]interface{})
									for i, v := range prereqToolsSlice {
										prereqTools[i] = v.(string)
									}
									blstruct.prereq.tools = prereqTools
								}
							case "Files": // for loop
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.prereq.files.urls = []string{""}
									blstruct.prereq.files.local.src = ""
									blstruct.prereq.files.local.dest = ""
									blstruct.prereq.files.remote.mounttype = ""
									blstruct.prereq.files.remote.address = ""
									blstruct.prereq.files.remote.username = ""
									blstruct.prereq.files.remote.pwd = ""
									blstruct.prereq.files.remote.src = ""
									blstruct.prereq.files.remote.dest = ""
								} else {
									for _, blItem = range nextblValues {
										switch blItem.Key {
										case "URLs":
											if blItem.Value == nil {
												datawarnings++
												blstruct.prereq.files.urls = []string{""}
											} else {
												prereqURLs := make([]string, len(blItem.Value.([]interface{})))
												prereqURLsSlice := blItem.Value.([]interface{})
												for i, v := range prereqURLsSlice {
													prereqURLs[i] = v.(string)
												}
												blstruct.prereq.files.urls = prereqURLs
											}
										case "Local":
											if blItem.Value == nil {
												datawarnings++
												blstruct.prereq.files.local.src = ""
												blstruct.prereq.files.local.dest = ""
											} else {
												extrablValues, ok := blItem.Value.(yaml.MapSlice)
												if !ok {
													fmt.Println("Error parsing baseline. Please check the baseline you specified or generate a template")
												}
												var nextblStep yaml.MapItem
												for _, nextblStep = range extrablValues {
													switch nextblStep.Key {
													case "Source":
														if nextblStep.Value == nil {
															datawarnings++
															blstruct.prereq.files.local.src = ""
														} else {
															blstruct.prereq.files.local.src = nextblStep.Value.(string)
														}
													case "Destination":
														if nextblStep.Value == nil {
															datawarnings++
															blstruct.prereq.files.local.dest = ""
														} else {
															blstruct.prereq.files.local.dest = nextblStep.Value.(string)
														}
													}
												}
											}
										case "Remote":
											if blItem.Value == nil {
												datawarnings++
												blstruct.prereq.files.remote.mounttype = ""
												blstruct.prereq.files.remote.address = ""
												blstruct.prereq.files.remote.username = ""
												blstruct.prereq.files.remote.pwd = ""
												blstruct.prereq.files.remote.src = ""
												blstruct.prereq.files.remote.dest = ""
											} else {
												extrablValues, ok := blItem.Value.(yaml.MapSlice)
												if !ok {
													warnings++
												}
												var nextblStep yaml.MapItem
												for _, nextblStep = range extrablValues {
													switch nextblStep.Key {
													case "Type":
														if nextblStep.Value == nil {
															datawarnings++
															blstruct.prereq.files.remote.mounttype = ""
														} else {
															blstruct.prereq.files.remote.mounttype = nextblStep.Value.(string)
														}
													case "Address":
														if nextblStep.Value == nil {
															datawarnings++
															blstruct.prereq.files.remote.address = ""
														} else {
															blstruct.prereq.files.remote.address = nextblStep.Value.(string)
														}
													case "Username":
														if nextblStep.Value == nil {
															datawarnings++
															blstruct.prereq.files.remote.username = ""
														} else {
															blstruct.prereq.files.remote.username = nextblStep.Value.(string)
														}
													case "Password":
														if nextblStep.Value == nil {
															datawarnings++
															blstruct.prereq.files.remote.pwd = ""
														} else {
															blstruct.prereq.files.remote.pwd = nextblStep.Value.(string)
														}
													case "Source":
														if nextblStep.Value == nil {
															datawarnings++
															blstruct.prereq.files.remote.src = ""
														} else {
															blstruct.prereq.files.remote.src = nextblStep.Value.(string)
														}
													case "Destination":
														if nextblStep.Value == nil {
															datawarnings++
															blstruct.prereq.files.remote.dest = ""
														} else {
															blstruct.prereq.files.remote.dest = nextblStep.Value.(string)
														}
													case "Files":
														if nextblStep.Value == nil {
															datawarnings++
															blstruct.prereq.files.remote.files = []string{""}
														} else {
															remoteFiles := make([]string, len(nextblStep.Value.([]interface{})))
															remoteFilesSlice := nextblStep.Value.([]interface{})
															for i, v := range remoteFilesSlice {
																remoteFiles[i] = v.(string)
															}
															blstruct.prereq.files.remote.files = remoteFiles
														}
													}
												}
											}
										}
									}
								}
							case "VCS":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.prereq.vcs.urls = []string{""}
									blstruct.prereq.vcs.execute = []string{""}
								} else {
									for _, blItem = range nextblValues {
										switch blItem.Key {
										case "URLs":
											if blItem.Value == nil {
												datawarnings++
												blstruct.prereq.vcs.urls = []string{""}
											} else {
												vcsURLs := make([]string, len(blItem.Value.([]interface{})))
												vcsURLsSlice := blItem.Value.([]interface{})
												for i, v := range vcsURLsSlice {
													vcsURLs[i] = v.(string)
												}
												blstruct.prereq.vcs.urls = vcsURLs
											}
										case "Execute":
											if blItem.Value == nil {
												datawarnings++
												blstruct.prereq.vcs.execute = []string{""}
											} else {
												vcsCMDS := make([]string, len(blItem.Value.([]interface{})))
												vcsCMDSSlice := blItem.Value.([]interface{})
												for i, v := range vcsCMDSSlice {
													vcsCMDS[i] = v.(string)
												}
												blstruct.prereq.vcs.execute = vcsCMDS
											}
										}
									}
								}
							case "Script":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.prereq.script = ""
								} else {
									blstruct.prereq.script = thirdStep.Value.(string)
								}
							case "Clean-up":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.prereq.cleanup = false
								} else {
									blstruct.prereq.cleanup = thirdStep.Value.(bool)
								}
							}
						case "Must-Have":
							switch thirdStep.Key {
							case "Installed":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.musthave.installed = []string{""}
								} else {
									mhInst := make([]string, len(thirdStep.Value.([]interface{})))
									mhInstSlice := thirdStep.Value.([]interface{})
									for i, v := range mhInstSlice {
										mhInst[i] = v.(string)
									}
									blstruct.musthave.installed = mhInst
								}
							case "Enabled":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.musthave.enabled = []string{""}
								} else {
									mhEnabled := make([]string, len(thirdStep.Value.([]interface{})))
									mhEnabledSlice := thirdStep.Value.([]interface{})
									for i, v := range mhEnabledSlice {
										mhEnabled[i] = v.(string)
									}
									blstruct.musthave.enabled = mhEnabled
								}
							case "Disabled":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.musthave.disabled = []string{""}
								} else {
									mhDisabled := make([]string, len(thirdStep.Value.([]interface{})))
									mhDisabledSlice := thirdStep.Value.([]interface{})
									for i, v := range mhDisabledSlice {
										mhDisabled[i] = v.(string)
									}
									blstruct.musthave.enabled = mhDisabled
								}
							case "Configured":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.musthave.configured.services[""] = musthaveconfiguredservices{source: []string{""}, destination: []string{""}}
								} else {
									for _, blItem = range nextblValues {
										service := blItem.Key.(string)
										var mhConfSrc []string
										var mhConfDest []string
										confValues, ok := blItem.Value.(yaml.MapSlice)
										if !ok {
											warnings++
										}
										for _, confItems := range confValues {
											switch confItems.Key {
											case "Source":
												if confItems.Value == nil {
													datawarnings++
													blstruct.musthave.configured.services[service] = musthaveconfiguredservices{source: []string{""}}
												} else {
													mhConfSrc = make([]string, len(confItems.Value.([]interface{})))
													mhConfSrcSlice := confItems.Value.([]interface{})
													for i, v := range mhConfSrcSlice {
														mhConfSrc[i] = v.(string)
													}
												}
											case "Destination":
												if confItems.Value == nil {
													datawarnings++
													blstruct.musthave.configured.services[service] = musthaveconfiguredservices{destination: []string{""}}
												} else {
													mhConfDest = make([]string, len(confItems.Value.([]interface{})))
													mhConfDestSlice := confItems.Value.([]interface{})
													for i, v := range mhConfDestSlice {
														mhConfDest[i] = v.(string)
													}
												}
											}
											blstruct.musthave.configured.services[service] = musthaveconfiguredservices{source: mhConfSrc, destination: mhConfDest}
										}
									}
								}
							case "Users":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.musthave.users.users[""] = musthaveusersstruct{groups: []string{""}, shell: "", home: "", sudoer: false}
								} else {
									for _, blItem = range nextblValues {
										user := blItem.Key.(string)
										var mhUsergroup []string
										var mhUsershell string
										var mhUserhome string
										var mhUsersudo bool
										userValues, ok := blItem.Value.(yaml.MapSlice)
										if !ok {
											warnings++
										}
										for _, userItems := range userValues {
											switch userItems.Key {
											case "Groups":
												if userItems.Value == nil {
													datawarnings++
													mhUsergroup = []string{""}
												} else {
													mhUsergroup = make([]string, len(userItems.Value.([]interface{})))
													mhUsergroupSlice := userItems.Value.([]interface{})
													for i, v := range mhUsergroupSlice {
														mhUsergroup[i] = v.(string)
													}
												}
											case "Shell":
												if userItems.Value == nil {
													datawarnings++
													mhUsershell = ""
												} else {
													mhUsershell = userItems.Value.(string)
												}
											case "Home-Dir":
												if userItems.Value == nil {
													datawarnings++
													mhUserhome = ""
												} else {
													mhUserhome = userItems.Value.(string)
												}
											case "Sudoer":
												if userItems.Value == nil {
													datawarnings++
													mhUsersudo = false
												} else {
													mhUsersudo = userItems.Value.(bool)
												}
											}
											blstruct.musthave.users.users[user] = musthaveusersstruct{
												groups: mhUsergroup,
												shell:  mhUsershell,
												home:   mhUserhome,
												sudoer: mhUsersudo,
											}
										}
									}
								}
							case "Policies":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.musthave.policies.polimport = ""
									blstruct.musthave.policies.polstatus = ""
									blstruct.musthave.policies.polreboot = false
								} else {
									for _, blItem = range nextblValues {
										switch blItem.Key {
										case "Status":
											if blItem.Value == nil {
												datawarnings++
												blstruct.musthave.policies.polstatus = ""
											} else {
												blstruct.musthave.policies.polstatus = blItem.Value.(string)
											}
										case "Import":
											if blItem.Value == nil {
												datawarnings++
												blstruct.musthave.policies.polimport = ""
											} else {
												blstruct.musthave.policies.polimport = blItem.Value.(string)
											}
										case "Reboot":
											if blItem.Value == nil {
												datawarnings++
												blstruct.musthave.policies.polreboot = false
											} else {
												blstruct.musthave.policies.polreboot = blItem.Value.(bool)
											}
										}
									}
								}
							case "Rules":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.musthave.rules.fwopen.ports = []string{""}
									blstruct.musthave.rules.fwopen.protocols = []string{""}
									blstruct.musthave.rules.fwclosed.ports = []string{""}
									blstruct.musthave.rules.fwclosed.protocols = []string{""}
									blstruct.musthave.rules.fwzones = []string{""}
								} else {
									for _, blItem = range nextblValues {
										rulesValues, ok := blItem.Value.(yaml.MapSlice)
										if !ok {
											warnings++
										}
										switch blItem.Key {
										case "Open":
											for _, rulesItems := range rulesValues {
												switch rulesItems.Key {
												case "Ports":
													if rulesItems.Value == nil {
														datawarnings++
														blstruct.musthave.rules.fwopen.ports = []string{""}
													} else {
														mhRulesOPorts := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesOPortsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesOPortsSlice {
															mhRulesOPorts[i] = strconv.Itoa(v.(int))
														}
														blstruct.musthave.rules.fwopen.ports = mhRulesOPorts
													}
												case "Protocols":
													if rulesItems.Value == nil {
														datawarnings++
														blstruct.musthave.rules.fwopen.ports = []string{""}
													} else {
														mhRulesOProtocols := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesOProtocolsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesOProtocolsSlice {
															mhRulesOProtocols[i] = v.(string)
														}
														blstruct.musthave.rules.fwopen.protocols = mhRulesOProtocols
													}
												}
											}
										case "Closed":
											for _, rulesItems := range rulesValues {
												switch rulesItems.Key {
												case "Ports":
													if rulesItems.Value == nil {
														datawarnings++
														blstruct.musthave.rules.fwclosed.ports = []string{""}
													} else {
														mhRulesCPorts := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesCPortsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesCPortsSlice {
															mhRulesCPorts[i] = strconv.Itoa(v.(int))
														}
														blstruct.musthave.rules.fwclosed.ports = mhRulesCPorts
													}
												case "Protocols":
													if rulesItems.Value == nil {
														datawarnings++
														blstruct.musthave.rules.fwclosed.protocols = []string{""}
													} else {
														mhRulesCProtocols := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesCProtocolsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesCProtocolsSlice {
															mhRulesCProtocols[i] = v.(string)
														}
														blstruct.musthave.rules.fwclosed.protocols = mhRulesCProtocols
													}
												}
											}
										case "Zones":
											if blItem.Value == nil {
												datawarnings++
												blstruct.musthave.rules.fwzones = []string{""}
											} else {
												mhRulesZones := make([]string, len(blItem.Value.([]interface{})))
												mhRulesZonesSlice := blItem.Value.([]interface{})
												for i, v := range mhRulesZonesSlice {
													mhRulesZones[i] = v.(string)
												}
												blstruct.musthave.rules.fwzones = mhRulesZones
											}
										}
									}
								}
							case "Mounts":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.musthave.mounts.mountname[""] = mountdetails{
										mounttype: "",
										address:   "",
										username:  "",
										pwd:       "",
										src:       "",
										dest:      "",
									}
								} else {
									var mhMounts string
									var mhMountType string
									var mhAddress string
									var mhUsername string
									var mhPassword string
									var mhMountSource string
									var mhMountDest string
									if !ok {
										warnings++
									}
									for _, blItem = range nextblValues {
										mhMounts = blItem.Key.(string)
										mountValues, ok := blItem.Value.(yaml.MapSlice)
										if !ok {
											warnings++
										}
										for _, mountItems := range mountValues {
											switch mountItems.Key {
											case "Type":
												if mountItems.Value == nil {
													datawarnings++
													mhMountType = ""
												} else {
													mhMountType = mountItems.Value.(string)
												}
											case "Address":
												if mountItems.Value == nil {
													datawarnings++
													mhAddress = ""
												} else {
													mhAddress = mountItems.Value.(string)
												}
											case "Username":
												if mountItems.Value == nil {
													datawarnings++
													mhUsername = ""
												} else {
													mhUsername = mountItems.Value.(string)
												}
											case "Password":
												if mountItems.Value == nil {
													datawarnings++
													mhPassword = ""
												} else {
													mhPassword = mountItems.Value.(string)
												}
											case "Source":
												if mountItems.Value == nil {
													datawarnings++
													mhMountSource = ""
												} else {
													mhMountSource = mountItems.Value.(string)
												}
											case "Destination":
												if mountItems.Value == nil {
													datawarnings++
													mhMountDest = ""
												} else {
													mhMountDest = mountItems.Value.(string)
												}
											}
										}
										blstruct.musthave.mounts.mountname[mhMounts] = mountdetails{
											mounttype: mhMountType,
											address:   mhAddress,
											username:  mhUsername,
											pwd:       mhPassword,
											src:       mhMountSource,
											dest:      mhMountDest,
										}
									}
								}
							}
						case "Must-Not-Have":
							switch thirdStep.Key {
							case "Installed":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.mustnothave.installed = []string{""}
								} else {
									mnhInst := make([]string, len(thirdStep.Value.([]interface{})))
									mnhInstSlice := thirdStep.Value.([]interface{})
									for i, v := range mnhInstSlice {
										mnhInst[i] = v.(string)
									}
									blstruct.mustnothave.installed = mnhInst
								}
							case "Enabled":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.mustnothave.enabled = []string{""}
								} else {
									mnhEnabled := make([]string, len(thirdStep.Value.([]interface{})))
									mnhEnabledSlice := thirdStep.Value.([]interface{})
									for i, v := range mnhEnabledSlice {
										mnhEnabled[i] = v.(string)
									}
									blstruct.mustnothave.installed = mnhEnabled
								}
							case "Disabled":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.mustnothave.disabled = []string{""}
								} else {
									mnhDisabled := make([]string, len(thirdStep.Value.([]interface{})))
									mnhDisabledSlice := thirdStep.Value.([]interface{})
									for i, v := range mnhDisabledSlice {
										mnhDisabled[i] = v.(string)
									}
									blstruct.mustnothave.enabled = mnhDisabled
								}
							case "Users":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.mustnothave.users = []string{""}
								} else {
									mnhUsers := make([]string, len(thirdStep.Value.([]interface{})))
									mnhUsersSlice := thirdStep.Value.([]interface{})
									for i, v := range mnhUsersSlice {
										mnhUsers[i] = v.(string)
									}
									blstruct.mustnothave.users = mnhUsers
								}
							case "Rules":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.mustnothave.rules.fwopen.ports = []string{""}
									blstruct.mustnothave.rules.fwopen.protocols = []string{""}
									blstruct.mustnothave.rules.fwclosed.ports = []string{""}
									blstruct.mustnothave.rules.fwclosed.protocols = []string{""}
									blstruct.mustnothave.rules.fwzones = []string{""}
								} else {
									for _, blItem = range nextblValues {
										rulesValues, ok := blItem.Value.(yaml.MapSlice)
										if !ok {
											warnings++
										}
										switch blItem.Key {
										case "Open":
											for _, rulesItems := range rulesValues {
												switch rulesItems.Key {
												case "Ports":
													if rulesItems.Value == nil {
														datawarnings++
														blstruct.mustnothave.rules.fwopen.ports = []string{""}
													} else {
														mhRulesOPorts := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesOPortsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesOPortsSlice {
															mhRulesOPorts[i] = strconv.Itoa(v.(int))
														}
														blstruct.mustnothave.rules.fwopen.ports = mhRulesOPorts
													}
												case "Protocols":
													if rulesItems.Value == nil {
														datawarnings++
														blstruct.mustnothave.rules.fwopen.protocols = []string{""}
													} else {
														mhRulesOProtocols := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesOProtocolsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesOProtocolsSlice {
															mhRulesOProtocols[i] = v.(string)
														}
														blstruct.mustnothave.rules.fwopen.protocols = mhRulesOProtocols
													}
												}
											}
										case "Closed":
											for _, rulesItems := range rulesValues {
												switch rulesItems.Key {
												case "Ports":
													if rulesItems.Value == nil {
														datawarnings++
														blstruct.mustnothave.rules.fwclosed.ports = []string{""}
													} else {
														mhRulesCPorts := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesCPortsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesCPortsSlice {
															mhRulesCPorts[i] = strconv.Itoa(v.(int))
														}
														blstruct.mustnothave.rules.fwclosed.ports = mhRulesCPorts
													}
												case "Protocols":
													if rulesItems.Value == nil {
														datawarnings++
														blstruct.mustnothave.rules.fwclosed.protocols = []string{""}
													} else {
														mhRulesCProtocols := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesCProtocolsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesCProtocolsSlice {
															mhRulesCProtocols[i] = v.(string)
														}
														blstruct.mustnothave.rules.fwclosed.protocols = mhRulesCProtocols
													}
												}
											}
										case "Zones":
											if blItem.Value == nil {
												datawarnings++
												blstruct.mustnothave.rules.fwzones = []string{""}
											} else {
												mhRulesZones := make([]string, len(blItem.Value.([]interface{})))
												mhRulesZonesSlice := blItem.Value.([]interface{})
												for i, v := range mhRulesZonesSlice {
													mhRulesZones[i] = v.(string)
												}
												blstruct.mustnothave.rules.fwzones = mhRulesZones
											}
										}
									}
								}
							case "Mounts":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.mustnothave.mounts = []string{""}
								} else {
									mnhMounts := make([]string, len(thirdStep.Value.([]interface{})))
									mnhMountsSlice := thirdStep.Value.([]interface{})
									for i, v := range mnhMountsSlice {
										mnhMounts[i] = v.(string)
									}
									blstruct.mustnothave.mounts = mnhMounts
								}
							}
						case "Final":
							switch thirdStep.Key {
							case "Scripts":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.final.scripts = []string{""}
								} else {
									fnlScripts := make([]string, len(thirdStep.Value.([]interface{})))
									fnlScriptsSlice := thirdStep.Value.([]interface{})
									for i, v := range fnlScriptsSlice {
										fnlScripts[i] = v.(string)
									}
									blstruct.final.scripts = fnlScripts
								}
							case "Commands":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.final.commands = []string{""}
								} else {
									fnlCommands := make([]string, len(thirdStep.Value.([]interface{})))
									fnlCommandsSlice := thirdStep.Value.([]interface{})
									for i, v := range fnlCommandsSlice {
										fnlCommands[i] = v.(string)
									}
									blstruct.final.commands = fnlCommands
								}
							case "Collect":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.final.collect.logs = []string{""}
									blstruct.final.collect.stats = []string{""}
									blstruct.final.collect.files = []string{""}
									blstruct.final.collect.users = false
								} else {
									for _, blItem = range nextblValues {
										switch blItem.Key {
										case "Logs":
											if blItem.Value == nil {
												datawarnings++
												blstruct.final.collect.logs = []string{""}
											} else {
												fnlColLog := make([]string, len(blItem.Value.([]interface{})))
												fnlColLogSlice := blItem.Value.([]interface{})
												for i, v := range fnlColLogSlice {
													fnlColLog[i] = v.(string)
												}
												blstruct.final.collect.logs = fnlColLog
											}
										case "Stats":
											if blItem.Value == nil {
												datawarnings++
												blstruct.final.collect.stats = []string{""}
											} else {
												fnlColstats := make([]string, len(blItem.Value.([]interface{})))
												fnlColstatsSlice := blItem.Value.([]interface{})
												for i, v := range fnlColstatsSlice {
													fnlColstats[i] = v.(string)
												}
												blstruct.final.collect.stats = fnlColstats
											}
										case "Files":
											if blItem.Value == nil {
												datawarnings++
												blstruct.final.collect.files = []string{""}
											} else {
												fnlColFiles := make([]string, len(blItem.Value.([]interface{})))
												fnlColFilesSlice := blItem.Value.([]interface{})
												for i, v := range fnlColFilesSlice {
													fnlColFiles[i] = v.(string)
												}
												blstruct.final.collect.files = fnlColFiles
											}
										case "Users":
											if blItem.Value == nil {
												datawarnings++
												blstruct.final.collect.users = false
											} else {
												blstruct.final.collect.users = blItem.Value.(bool)
											}
										}
									}
								}
							case "Restart":
								if thirdStep.Value == nil {
									datawarnings++
									blstruct.final.restart.services = false
									blstruct.final.restart.servers = false
								} else {
									for _, blItem = range nextblValues {
										switch blItem.Key {
										case "Services":
											if blItem.Value == nil {
												datawarnings++
												blstruct.final.restart.services = false
											} else {
												blstruct.final.restart.services = blItem.Value.(bool)
											}
										case "Servers":
											if blItem.Value == nil {
												datawarnings++
												blstruct.final.restart.servers = false
											} else {
												blstruct.final.restart.servers = blItem.Value.(bool)
											}
										}
									}
								}
							}
						default:
							blerrors++
							panic("\nError parsing baseline.\n Please check your baseline.\nAborting...\n")
						}
					}
				}
			}

			// TODO apply baseline
			sshList := blstruct.applyOSExcludes(servergroupname, configs)
			//fmt.Println(sshList)
			// establish ssh connections to servers via goroutines and maintain sessions
			//commandChannel := make(chan map[string]string)
			var rebootBool bool
			blstruct.applyPrereq(&sshList)
			// commandset via channel to servers and wait for it to complete
			blstruct.applyMustHaves(&sshList, &rebootBool)
			// commandset via channel to servers and wait for it to complete
			blstruct.applyMustNotHaves(&sshList)
			// commandset via channel to servers and wait for it to complete
			blstruct.applyFinals(&sshList, &rebootBool)
			// commandset via channel to servers and wait for it to complete
			// once all checks are completed pass disconnect via channels to open sessions
			if rebootBool {
				// send reboot command to servers
			} else {
				// close connections without rebooting
			}
		}
	}
}