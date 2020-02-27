package yamlparser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/color"
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


// Rollcall Parse and tally up the number of servers and groups and servers in groups
func Rollcall() {
	cyan := color.Cyan.Render
	green := color.Green.Render
	red := color.Red.Render
	yellow := color.Yellow.Render
	fmt.Println(yellow("Parsing data..."))
	data := ParseServersList()
	err := yaml.Unmarshal([]byte(data), &Config)
	if err != nil {
		fmt.Println(yellow("Data could not be parsed, "), red("encountered some errors..."))
		fmt.Println(yellow("Please check your pool.yml file, or generate one with:"))
		fmt.Println(cyan("gossh generate pool template"))
		fmt.Println(yellow("or display a example with:"))
		fmt.Println(cyan("gossh generate pool example"))
	} else {
		fmt.Println(yellow("Data parsed, "), green("no errors encountered..."))
	}
	Waittotal = TotalServercount(Config)
	ServersPerGroup = ServersPerGroupcount(Config)
	Grouptotal = len(Config)
	fmt.Printf(yellow("Total number of servers: %s\n"), Waittotal)
	fmt.Printf(yellow("Total number of servers per group: "))
	for _, totalItem := range ServersPerGroup {
		fmt.Printf(" %s ", totalItem)
	}
	fmt.Printf("\n")
	fmt.Printf(yellow("Total groups of servers: %s\n"), Grouptotal)
	fmt.Printf(yellow("Total number of logical cores: %s\n"), runtime.NumCPU())
	fmt.Printf(yellow("======================================================\n"))
}

// ParseServersList server list parser, parses it to a map of structs in main function
func ParseServersList() string {
	cyan := color.Cyan.Render
	// green := color.Green.Render
	// red := color.Red.Render
	yellow := color.Yellow.Render
	yamlLocation, _ := filepath.Abs("./config/pool-testing12.yml") // remove -testing to revert changes, using custom server list
	bufRead, err := os.Open(yamlLocation)
	if err != nil {
		fmt.Println(yellow("Please check your pool.yml file, or generate one with:"))
		fmt.Println(cyan("gossh generate pool template"))
		fmt.Println(yellow("or display a example with:"))
		fmt.Println(cyan("gossh generate pool example"))
		os.Exit(-1)
	}
	defer bufRead.Close()

	scanner := bufio.NewScanner(bufRead)
	var configYaml []string

	for scanner.Scan() {
		configYaml = append(configYaml, scanner.Text())
	}
	parse := strings.Join(configYaml, "\n")

	return parse
}
