package sshlib

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// VerifyBaselines Baseline verification. Check if baselines has no errors and provide steps, as to what will be done, if applied
func VerifyBaselines(baselineyaml *yaml.MapSlice) {
	// Baseline
	for _, blItem := range *baselineyaml {
		fmt.Printf("%s:\n", blItem.Key)
		groupValues, ok := blItem.Value.(yaml.MapSlice)
		if !ok {
			panic("\nError:\nCheck your baseline for issues\nAlternatively generate a template to see what is missing/wrong\n")
		}
		// Server groups
		for _, groupItem := range groupValues {
			// initialize the data
			servergroupname := groupItem.Key.(string)
			blstepsValue, ok := groupItem.Value.(yaml.MapSlice)
			if !ok {
				panic("\nError parsing server groups.\nAborting...\n")
			}
			if strings.ToLower(servergroupname) == "all" {
				fmt.Println("Checking baseline on all servers:")
			} else {
				fmt.Println("Checking baseline on", servergroupname+":")
			}
			blstruct, _ := parseBaselineGroup(blstepsValue)
			blstruct.verification(servergroupname)
		}
	}
}
