package yamlparser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
var (
	Config          yaml.MapSlice
	Waittotal       int
	Grouptotal      int
	ServersPerGroup []int
)

// func init() {
// }

// Rollcall Tally up the number of servers and groups and servers in groups
func Rollcall() {
	fmt.Println("Parsing data...")
	data := ParseServersList()
	err := yaml.Unmarshal([]byte(data), &Config)
	generalError(err)

	fmt.Println("Data parsed, no errors encountered...")
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
	fmt.Printf("Total number of logical cores: %v\n", runtime.NumCPU())
	fmt.Printf("======================================================\n")
}

// ParseServersList server list parser, parses it to a map of structs in main function
func ParseServersList() string {
	yamlLocation, _ := filepath.Abs("./config/pool-testing1.yml") // remove -testing to revert changes, using custom server list
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

