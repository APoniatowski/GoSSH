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
	for _, baselineItem := range *baselineYAML {
		fmt.Printf("%s:\n", baselineItem.Key)
		groupValues, ok := baselineItem.Value.(yaml.MapSlice)
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
					nextBaselineValues, ok := thirdStep.Value.(yaml.MapSlice)
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
									OsSlice := thirdStep.Value.([]interface{})
									for i, v := range OsSlice {
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
									prerequisiteTools := make([]string, len(thirdStep.Value.([]interface{})))
									prerequisiteToolsSlice := thirdStep.Value.([]interface{})
									for i, v := range prerequisiteToolsSlice {
										prerequisiteTools[i] = v.(string)
									}
									baselineStruct.prereq.tools = prerequisiteTools
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
									for _, baselineItem = range nextBaselineValues {
										switch baselineItem.Key {
										case "URLs":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.prereq.files.urls = []string{""}
											} else {
												prerequisiteURLs := make([]string, len(baselineItem.Value.([]interface{})))
												prerequisiteURLsSlice := baselineItem.Value.([]interface{})
												for i, v := range prerequisiteURLsSlice {
													prerequisiteURLs[i] = v.(string)
												}
												baselineStruct.prereq.files.urls = prerequisiteURLs
											}
										case "Local":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.prereq.files.local.src = ""
												baselineStruct.prereq.files.local.dest = ""
											} else {
												extraBaselineValues, ok := baselineItem.Value.(yaml.MapSlice)
												if !ok {
													fmt.Println("Error parsing baseline. Please check the baseline you specified or generate a template")
												}
												var nextBaselineStep yaml.MapItem
												for _, nextBaselineStep = range extraBaselineValues {
													switch nextBaselineStep.Key {
													case "Source":
														if nextBaselineStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.local.src = ""
														} else {
															baselineStruct.prereq.files.local.src = nextBaselineStep.Value.(string)
														}
													case "Destination":
														if nextBaselineStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.local.dest = ""
														} else {
															baselineStruct.prereq.files.local.dest = nextBaselineStep.Value.(string)
														}
													}
												}
											}
										case "Remote":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.prereq.files.remote.mounttype = ""
												baselineStruct.prereq.files.remote.address = ""
												baselineStruct.prereq.files.remote.username = ""
												baselineStruct.prereq.files.remote.pwd = ""
												baselineStruct.prereq.files.remote.src = ""
												baselineStruct.prereq.files.remote.dest = ""
											} else {
												extraBaselineValues, ok := baselineItem.Value.(yaml.MapSlice)
												if !ok {
													warnings++
												}
												var nextBaselineStep yaml.MapItem
												for _, nextBaselineStep = range extraBaselineValues {
													switch nextBaselineStep.Key {
													case "Type":
														if nextBaselineStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.mounttype = ""
														} else {
															baselineStruct.prereq.files.remote.mounttype = nextBaselineStep.Value.(string)
														}
													case "Address":
														if nextBaselineStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.address = ""
														} else {
															baselineStruct.prereq.files.remote.address = nextBaselineStep.Value.(string)
														}
													case "Username":
														if nextBaselineStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.username = ""
														} else {
															baselineStruct.prereq.files.remote.username = nextBaselineStep.Value.(string)
														}
													case "Password":
														if nextBaselineStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.pwd = ""
														} else {
															baselineStruct.prereq.files.remote.pwd = nextBaselineStep.Value.(string)
														}
													case "Source":
														if nextBaselineStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.src = ""
														} else {
															baselineStruct.prereq.files.remote.src = nextBaselineStep.Value.(string)
														}
													case "Destination":
														if nextBaselineStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.dest = ""
														} else {
															baselineStruct.prereq.files.remote.dest = nextBaselineStep.Value.(string)
														}
													case "Files":
														if nextBaselineStep.Value == nil {
															dataWarnings++
															baselineStruct.prereq.files.remote.files = []string{""}
														} else {
															remoteFiles := make([]string, len(nextBaselineStep.Value.([]interface{})))
															remoteFilesSlice := nextBaselineStep.Value.([]interface{})
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
									for _, baselineItem = range nextBaselineValues {
										switch baselineItem.Key {
										case "URLs":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.prereq.vcs.urls = []string{""}
											} else {
												vcsURLs := make([]string, len(baselineItem.Value.([]interface{})))
												vcsURLsSlice := baselineItem.Value.([]interface{})
												for i, v := range vcsURLsSlice {
													vcsURLs[i] = v.(string)
												}
												baselineStruct.prereq.vcs.urls = vcsURLs
											}
										case "Execute":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.prereq.vcs.execute = []string{""}
											} else {
												vcsCommands := make([]string, len(baselineItem.Value.([]interface{})))
												vcsCommandsSlice := baselineItem.Value.([]interface{})
												for i, v := range vcsCommandsSlice {
													vcsCommands[i] = v.(string)
												}
												baselineStruct.prereq.vcs.execute = vcsCommands
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
							case "Commands":
								if thirdStep.Value == nil {
									dataWarnings++
									baselineStruct.prereq.commands = []string{""}
								} else {
									prereqCommands := make([]string, len(thirdStep.Value.([]interface{})))
									prereqCommandsSlice := thirdStep.Value.([]interface{})
									for i, v := range prereqCommandsSlice {
										prereqCommands[i] = v.(string)
									}
									baselineStruct.prereq.commands = prereqCommands
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
									for _, baselineItem = range nextBaselineValues {
										service := baselineItem.Key.(string)
										var mhConfSrc []string
										var mhConfDest []string
										confValues, ok := baselineItem.Value.(yaml.MapSlice)
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
									for _, baselineItem = range nextBaselineValues {
										user := baselineItem.Key.(string)
										var mhUserGroup []string
										var mhUserShell string
										var mhUserHome string
										var mhUserSudo bool
										userValues, ok := baselineItem.Value.(yaml.MapSlice)
										if !ok {
											warnings++
										}
										for _, userItems := range userValues {
											switch userItems.Key {
											case "Groups":
												if userItems.Value == nil {
													dataWarnings++
													mhUserGroup = []string{""}
												} else {
													mhUserGroup = make([]string, len(userItems.Value.([]interface{})))
													mhUserGroupSlice := userItems.Value.([]interface{})
													for i, v := range mhUserGroupSlice {
														mhUserGroup[i] = v.(string)
													}
												}
											case "Shell":
												if userItems.Value == nil {
													dataWarnings++
													mhUserShell = ""
												} else {
													mhUserShell = userItems.Value.(string)
												}
											case "Home-Dir":
												if userItems.Value == nil {
													dataWarnings++
													mhUserHome = ""
												} else {
													mhUserHome = userItems.Value.(string)
												}
											case "Sudoer":
												if userItems.Value == nil {
													dataWarnings++
													mhUserSudo = false
												} else {
													mhUserSudo = userItems.Value.(bool)
												}
											}
											baselineStruct.musthave.users.users[user] = musthaveusersstruct{
												groups: mhUserGroup,
												shell:  mhUserShell,
												home:   mhUserHome,
												sudoer: mhUserSudo,
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
									for _, baselineItem = range nextBaselineValues {
										switch baselineItem.Key {
										case "Status":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.musthave.policies.polstatus = ""
											} else {
												baselineStruct.musthave.policies.polstatus = baselineItem.Value.(string)
											}
										case "Import":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.musthave.policies.polimport = ""
											} else {
												baselineStruct.musthave.policies.polimport = baselineItem.Value.(string)
											}
										case "Reboot":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.musthave.policies.polreboot = false
											} else {
												baselineStruct.musthave.policies.polreboot = baselineItem.Value.(bool)
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
									for _, baselineItem = range nextBaselineValues {
										rulesValues, ok := baselineItem.Value.(yaml.MapSlice)
										if !ok {
											warnings++
										}
										switch baselineItem.Key {
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
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.musthave.rules.fwzones = []string{""}
											} else {
												mhRulesZones := make([]string, len(baselineItem.Value.([]interface{})))
												mhRulesZonesSlice := baselineItem.Value.([]interface{})
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
									for _, baselineItem = range nextBaselineValues {
										mhMounts = baselineItem.Key.(string)
										mountValues, ok := baselineItem.Value.(yaml.MapSlice)
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
									for _, baselineItem = range nextBaselineValues {
										rulesValues, ok := baselineItem.Value.(yaml.MapSlice)
										if !ok {
											warnings++
										}
										switch baselineItem.Key {
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
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.mustnothave.rules.fwzones = []string{""}
											} else {
												mhRulesZones := make([]string, len(baselineItem.Value.([]interface{})))
												mhRulesZonesSlice := baselineItem.Value.([]interface{})
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
									for _, baselineItem = range nextBaselineValues {
										switch baselineItem.Key {
										case "Logs":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.final.collect.logs = []string{""}
											} else {
												fnlColLog := make([]string, len(baselineItem.Value.([]interface{})))
												fnlColLogSlice := baselineItem.Value.([]interface{})
												for i, v := range fnlColLogSlice {
													fnlColLog[i] = v.(string)
												}
												baselineStruct.final.collect.logs = fnlColLog
											}
										case "Stats":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.final.collect.stats = []string{""}
											} else {
												finalCollectStats := make([]string, len(baselineItem.Value.([]interface{})))
												finalCollectStatsSlice := baselineItem.Value.([]interface{})
												for i, v := range finalCollectStatsSlice {
													finalCollectStats[i] = v.(string)
												}
												baselineStruct.final.collect.stats = finalCollectStats
											}
										case "Files":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.final.collect.files = []string{""}
											} else {
												fnlColFiles := make([]string, len(baselineItem.Value.([]interface{})))
												fnlColFilesSlice := baselineItem.Value.([]interface{})
												for i, v := range fnlColFilesSlice {
													fnlColFiles[i] = v.(string)
												}
												baselineStruct.final.collect.files = fnlColFiles
											}
										case "Users":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.final.collect.users = false
											} else {
												baselineStruct.final.collect.users = baselineItem.Value.(bool)
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
									for _, baselineItem = range nextBaselineValues {
										switch baselineItem.Key {
										case "Services":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.final.restart.services = false
											} else {
												baselineStruct.final.restart.services = baselineItem.Value.(bool)
											}
										case "Servers":
											if baselineItem.Value == nil {
												dataWarnings++
												baselineStruct.final.restart.servers = false
											} else {
												baselineStruct.final.restart.servers = baselineItem.Value.(bool)
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
					commandChannel := make(chan map[string]string)
					var wg sync.WaitGroup
					var commandSync sync.WaitGroup
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

						for i := range sshList {
							if i == pp.fqdn.(string) {
								wg.Add(1)
								disconnectSessions <- false // To keep sessions alive and not disconnect
								go pp.connectAndRunBaseline(commandChannel,
									servername.(string),
									output,
									disconnectSessions,
									readyState,
									&wg,
									&commandSync)
							}
							go func() {
								wg.Wait()
								close(output)
							}()
						}
						fmt.Printf("%d  %d\n", len(readyState), len(sshList)) // delete later, debugging
						for {
							if len(readyState) == len(sshList) {
								break
							}
						}
						baselineStruct.applyPrereq(&sshList, commandChannel)
						baselineStruct.applyMustHaves(&sshList, &rebootBool, commandChannel)
						baselineStruct.applyMustNotHaves(&sshList, commandChannel)
						baselineStruct.applyFinals(&sshList, &rebootBool, commandChannel)
						channelreaderlib.ChannelReaderGroups(output, &wg) // move to each apply stage
					}
				} else if strings.ToLower(serverGroupName) == "all" {
					var allServers yaml.MapSlice
					output := make(chan string)
					readyState := make(chan []bool)
					commandChannel := make(chan map[string]string)
					var wg sync.WaitGroup
					var commandSync sync.WaitGroup
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
						for i := range sshList {
							if i == pp.fqdn.(string) {
								wg.Add(1)
								disconnectSessions <- false // To keep sessions alive and not disconnect
								go pp.connectAndRunBaseline(commandChannel,
									servername.(string),
									output,
									disconnectSessions,
									readyState,
									&wg,
									&commandSync)
							}
						}
						for {
							if len(readyState) == len(sshList) {
								break // wait for all servers to be connected and ready to accept commands
							}
						}
						go func() {
							wg.Wait()
							close(output)
						}()
						baselineStruct.applyPrereq(&sshList, commandChannel)
						baselineStruct.applyMustHaves(&sshList, &rebootBool, commandChannel)
						baselineStruct.applyMustNotHaves(&sshList, commandChannel)
						baselineStruct.applyFinals(&sshList, &rebootBool, commandChannel)
						channelreaderlib.ChannelReaderBaselines(output, &wg, &commandSync) // move to each apply stage
					}
				}
			}
			if rebootBool {
				// send reboot command to servers
			}
			disconnectSessions <- true // When read, channels/sessions will be closed
		}
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// connectAndRunBaseline Establish a connection and run command(s), will add CLI args in the near future
func (parsedData *ParsedPool) connectAndRunBaseline(command chan map[string]string,
	servername string,
	output chan<- string,
	disconnect <-chan bool,
	ready chan []bool,
	wg *sync.WaitGroup,
	commandSync *sync.WaitGroup) {
	pp := parsedData
	authMethodCheck := []ssh.AuthMethod{}
	key, err := ioutil.ReadFile(pp.keypath.(string))
	if err != nil {
		authMethodCheck = append(authMethodCheck, ssh.Password(pp.password.(string)))
	} else {
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			loggerlib.GeneralError(servername, "[INFO: Failed To Parse Private Key] ", err)
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
			disconnectCheck := <-disconnect
			if !disconnectCheck {
				if len(commandCheck) != 0 {
					for i, j := range commandCheck {
						if i == pp.fqdn && j != "" {
							commandSync.Add(1)
							output <- executeBaselines(servername, j, pp.password.(string), connection)
							commandCheck[i] = ""
						}
					}
					command <- commandCheck
					commandSync.Done()
				}
			} else {
				break
			}
			// close sync group here?
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
	var terminalOutput []byte
	var waitOutput sync.WaitGroup
	// it does not wait for output on some machines that are taking too long to respond. I'd like to avoid using Rlocks/Runlocks for this
	waitOutput.Add(1)
	go func(in io.WriteCloser, out io.Reader, terminalOutput *[]byte) {
		var (
			line string
			read = bufio.NewReader(out)
		)
		for {
			buffer, err := read.ReadByte()
			if err != nil {
				break
			}
			*terminalOutput = append(*terminalOutput, buffer)
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
		waitOutput.Done()
	}(in, out, &terminalOutput)
	_, err = session.Output(cmd)
	waitOutput.Wait()
	if err != nil {
		validator = "NOK\n"
		loggerlib.ErrorLogger(servername, "[INFO: Failed] ", terminalOutput)
	} else {
		validator = "OK\n"
		loggerlib.OutputLogger(servername, "[INFO: Success] ", terminalOutput)
	}
	return validator
}
