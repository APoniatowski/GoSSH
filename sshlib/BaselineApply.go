package sshlib

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/APoniatowski/GoSSH/channelreaderlib"
	"github.com/APoniatowski/GoSSH/loggerlib"
	"golang.org/x/crypto/ssh"

	//"golang.org/x/crypto/ssh"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// ApplyBaselines Apply defined baselines
func ApplyBaselines(baselineYAML *yaml.MapSlice, configs *yaml.MapSlice) {
	var warnings int
	var mainCategoryWarnings int
	var dataWarnings int
	var baselineErrors int
	// Baseline
	for _, blItem := range *baselineYAML {
		fmt.Printf("%s:\n", blItem.Key)
		groupValues, ok := blItem.Value.(yaml.MapSlice)
		if !ok {
			panic("\nCheck your baseline for issues\nAlternatively generate a template to see what is missing/wrong\n")
		}
		// Server groups
		for _, groupItem := range groupValues {
			// initialize the data
			serverGroupName := groupItem.Key.(string)
			var baselineStruct ParsedBaseline
			baselineStruct.musthave.configured.services = make(map[string]musthaveconfiguredservices)
			baselineStruct.musthave.users.users = make(map[string]musthaveusersstruct)
			baselineStruct.musthave.mounts.mountname = make(map[string]mountdetails)
			baselineStepsValue, ok := groupItem.Value.(yaml.MapSlice)
			if !ok {
				baselineErrors++
				panic("\nError parsing server groups.\nAborting...\n")
			}
			if strings.ToLower(serverGroupName) == "all" {
				fmt.Println("Applying baseline on all servers:")
			} else {
				fmt.Println("Applying baseline on", serverGroupName+":")
			}
			// Exclude, Prerequisites, Must-Have, Must-Not-Have, Final
			for _, baselineStepItem := range baselineStepsValue {
				nextValues, ok := baselineStepItem.Value.(yaml.MapSlice)
				if !ok {
					mainCategoryWarnings++
				}
				baselineStepCheck := baselineStepItem.Key
				if baselineStepItem.Key == nil {
					baselineErrors++
				}

				// OS, Servers, Tools, Files, VCS, etc
				for _, thirdStep := range nextValues {
					nextblValues, ok := thirdStep.Value.(yaml.MapSlice)
					if !ok {
						warnings++
					}
					if thirdStep.Key == nil {
						warnings++
						baselineErrors++
					} else {
						switch baselineStepCheck {
						case "Exclude":
							switch thirdStep.Key {
							case "OS":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.exclude.osExcl = []string{""}
								} else {
									exclOS := make([]string, len(thirdStep.Value.([]interface{})))
									OSslice := thirdStep.Value.([]interface{})
									for i, v := range OSslice {
										exclOS[i] = v.(string)
									}
									baselineStruct.exclude.osExcl = exclOS
								}
							case "Servers":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.exclude.serversExcl = []string{""}
								} else {
									exclServers := make([]string, len(thirdStep.Value.([]interface{})))
									serverSlice := thirdStep.Value.([]interface{})
									for i, v := range serverSlice {
										exclServers[i] = v.(string)
									}
									baselineStruct.exclude.serversExcl = exclServers
								}
							}
						case "Prerequisites":
							switch thirdStep.Key {
							case "Tools":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.prereq.tools = []string{""}
								} else {
									prereqTools := make([]string, len(thirdStep.Value.([]interface{})))
									prereqToolsSlice := thirdStep.Value.([]interface{})
									for i, v := range prereqToolsSlice {
										prereqTools[i] = v.(string)
									}
									baselineStruct.prereq.tools = prereqTools
								}
							case "Files": // for loop
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.prereq.files.urls = []string{""}
									baselineStruct.prereq.files.local.src = ""
									baselineStruct.prereq.files.local.dest = ""
									baselineStruct.prereq.files.remote.mounttype = ""
									baselineStruct.prereq.files.remote.address = ""
									baselineStruct.prereq.files.remote.username = ""
									baselineStruct.prereq.files.remote.pwd = ""
									baselineStruct.prereq.files.remote.src = ""
									baselineStruct.prereq.files.remote.dest = ""
								} else {
									for _, blItem = range nextblValues {
										switch blItem.Key {
										case "URLs":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.prereq.files.urls = []string{""}
											} else {
												prereqURLs := make([]string, len(blItem.Value.([]interface{})))
												prereqURLsSlice := blItem.Value.([]interface{})
												for i, v := range prereqURLsSlice {
													prereqURLs[i] = v.(string)
												}
												baselineStruct.prereq.files.urls = prereqURLs
											}
										case "Local":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.prereq.files.local.src = ""
												baselineStruct.prereq.files.local.dest = ""
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
															dataWarnings++
															baselineStruct.prereq.files.local.src = ""
														} else {
															baselineStruct.prereq.files.local.src = nextblStep.Value.(string)
														}
													case "Destination":
														if nextblStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.local.dest = ""
														} else {
															baselineStruct.prereq.files.local.dest = nextblStep.Value.(string)
														}
													}
												}
											}
										case "Remote":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.prereq.files.remote.mounttype = ""
												baselineStruct.prereq.files.remote.address = ""
												baselineStruct.prereq.files.remote.username = ""
												baselineStruct.prereq.files.remote.pwd = ""
												baselineStruct.prereq.files.remote.src = ""
												baselineStruct.prereq.files.remote.dest = ""
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
															dataWarnings++
															baselineStruct.prereq.files.remote.mounttype = ""
														} else {
															baselineStruct.prereq.files.remote.mounttype = nextblStep.Value.(string)
														}
													case "Address":
														if nextblStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.address = ""
														} else {
															baselineStruct.prereq.files.remote.address = nextblStep.Value.(string)
														}
													case "Username":
														if nextblStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.username = ""
														} else {
															baselineStruct.prereq.files.remote.username = nextblStep.Value.(string)
														}
													case "Password":
														if nextblStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.pwd = ""
														} else {
															baselineStruct.prereq.files.remote.pwd = nextblStep.Value.(string)
														}
													case "Source":
														if nextblStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.src = ""
														} else {
															baselineStruct.prereq.files.remote.src = nextblStep.Value.(string)
														}
													case "Destination":
														if nextblStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.dest = ""
														} else {
															baselineStruct.prereq.files.remote.dest = nextblStep.Value.(string)
														}
													case "Files":
														if nextblStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.files = []string{""}
														} else {
															remoteFiles := make([]string, len(nextblStep.Value.([]interface{})))
															remoteFilesSlice := nextblStep.Value.([]interface{})
															for i, v := range remoteFilesSlice {
																remoteFiles[i] = v.(string)
															}
															baselineStruct.prereq.files.remote.files = remoteFiles
														}
													}
												}
											}
										}
									}
								}
							case "VCS":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.prereq.vcs.urls = []string{""}
									baselineStruct.prereq.vcs.execute = []string{""}
								} else {
									for _, blItem = range nextblValues {
										switch blItem.Key {
										case "URLs":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.prereq.vcs.urls = []string{""}
											} else {
												vcsURLs := make([]string, len(blItem.Value.([]interface{})))
												vcsURLsSlice := blItem.Value.([]interface{})
												for i, v := range vcsURLsSlice {
													vcsURLs[i] = v.(string)
												}
												baselineStruct.prereq.vcs.urls = vcsURLs
											}
										case "Execute":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.prereq.vcs.execute = []string{""}
											} else {
												vcsCMDS := make([]string, len(blItem.Value.([]interface{})))
												vcsCMDSSlice := blItem.Value.([]interface{})
												for i, v := range vcsCMDSSlice {
													vcsCMDS[i] = v.(string)
												}
												baselineStruct.prereq.vcs.execute = vcsCMDS
											}
										}
									}
								}
							case "Script":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.prereq.script = ""
								} else {
									baselineStruct.prereq.script = thirdStep.Value.(string)
								}
							case "Clean-up":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.prereq.cleanup = false
								} else {
									baselineStruct.prereq.cleanup = thirdStep.Value.(bool)
								}
							}
						case "Must-Have":
							switch thirdStep.Key {
							case "Installed":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.musthave.installed = []string{""}
								} else {
									mhInst := make([]string, len(thirdStep.Value.([]interface{})))
									mhInstSlice := thirdStep.Value.([]interface{})
									for i, v := range mhInstSlice {
										mhInst[i] = v.(string)
									}
									baselineStruct.musthave.installed = mhInst
								}
							case "Enabled":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.musthave.enabled = []string{""}
								} else {
									mhEnabled := make([]string, len(thirdStep.Value.([]interface{})))
									mhEnabledSlice := thirdStep.Value.([]interface{})
									for i, v := range mhEnabledSlice {
										mhEnabled[i] = v.(string)
									}
									baselineStruct.musthave.enabled = mhEnabled
								}
							case "Disabled":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.musthave.disabled = []string{""}
								} else {
									mhDisabled := make([]string, len(thirdStep.Value.([]interface{})))
									mhDisabledSlice := thirdStep.Value.([]interface{})
									for i, v := range mhDisabledSlice {
										mhDisabled[i] = v.(string)
									}
									baselineStruct.musthave.enabled = mhDisabled
								}
							case "Configured":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.musthave.configured.services[""] = musthaveconfiguredservices{source: []string{""}, destination: []string{""}}
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
													dataWarnings++
													baselineStruct.musthave.configured.services[service] = musthaveconfiguredservices{source: []string{""}}
												} else {
													mhConfSrc = make([]string, len(confItems.Value.([]interface{})))
													mhConfSrcSlice := confItems.Value.([]interface{})
													for i, v := range mhConfSrcSlice {
														mhConfSrc[i] = v.(string)
													}
												}
											case "Destination":
												if confItems.Value == nil {
													dataWarnings++
													baselineStruct.musthave.configured.services[service] = musthaveconfiguredservices{destination: []string{""}}
												} else {
													mhConfDest = make([]string, len(confItems.Value.([]interface{})))
													mhConfDestSlice := confItems.Value.([]interface{})
													for i, v := range mhConfDestSlice {
														mhConfDest[i] = v.(string)
													}
												}
											}
											baselineStruct.musthave.configured.services[service] = musthaveconfiguredservices{source: mhConfSrc, destination: mhConfDest}
										}
									}
								}
							case "Users":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.musthave.users.users[""] = musthaveusersstruct{groups: []string{""}, shell: "", home: "", sudoer: false}
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
													dataWarnings++
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
													dataWarnings++
													mhUsershell = ""
												} else {
													mhUsershell = userItems.Value.(string)
												}
											case "Home-Dir":
												if userItems.Value == nil {
													dataWarnings++
													mhUserhome = ""
												} else {
													mhUserhome = userItems.Value.(string)
												}
											case "Sudoer":
												if userItems.Value == nil {
													dataWarnings++
													mhUsersudo = false
												} else {
													mhUsersudo = userItems.Value.(bool)
												}
											}
											baselineStruct.musthave.users.users[user] = musthaveusersstruct{
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
									dataWarnings++
									baselineStruct.musthave.policies.polimport = ""
									baselineStruct.musthave.policies.polstatus = ""
									baselineStruct.musthave.policies.polreboot = false
								} else {
									for _, blItem = range nextblValues {
										switch blItem.Key {
										case "Status":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.musthave.policies.polstatus = ""
											} else {
												baselineStruct.musthave.policies.polstatus = blItem.Value.(string)
											}
										case "Import":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.musthave.policies.polimport = ""
											} else {
												baselineStruct.musthave.policies.polimport = blItem.Value.(string)
											}
										case "Reboot":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.musthave.policies.polreboot = false
											} else {
												baselineStruct.musthave.policies.polreboot = blItem.Value.(bool)
											}
										}
									}
								}
							case "Rules":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.musthave.rules.fwopen.ports = []string{""}
									baselineStruct.musthave.rules.fwopen.protocols = []string{""}
									baselineStruct.musthave.rules.fwclosed.ports = []string{""}
									baselineStruct.musthave.rules.fwclosed.protocols = []string{""}
									baselineStruct.musthave.rules.fwzones = []string{""}
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
														dataWarnings++
														baselineStruct.musthave.rules.fwopen.ports = []string{""}
													} else {
														mhRulesOPorts := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesOPortsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesOPortsSlice {
															mhRulesOPorts[i] = strconv.Itoa(v.(int))
														}
														baselineStruct.musthave.rules.fwopen.ports = mhRulesOPorts
													}
												case "Protocols":
													if rulesItems.Value == nil {
														dataWarnings++
														baselineStruct.musthave.rules.fwopen.ports = []string{""}
													} else {
														mhRulesOProtocols := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesOProtocolsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesOProtocolsSlice {
															mhRulesOProtocols[i] = v.(string)
														}
														baselineStruct.musthave.rules.fwopen.protocols = mhRulesOProtocols
													}
												}
											}
										case "Closed":
											for _, rulesItems := range rulesValues {
												switch rulesItems.Key {
												case "Ports":
													if rulesItems.Value == nil {
														dataWarnings++
														baselineStruct.musthave.rules.fwclosed.ports = []string{""}
													} else {
														mhRulesCPorts := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesCPortsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesCPortsSlice {
															mhRulesCPorts[i] = strconv.Itoa(v.(int))
														}
														baselineStruct.musthave.rules.fwclosed.ports = mhRulesCPorts
													}
												case "Protocols":
													if rulesItems.Value == nil {
														dataWarnings++
														baselineStruct.musthave.rules.fwclosed.protocols = []string{""}
													} else {
														mhRulesCProtocols := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesCProtocolsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesCProtocolsSlice {
															mhRulesCProtocols[i] = v.(string)
														}
														baselineStruct.musthave.rules.fwclosed.protocols = mhRulesCProtocols
													}
												}
											}
										case "Zones":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.musthave.rules.fwzones = []string{""}
											} else {
												mhRulesZones := make([]string, len(blItem.Value.([]interface{})))
												mhRulesZonesSlice := blItem.Value.([]interface{})
												for i, v := range mhRulesZonesSlice {
													mhRulesZones[i] = v.(string)
												}
												baselineStruct.musthave.rules.fwzones = mhRulesZones
											}
										}
									}
								}
							case "Mounts":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.musthave.mounts.mountname[""] = mountdetails{
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
													dataWarnings++
													mhMountType = ""
												} else {
													mhMountType = mountItems.Value.(string)
												}
											case "Address":
												if mountItems.Value == nil {
													dataWarnings++
													mhAddress = ""
												} else {
													mhAddress = mountItems.Value.(string)
												}
											case "Username":
												if mountItems.Value == nil {
													dataWarnings++
													mhUsername = ""
												} else {
													mhUsername = mountItems.Value.(string)
												}
											case "Password":
												if mountItems.Value == nil {
													dataWarnings++
													mhPassword = ""
												} else {
													mhPassword = mountItems.Value.(string)
												}
											case "Source":
												if mountItems.Value == nil {
													dataWarnings++
													mhMountSource = ""
												} else {
													mhMountSource = mountItems.Value.(string)
												}
											case "Destination":
												if mountItems.Value == nil {
													dataWarnings++
													mhMountDest = ""
												} else {
													mhMountDest = mountItems.Value.(string)
												}
											}
										}
										baselineStruct.musthave.mounts.mountname[mhMounts] = mountdetails{
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
									dataWarnings++
									baselineStruct.mustnothave.installed = []string{""}
								} else {
									mnhInst := make([]string, len(thirdStep.Value.([]interface{})))
									mnhInstSlice := thirdStep.Value.([]interface{})
									for i, v := range mnhInstSlice {
										mnhInst[i] = v.(string)
									}
									baselineStruct.mustnothave.installed = mnhInst
								}
							case "Enabled":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.mustnothave.enabled = []string{""}
								} else {
									mnhEnabled := make([]string, len(thirdStep.Value.([]interface{})))
									mnhEnabledSlice := thirdStep.Value.([]interface{})
									for i, v := range mnhEnabledSlice {
										mnhEnabled[i] = v.(string)
									}
									baselineStruct.mustnothave.installed = mnhEnabled
								}
							case "Disabled":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.mustnothave.disabled = []string{""}
								} else {
									mnhDisabled := make([]string, len(thirdStep.Value.([]interface{})))
									mnhDisabledSlice := thirdStep.Value.([]interface{})
									for i, v := range mnhDisabledSlice {
										mnhDisabled[i] = v.(string)
									}
									baselineStruct.mustnothave.enabled = mnhDisabled
								}
							case "Users":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.mustnothave.users = []string{""}
								} else {
									mnhUsers := make([]string, len(thirdStep.Value.([]interface{})))
									mnhUsersSlice := thirdStep.Value.([]interface{})
									for i, v := range mnhUsersSlice {
										mnhUsers[i] = v.(string)
									}
									baselineStruct.mustnothave.users = mnhUsers
								}
							case "Rules":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.mustnothave.rules.fwopen.ports = []string{""}
									baselineStruct.mustnothave.rules.fwopen.protocols = []string{""}
									baselineStruct.mustnothave.rules.fwclosed.ports = []string{""}
									baselineStruct.mustnothave.rules.fwclosed.protocols = []string{""}
									baselineStruct.mustnothave.rules.fwzones = []string{""}
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
														dataWarnings++
														baselineStruct.mustnothave.rules.fwopen.ports = []string{""}
													} else {
														mhRulesOPorts := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesOPortsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesOPortsSlice {
															mhRulesOPorts[i] = strconv.Itoa(v.(int))
														}
														baselineStruct.mustnothave.rules.fwopen.ports = mhRulesOPorts
													}
												case "Protocols":
													if rulesItems.Value == nil {
														dataWarnings++
														baselineStruct.mustnothave.rules.fwopen.protocols = []string{""}
													} else {
														mhRulesOProtocols := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesOProtocolsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesOProtocolsSlice {
															mhRulesOProtocols[i] = v.(string)
														}
														baselineStruct.mustnothave.rules.fwopen.protocols = mhRulesOProtocols
													}
												}
											}
										case "Closed":
											for _, rulesItems := range rulesValues {
												switch rulesItems.Key {
												case "Ports":
													if rulesItems.Value == nil {
														dataWarnings++
														baselineStruct.mustnothave.rules.fwclosed.ports = []string{""}
													} else {
														mhRulesCPorts := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesCPortsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesCPortsSlice {
															mhRulesCPorts[i] = strconv.Itoa(v.(int))
														}
														baselineStruct.mustnothave.rules.fwclosed.ports = mhRulesCPorts
													}
												case "Protocols":
													if rulesItems.Value == nil {
														dataWarnings++
														baselineStruct.mustnothave.rules.fwclosed.protocols = []string{""}
													} else {
														mhRulesCProtocols := make([]string, len(rulesItems.Value.([]interface{})))
														mhRulesCProtocolsSlice := rulesItems.Value.([]interface{})
														for i, v := range mhRulesCProtocolsSlice {
															mhRulesCProtocols[i] = v.(string)
														}
														baselineStruct.mustnothave.rules.fwclosed.protocols = mhRulesCProtocols
													}
												}
											}
										case "Zones":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.mustnothave.rules.fwzones = []string{""}
											} else {
												mhRulesZones := make([]string, len(blItem.Value.([]interface{})))
												mhRulesZonesSlice := blItem.Value.([]interface{})
												for i, v := range mhRulesZonesSlice {
													mhRulesZones[i] = v.(string)
												}
												baselineStruct.mustnothave.rules.fwzones = mhRulesZones
											}
										}
									}
								}
							case "Mounts":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.mustnothave.mounts = []string{""}
								} else {
									mnhMounts := make([]string, len(thirdStep.Value.([]interface{})))
									mnhMountsSlice := thirdStep.Value.([]interface{})
									for i, v := range mnhMountsSlice {
										mnhMounts[i] = v.(string)
									}
									baselineStruct.mustnothave.mounts = mnhMounts
								}
							}
						case "Final":
							switch thirdStep.Key {
							case "Scripts":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.final.scripts = []string{""}
								} else {
									fnlScripts := make([]string, len(thirdStep.Value.([]interface{})))
									fnlScriptsSlice := thirdStep.Value.([]interface{})
									for i, v := range fnlScriptsSlice {
										fnlScripts[i] = v.(string)
									}
									baselineStruct.final.scripts = fnlScripts
								}
							case "Commands":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.final.commands = []string{""}
								} else {
									fnlCommands := make([]string, len(thirdStep.Value.([]interface{})))
									fnlCommandsSlice := thirdStep.Value.([]interface{})
									for i, v := range fnlCommandsSlice {
										fnlCommands[i] = v.(string)
									}
									baselineStruct.final.commands = fnlCommands
								}
							case "Collect":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.final.collect.logs = []string{""}
									baselineStruct.final.collect.stats = []string{""}
									baselineStruct.final.collect.files = []string{""}
									baselineStruct.final.collect.users = false
								} else {
									for _, blItem = range nextblValues {
										switch blItem.Key {
										case "Logs":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.final.collect.logs = []string{""}
											} else {
												fnlColLog := make([]string, len(blItem.Value.([]interface{})))
												fnlColLogSlice := blItem.Value.([]interface{})
												for i, v := range fnlColLogSlice {
													fnlColLog[i] = v.(string)
												}
												baselineStruct.final.collect.logs = fnlColLog
											}
										case "Stats":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.final.collect.stats = []string{""}
											} else {
												fnlColstats := make([]string, len(blItem.Value.([]interface{})))
												fnlColstatsSlice := blItem.Value.([]interface{})
												for i, v := range fnlColstatsSlice {
													fnlColstats[i] = v.(string)
												}
												baselineStruct.final.collect.stats = fnlColstats
											}
										case "Files":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.final.collect.files = []string{""}
											} else {
												fnlColFiles := make([]string, len(blItem.Value.([]interface{})))
												fnlColFilesSlice := blItem.Value.([]interface{})
												for i, v := range fnlColFilesSlice {
													fnlColFiles[i] = v.(string)
												}
												baselineStruct.final.collect.files = fnlColFiles
											}
										case "Users":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.final.collect.users = false
											} else {
												baselineStruct.final.collect.users = blItem.Value.(bool)
											}
										}
									}
								}
							case "Restart":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.final.restart.services = false
									baselineStruct.final.restart.servers = false
								} else {
									for _, blItem = range nextblValues {
										switch blItem.Key {
										case "Services":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.final.restart.services = false
											} else {
												baselineStruct.final.restart.services = blItem.Value.(bool)
											}
										case "Servers":
											if blItem.Value == nil {
												dataWarnings++
												baselineStruct.final.restart.servers = false
											} else {
												baselineStruct.final.restart.servers = blItem.Value.(bool)
											}
										}
									}
								}
							}
						default:
							baselineErrors++
							panic("\nError parsing baseline.\n Please check your baseline.\nAborting...\n")
						}
					}
				}
			}

			// TODO apply baseline
			sshList := baselineStruct.applyOSExcludes(serverGroupName, configs)
			disconnectSessions := make(chan bool)
			var rebootBool bool
			//fmt.Println(sshList)
			// establish ssh connections to servers via goroutines and maintain sessions
			for _, groupItem := range *configs {
				if serverGroupName == groupItem.Key {
					output := make(chan string)
					readyState := make(chan []bool)
					var wg sync.WaitGroup

					groupValue, ok := groupItem.Value.(yaml.MapSlice)
					if !ok {
						panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
					}
					for _, serverItem := range groupValue {
						servername := serverItem.Key
						serverValue, ok := serverItem.Value.(yaml.MapSlice)
						if !ok {
							panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
						}
						var pp ParsedPool
						pp.fqdn = serverValue[0].Value
						pp.username = serverValue[1].Value
						pp.password = serverValue[2].Value
						pp.keypath = serverValue[3].Value
						pp.port = serverValue[4].Value
						pp.os = serverValue[5].Value
						pp.defaulter()
						commandChannel := make(chan map[string]string)
						for i := range sshList {
							if i == pp.fqdn.(string) {
								wg.Add(1)
								fmt.Println("sending false to disconnectsessions channel")
								disconnectSessions <- false // To keep sessions alive and not disconnect
								fmt.Println("starting ssh goroutine")
								go pp.connectAndRunBaseline(commandChannel,
									servername.(string),
									output,
									disconnectSessions,
									readyState,
									&wg)
							}
							fmt.Println("started goroutines, starting output goroutine")
							go func() {
								wg.Wait()
								close(output)
							}()

						}
						fmt.Println("readystate/sshlist")
						fmt.Printf("%d  %d\n", len(readyState), len(sshList))
						for {
							if len(readyState) == len(sshList) {
								break
							}
						}
						fmt.Println("starting stages...")
						baselineStruct.applyPrereq(&sshList, commandChannel)
						baselineStruct.applyMustHaves(&sshList, &rebootBool)
						baselineStruct.applyMustNotHaves(&sshList)
						baselineStruct.applyFinals(&sshList, &rebootBool)

						channelreaderlib.ChannelReaderGroups(output, &wg) // move to each apply stage
					}
				} else if strings.ToLower(serverGroupName) == "all" {
					var allServers yaml.MapSlice
					output := make(chan string)
					readyState := make(chan []bool)
					var wg sync.WaitGroup
					// Concatenates the groups to create a single group
					for _, groupItem := range *configs {
						groupValue, ok := groupItem.Value.(yaml.MapSlice)
						if !ok {
							panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
						}

						allServers = append(allServers, groupValue...)
					}
					for _, serverItem := range allServers {
						wg.Add(1)
						servername := serverItem.Key
						serverValue, ok := serverItem.Value.(yaml.MapSlice)
						if !ok {
							panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
						}
						var pp ParsedPool
						pp.fqdn = serverValue[0].Value
						pp.username = serverValue[1].Value
						pp.password = serverValue[2].Value
						pp.keypath = serverValue[3].Value
						pp.port = serverValue[4].Value
						pp.os = serverValue[5].Value
						pp.defaulter()
						commandChannel := make(chan map[string]string)
						for i := range sshList {
							if i == pp.fqdn.(string) {
								wg.Add(1)
								disconnectSessions <- false // To keep sessions alive and not disconnect

								go pp.connectAndRunBaseline(commandChannel,
									servername.(string),
									output,
									disconnectSessions,
									readyState,
									&wg)
							}
							go func() {
								wg.Wait()
								close(output)
							}()

						}
						for {
							if len(readyState) == len(sshList) {
								break
							}
						}
						baselineStruct.applyPrereq(&sshList, commandChannel)
						baselineStruct.applyMustHaves(&sshList, &rebootBool)
						baselineStruct.applyMustNotHaves(&sshList)
						baselineStruct.applyFinals(&sshList, &rebootBool)

						channelreaderlib.ChannelReaderGroups(output, &wg) // move to each apply stage
					}
				}
			}

			if rebootBool {
				// send reboot command to servers
			} else {
				// close connections without rebooting
			}
			disconnectSessions <- true // When read, channels/sessions will be closed
		}
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// connectAndRun Establish a connection and run command(s), will add CLI args in the near future
func (parseddata *ParsedPool) connectAndRunBaseline(command chan map[string]string,
	servername string,
	output chan<- string,
	disconnect <-chan bool,
	ready chan []bool,
	wg *sync.WaitGroup) {
	pp := parseddata
	authMethodCheck := []ssh.AuthMethod{}
	key, err := ioutil.ReadFile(pp.keypath.(string))
	if err != nil {
		authMethodCheck = append(authMethodCheck, ssh.Password(pp.password.(string)))
	} else {
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			loggerlib.GeneralError(servername, "[INFO: Failed To Parse PrivKey] ", err)
		}
		authMethodCheck = append(authMethodCheck, ssh.PublicKeys(signer))
	}
	// hostKeyCallback, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	// if err != nil {
	hostKeyCallback := ssh.InsecureIgnoreHostKey()
	// }
	sshConfig := &ssh.ClientConfig{
		User:            pp.username.(string),
		Auth:            authMethodCheck,
		HostKeyCallback: hostKeyCallback,
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
		},
		Timeout: 5 * time.Second,
	}
	defer func() {
		if recv := recover(); recv != nil {
			recoveries = recv
		}
	}()
	connection, err := ssh.Dial("tcp", pp.fqdn.(string)+":"+pp.port.(string), sshConfig)
	if err != nil {
		loggerlib.GeneralError(servername, "[ERROR: Connection Failed] ", err)
		validator = "NOK\n"
		output <- validator
		wg.Done()
	} else {
		yesReady := <-ready
		yesReady = append(yesReady, true)
		ready <- yesReady
		for {
			commandCheck := <-command
			if len(commandCheck) != 0 {
				disconnectCheck := <-disconnect
				if !disconnectCheck {
					for i, j := range commandCheck {
						if i == pp.fqdn && j != "" {
							output <- executeBaselines(servername, j, pp.password.(string), connection)
							commandCheck[i] = ""
						}
					}
					command <- commandCheck
				} else {
					break
				}
			}
		}
		defer connection.Close()
		defer wg.Done()
	}
}

func executeBaselines(servername string, cmd string, password string, connection *ssh.Client) string {
	// adding recover to avoid panics during a run. Logs are written, so no need to panic when it its one of
	// the errors below.
	defer func() {
		if recv := recover(); recv != nil {
			recoveries = recv
		}
	}()
	session, err := connection.NewSession()
	if err != nil {
		loggerlib.GeneralError(servername, "[ERROR: Failed To Create Session] ", err)
	}
	defer session.Close()
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 50, 100, modes); err != nil {
		err := session.Close()
		if err != nil {
			loggerlib.GeneralError(servername, "[ERROR: Closing Pty session] ", err)
		}
		loggerlib.GeneralError(servername, "[ERROR: Pty Request Failed] ", err)
	}
	in, err := session.StdinPipe()
	if err != nil {
		loggerlib.GeneralError(servername, "[ERROR: Stdin Error] ", err)
	}
	out, err := session.StdoutPipe()
	if err != nil {
		loggerlib.GeneralError(servername, "[ERROR: Stdout Error] ", err)
	}
	var terminaloutput []byte
	var waitoutput sync.WaitGroup
	// it does not wait for output on some machines that are taking too long to respond. I'd like to avoid using Rlocks/Runlocks for this
	waitoutput.Add(1)
	go func(in io.WriteCloser, out io.Reader, terminaloutput *[]byte) {
		var (
			line string
			read = bufio.NewReader(out)
		)
		for {
			buffer, err := read.ReadByte()
			if err != nil {
				break
			}
			*terminaloutput = append(*terminaloutput, buffer)
			if buffer == byte('\n') {
				line = ""
				continue
			}
			line += string(buffer)
			if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
				_, err = in.Write([]byte(password + "\n"))
				if err != nil {
					break
				}
			}
		}
		waitoutput.Done()
	}(in, out, &terminaloutput)
	_, err = session.Output(cmd)
	waitoutput.Wait()
	if err != nil {
		validator = "NOK\n"
		loggerlib.ErrorLogger(servername, "[INFO: Failed] ", terminaloutput)
	} else {
		validator = "OK\n"
		loggerlib.OutputLogger(servername, "[INFO: Success] ", terminaloutput)
	}
	return validator
}
