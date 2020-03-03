package yamlparser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"gopkg.in/yaml.v2"
)

// Rollcall Parse and tally up the number of servers and groups and servers in groups
func ParsePool() {
	cyan := color.Cyan.Render
	green := color.Green.Render
	red := color.Red.Render
	yellow := color.Yellow.Render
	fmt.Println(yellow("Parsing data..."))
	data := ReadPool()
	err := yaml.Unmarshal([]byte(data), &Pool)
	if err != nil {
		fmt.Println(yellow("Data could not be parsed, "), red("encountered some errors..."))
		fmt.Println(yellow("Please check your pool.yml file, or generate one with:"))
		fmt.Println(cyan("gossh generate pool template"))
		fmt.Println(yellow("or display a example with:"))
		fmt.Println(cyan("gossh generate pool example"))
	} else {
		fmt.Println(yellow("Data parsed, "), green("no errors encountered..."))
	}
	Waittotal = TotalServercount(Pool)
	ServersPerGroup = ServersPerGroupcount(Pool)
	Grouptotal = len(Pool)
	fmt.Printf(yellow("Total number of servers: %d\n"), Waittotal)
	fmt.Printf(yellow("Total number of servers per group: "))
	for _, totalItem := range ServersPerGroup {
		fmt.Printf("[%d] ", totalItem)
	}
	fmt.Printf("\n")
	fmt.Printf(yellow("Total groups of servers: %d\n"), Grouptotal)
	fmt.Printf(yellow("Total number of logical cores: %d\n"), runtime.NumCPU())
	fmt.Printf(yellow("======================================================\n"))
}

// ParseServersList server list parser, parses it to a map of structs in main function
func ReadPool() string {
	cyan := color.Cyan.Render
	// green := color.Green.Render
	// red := color.Red.Render
	yellow := color.Yellow.Render
	yamlLocation, _ := filepath.Abs("./config/pool-testing.yml") // remove -testing to revert changes, using custom server list
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
	var poolYaml []string

	for scanner.Scan() {
		poolYaml = append(poolYaml, scanner.Text())
	}
	parse := strings.Join(poolYaml, "\n")

	return parse
}
