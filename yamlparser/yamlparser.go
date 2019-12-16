package yamlparser

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Error checking function
func generalError(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

// ParseServersList server list parser, parses it to a map of structs in main function
func ParseServersList() string {
	yamlLocation, _ := filepath.Abs("./config/config.yml") // remove -testing to revert changes, using custom server list
	bufRead, err := os.Open(yamlLocation)
	generalError(err)
	defer bufRead.Close()

	scanner := bufio.NewScanner(bufRead)
	var configYaml []string

	for scanner.Scan() {
		configYaml = append(configYaml, scanner.Text())
	}
	parse := strings.Join(configYaml, "\n")

	return parse
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
