package main

import (
	"fmt"
	"log"
	"os"

	"github.com/APoniatowski/GoSSH/clioptions"
	"github.com/APoniatowski/GoSSH/sshlib"
	"github.com/APoniatowski/GoSSH/yamlparser"
	"github.com/urfave/cli"
)

// Main function to carry out operations
func main() {
	var cmd []string

	app := cli.NewApp()
	app.Name = "GoSSH"
	app.Version = "0.9.0"
	app.Usage = "Open Source Go Infrastucture Automation Tool"
	app.UsageText = "gossh [global options] command [subcommand] [script or arguments...]"
	app.Commands = []cli.Command{
		{
			Name:    "sequential",
			Aliases: []string{"s"},
			Usage:   "Run the command sequentially on all servers in your config file",
			Action: func(c *cli.Context) error {
				yamlparser.Rollcall()
				cmd = os.Args[2:]
				command := clioptions.GeneralCommandParse(cmd)
				sshlib.RunSequentially(&yamlparser.Config, &command)
				return nil
			},
		},
		{
			Name:    "groups",
			Aliases: []string{"g"},
			Usage:   "Run the command on all servers per group concurrently in your config file",
			Action: func(c *cli.Context) error {
				yamlparser.Rollcall()
				cmd = os.Args[2:]
				command := clioptions.GeneralCommandParse(cmd)
				sshlib.RunGroups(&yamlparser.Config, &command)
				return nil
			},
		},
		{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "Run the command on all servers concurrently in your config file",
			Action: func(c *cli.Context) error {
				yamlparser.Rollcall()
				cmd = os.Args[2:]
				command := clioptions.GeneralCommandParse(cmd)
				sshlib.RunAllServers(&yamlparser.Config, &command)
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "run",
					Usage: "Run a bash script on the defined servers",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						yamlparser.Rollcall()
						cmd := os.Args[3]
						cmdargs := os.Args[4:]
						command := clioptions.BashScriptParse(cmd, cmdargs)
						sshlib.RunAllServers(&yamlparser.Config, &command)
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{
	// 		Name:  "lang, l",
	// 		Value: "english",
	// 		Usage: "language for the greeting",
	// 	},
	// }

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}


// switch serveroptions := os.Args[1]; serveroptions {
// case "seq":
// 	command := cli.GeneralCommandParse(cmd)
// 	sshlib.RunSequentially(&yamlparser.Config, &command)
// case "groups":
// 	command := cli.GeneralCommandParse(cmd)
// 	sshlib.RunGroups(&yamlparser.Config, &command)
// case "all":
// 	command := cli.GeneralCommandParse(cmd)
// 	sshlib.RunAllServers(&yamlparser.Config, &command)
// default:
// 	fmt.Println("Usage: gossh [option] [command]")
// 	fmt.Println("Options:")
// 	fmt.Println("  seq		- Run the command sequentially on all servers in your config file")
// 	fmt.Println("  groups	- Run the command on all servers per group concurrently in your config file")
// 	fmt.Println("  all		- Run the command on all servers concurrently in your config file")
// 	fmt.Println("need to add another option [option1 or] [option 2]")
// }

// switch scriptingoptions := os.Args[2]; scriptingoptions {
// case "run":
// 	command := cli.BashScriptParse(cmd)
// 	sshlib.RunAllServers(&yamlparser.Config, &command)
// default:
// 	fmt.Println("Usage: gossh [option] [command]")
// 	fmt.Println("Options:")
// 	fmt.Println("  seq		- Run the command sequentially on all servers in your config file")
// 	fmt.Println("  groups	- Run the command on all servers per group concurrently in your config file")
// 	fmt.Println("  all		- Run the command on all servers concurrently in your config file")
// 	fmt.Println("need to add another option [option1 or] [option 2]")
// }

// if len(os.Args) > 2 {
// 	cmd = os.Args[2:] // will change this to 3 later, when I see a need to expand on more arguments, eg. running only 1 group, or x amount of servers
// } else {
// 	fmt.Println("No command was specified, please specify a command.")
// 	fmt.Println("Usage: gossh [option] [command]")
// 	fmt.Println("Options:")
// 	fmt.Println("  seq		- Run the command sequentially on all servers in your config file")
// 	fmt.Println("  groups	- Run the command on all servers per group concurrently in your config file")
// 	fmt.Println("  all		- Run the command on all servers concurrently in your config file")
// 	fmt.Println("need to add another option [option1 or] [option 2]")
// 	os.Exit(-1)
// }
