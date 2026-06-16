package sshlib

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// CheckBaselines Baseline compliancy check. Provides info of if servers meet baseline compliance as defined in the chosen baseline file.
// Returns an aggregated ComplianceReport across every baseline group so the caller can derive a real process exit code.
func CheckBaselines(baselineyaml *yaml.MapSlice, configs *yaml.MapSlice) *ComplianceReport {
	report := &ComplianceReport{}
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
			groupReport := blstruct.runBaseline(servergroupname, configs, modeCheck)
			report.Results = append(report.Results, groupReport.Results...)
		}
	}
	fmt.Println(report.Summary())
	return report
}
