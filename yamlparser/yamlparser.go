package yamlparser

import (
	"bufio"
	"fmt"
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

// Config global var to be able to access it in all packages
var Config yaml.MapSlice

// Waittotal needed in sshlib
var Waittotal int

// Grouptotal needed in sshlib
var Grouptotal int

// ServersPerGroup needed in sshlib
var ServersPerGroup []int

func init() {
	fmt.Println("Parsing data...")
	data := ParseServersList()
	err := yaml.Unmarshal([]byte(data), &Config)
	generalError(err)
	fmt.Println("Data parsed, no errors encountered...")
}

// Rollcall Tally up the number of servers and groups and servers in groups
func Rollcall() {
	Waittotal = TotalServercount(Config)
	ServersPerGroup = ServersPerGroupcount(Config)
	Grouptotal = len(Config)
	fmt.Printf("Total number of servers: %d\n", Waittotal)
	fmt.Printf("Total number of servers per group: ")
	for _, totalItem := range ServersPerGroup {
		fmt.Printf(" %d ", totalItem)
	}
	fmt.Printf("\n")
	fmt.Printf("Total groups of servers: %d\n", Grouptotal)
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
