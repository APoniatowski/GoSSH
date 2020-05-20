package sshlib

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// ApplyBaselines Apply defined baselines
func ApplyBaselines(baselineyaml *yaml.MapSlice) {
	var warnings int
	var maincategorywarnings int
	var datawarnings int
	var blerrors int
	// Baseline
	for _, blItem := range *baselineyaml {
		fmt.Printf("%s:\n", blItem.Key)
		groupValues, ok := blItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("\nCheck your baseline for issues\nAlternatively generate a template to see what is missing/wrong\n"))
		}
		// Server groups
		for _, groupItem := range groupValues {
			// initialize the data
			servergroupname := groupItem.Key.(string)
			var blstruct ParsedBaseline
			blstruct.musthave.configured.services = make(map[string]musthaveconfiguredservices)
			blstruct.musthave.users.users = make(map[string]musthaveusersstruct)
			blstruct.musthave.mounts.mountname = make(map[string]mountdetails)
			//
			// fmt.Printf("%s:\n", groupItem.Key)
			blstepsValue, ok := groupItem.Value.(yaml.MapSlice)
			if !ok {
				blerrors++
				panic(fmt.Sprintf("\nError parsing server groups.\nAborting...\n"))
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
													fmt.Println("3 Error parsing baseline. Please check the baseline you specified or generate a template")
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
							panic(fmt.Sprintf("\nError parsing baseline.\n Please check your baseline.\nAborting...\n"))
						}
					}
				}
			}
			fmt.Println()
			// Start operations here, to run jobs per group
			// or in case if it is 'all', then it needs to run through all servers

			// for ke, ve := range blstruct.musthave.mounts.mountname {
			// 	fmt.Println(ke, ve)
			// }
			// for ke, ve := range blstruct.musthave.users.users {
			// 	fmt.Println(ke, ve)
			// }
			// for ke, ve := range blstruct.musthave.configured.services {
			// 	fmt.Println(ke, ve)
			// }
		}
	}
}

// CheckBaselines Baseline compliancy check. Provides info of if servers meet baseline compliance as defined in the chosen baseline file
func CheckBaselines(baselineyaml *yaml.MapSlice) {
	var warnings int
	var maincategorywarnings int
	var datawarnings int
	var blerrors int
	// Baseline
	for _, blItem := range *baselineyaml {
		fmt.Printf("%s:\n", blItem.Key)
		groupValues, ok := blItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("\nCheck your baseline for issues\nAlternatively generate a template to see what is missing/wrong\n"))
		}
		// Server groups
		for _, groupItem := range groupValues {
			// initialize the data
			servergroupname := groupItem.Key.(string)
			var blstruct ParsedBaseline
			blstruct.musthave.configured.services = make(map[string]musthaveconfiguredservices)
			blstruct.musthave.users.users = make(map[string]musthaveusersstruct)
			blstruct.musthave.mounts.mountname = make(map[string]mountdetails)
			//
			// fmt.Printf("%s:\n", groupItem.Key)
			blstepsValue, ok := groupItem.Value.(yaml.MapSlice)
			if !ok {
				blerrors++
				panic(fmt.Sprintf("\nError parsing server groups.\nAborting...\n"))
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
													fmt.Println("3 Error parsing baseline. Please check the baseline you specified or generate a template")
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
							panic(fmt.Sprintf("\nError parsing baseline.\n Please check your baseline.\nAborting...\n"))
						}
					}
				}
			}
			fmt.Println()
			// Start operations here, to run jobs per group
			// or in case if it is 'all', then it needs to run through all servers

			// for ke, ve := range blstruct.musthave.mounts.mountname {
			// 	fmt.Println(ke, ve)
			// }
			// for ke, ve := range blstruct.musthave.users.users {
			// 	fmt.Println(ke, ve)
			// }
			// for ke, ve := range blstruct.musthave.configured.services {
			// 	fmt.Println(ke, ve)
			// }
		}
	}
}

//VerifyBaselines Baseline verification. Check if baselines have no errors and provide steps, as to what will be done, if applied
func VerifyBaselines(baselineyaml *yaml.MapSlice) {
	var warnings int
	var maincategorywarnings int
	var datawarnings int
	var blerrors int
	// Baseline
	for _, blItem := range *baselineyaml {
		fmt.Printf("%s:\n", blItem.Key)
		groupValues, ok := blItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("\nError:\nCheck your baseline for issues\nAlternatively generate a template to see what is missing/wrong\n"))
		}
		// Server groups
		for _, groupItem := range groupValues {
			// initialize the data
			servergroupname := groupItem.Key.(string)
			var blstruct ParsedBaseline
			blstruct.musthave.configured.services = make(map[string]musthaveconfiguredservices)
			blstruct.musthave.users.users = make(map[string]musthaveusersstruct)
			blstruct.musthave.mounts.mountname = make(map[string]mountdetails)
			//
			// fmt.Printf("%s:\n", groupItem.Key)
			blstepsValue, ok := groupItem.Value.(yaml.MapSlice)
			if !ok {
				blerrors++
				panic(fmt.Sprintf("\nError parsing server groups.\nAborting...\n"))
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
							panic(fmt.Sprintf("\nError parsing baseline.\n Please check your baseline.\nAborting...\n"))
						}
					}
				}
			}
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
			if len(blstruct.musthave.installed) == 0 &&
				len(blstruct.musthave.enabled) == 0 &&
				len(blstruct.musthave.disabled) == 0 &&
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
				//will continue from here later
			}

			// Start operations here, to run jobs per group
			// or in case if it is 'all', then it needs to run through all servers

			// for ke, ve := range blstruct.musthave.mounts.mountname {
			// 	fmt.Println(ke, ve)
			// }
			// for ke, ve := range blstruct.musthave.users.users {
			// 	fmt.Println(ke, ve)
			// }
			// for ke, ve := range blstruct.musthave.configured.services {
			// 	fmt.Println(ke, ve)
			// }
		}
	}
}
