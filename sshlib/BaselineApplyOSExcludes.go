package sshlib

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

func (baselineStruct *ParsedBaseline) applyOSExcludes(serverGroupName string, configs *yaml.MapSlice) map[string]string {
	sshList := make(map[string]string)
	if strings.ToLower(serverGroupName) == "all" {
		if len(baselineStruct.exclude.osExcl) == 0 &&
			len(baselineStruct.exclude.serversExcl) == 0 {
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
				if groupItem.Key == serverGroupName {
					for _, serverItem := range groupValue {
						var osNameCheck bool
						var serverNameCheck bool
						serverValue, ok := serverItem.Value.(yaml.MapSlice)
						if !ok {
							panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
						}
						if len(baselineStruct.exclude.osExcl) > 0 {
							for _, ve := range baselineStruct.exclude.osExcl {
								if strings.EqualFold(serverValue[5].Value.(string), ve) {
									osNameCheck = true
								}
							}
						}
						if len(baselineStruct.exclude.serversExcl) > 0 {
							for _, ve := range baselineStruct.exclude.serversExcl {
								if strings.EqualFold(serverValue[0].Value.(string), ve) {
									serverNameCheck = true
								}
							}
						}
						if !serverNameCheck && !osNameCheck {
							sshList[serverValue[0].Value.(string)] = serverValue[5].Value.(string)
						}
					}
				}
			}
		}
	} else {
		if len(baselineStruct.exclude.osExcl) == 0 &&
			len(baselineStruct.exclude.serversExcl) == 0 {
			for _, groupItem := range *configs {
				groupValue, ok := groupItem.Value.(yaml.MapSlice)
				if !ok {
					panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
				}
				if strings.EqualFold(groupItem.Key.(string), serverGroupName) {
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
				if strings.EqualFold(groupItem.Key.(string), serverGroupName) {
					for _, serverItem := range groupValue {
						var osNameCheck bool
						var serverNameCheck bool
						serverValue, ok := serverItem.Value.(yaml.MapSlice)
						if !ok {
							panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
						}
						if len(baselineStruct.exclude.osExcl) > 0 {
							for _, ve := range baselineStruct.exclude.osExcl {
								if strings.EqualFold(serverValue[5].Value.(string), ve) {
									osNameCheck = true
								}
							}
						}
						if len(baselineStruct.exclude.serversExcl) > 0 {
							for _, ve := range baselineStruct.exclude.serversExcl {
								if strings.EqualFold(serverValue[0].Value.(string), ve) {
									serverNameCheck = true
								}
							}
						}
						if !serverNameCheck && !osNameCheck {
							sshList[serverValue[0].Value.(string)] = serverValue[5].Value.(string)
						}
					}
				}
			}
		}
	}
	return sshList
}