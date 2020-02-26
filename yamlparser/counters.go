package yamlparser

import "gopkg.in/yaml.v2"

// ServersPerGroupcount Will run on init, to get total number of servers per group
func ServersPerGroupcount(parsed yaml.MapSlice) []int {
	spgslice := make([]int, 0)
	var waitttl int
	for _, groupItem := range parsed {
		groupValue, _ := groupItem.Value.(yaml.MapSlice)
		for _, serverItem := range groupValue {
			_ = serverItem
			waitttl++
		}
		spgslice = append(spgslice, waitttl)
		waitttl = 0
	}
	return spgslice
}

// TotalServercount Will run on init, to get a total number of servers in all groups
func TotalServercount(parsed yaml.MapSlice) int {
	var waitttl int
	for _, groupItem := range parsed {
		groupValue, _ := groupItem.Value.(yaml.MapSlice)
		for _, serverItem := range groupValue {
			_ = serverItem
			waitttl++

		}
	}
	return waitttl
}
