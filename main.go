package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/APoniatowski/GoSSH/sshlib"
	"github.com/APoniatowski/GoSSH/yamlparser"
)

// Error checking function
func generalError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// Main function to carry out operations
func main() {
	var cmd []string
	if len(os.Args) > 2 {
		cmd = os.Args[2:] // will change this to 3 later, when I see a need to expand on more arguments, eg. running only 1 group, or x amount of servers
	} else {
		fmt.Println("No command was specified, please specify a command.")
		os.Exit(1)
	}

	command := strconv.Quote(strings.Join(cmd, " "))
	command = "sh -c " + command + " 2>&1"
	yamlparser.Rollcall()

	switch options := os.Args[1]; options {
	case "seq":
		sshlib.RunSequentially(&yamlparser.Config, &command)
	case "groups":
		sshlib.RunGroups(&yamlparser.Config, &command)
	case "all":
		sshlib.RunAllServers(&yamlparser.Config, &command)
	default:
		fmt.Println("Usage: gossh [option] [command]")
		fmt.Println("Options:")
		fmt.Println("  seq		- Run the command sequentially on all servers in your config file")
		fmt.Println("  groups	- Run the command on all servers per group concurrently in your config file")
		fmt.Println("  all		- Run the command on all servers concurrently in your config file")
	}

}
