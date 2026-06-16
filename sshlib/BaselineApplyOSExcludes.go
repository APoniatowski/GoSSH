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
				pp := parseServer(serverValue)
				sshList[pp.FQDN] = pp.OS
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
						pp := parseServer(serverValue)
						if len(baselineStruct.exclude.osExcl) > 0 {
							for _, ve := range baselineStruct.exclude.osExcl {
								if strings.EqualFold(pp.OS, ve) {
									osNameCheck = true
								}
							}
						}
						if len(baselineStruct.exclude.serversExcl) > 0 {
							for _, ve := range baselineStruct.exclude.serversExcl {
								if strings.EqualFold(pp.FQDN, ve) {
									serverNameCheck = true
								}
							}
						}
						if !serverNameCheck && !osNameCheck {
							sshList[pp.FQDN] = pp.OS
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
						pp := parseServer(serverValue)
						sshList[pp.FQDN] = pp.OS
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
						pp := parseServer(serverValue)
						if len(baselineStruct.exclude.osExcl) > 0 {
							for _, ve := range baselineStruct.exclude.osExcl {
								if strings.EqualFold(pp.OS, ve) {
									osNameCheck = true
								}
							}
						}
						if len(baselineStruct.exclude.serversExcl) > 0 {
							for _, ve := range baselineStruct.exclude.serversExcl {
								if strings.EqualFold(pp.FQDN, ve) {
									serverNameCheck = true
								}
							}
						}
						if !serverNameCheck && !osNameCheck {
							sshList[pp.FQDN] = pp.OS
						}
					}
				}
			}
		}
	}
	return sshList
}
