package sshlib

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// ApplyBaselines Apply defined baselines. Returns an aggregated ComplianceReport
// across every baseline group so the caller can derive a real process exit code.
func ApplyBaselines(baselineYAML *yaml.MapSlice, configs *yaml.MapSlice) *ComplianceReport {
	report := &ComplianceReport{}
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
			baselineStepsValue, ok := groupItem.Value.(yaml.MapSlice)
			if !ok {
				panic("\nError parsing server groups.\nAborting...\n")
			}
			if strings.ToLower(serverGroupName) == "all" {
				fmt.Println("Applying baseline on all servers:")
			} else {
				fmt.Println("Applying baseline on", serverGroupName+":")
			}
			baselineStruct, _ := parseBaselineGroup(baselineStepsValue)
			groupReport := baselineStruct.runBaseline(serverGroupName, configs, modeApply)
			report.Results = append(report.Results, groupReport.Results...)
		}
	}
	fmt.Println(report.Summary())
	return report
}
