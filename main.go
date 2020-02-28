package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/APoniatowski/GoSSH/clioptions"
	"github.com/APoniatowski/GoSSH/examplegenerator"
	"github.com/APoniatowski/GoSSH/sshlib"
	"github.com/APoniatowski/GoSSH/yamlparser"
	"github.com/gookit/color"
	"github.com/urfave/cli"
)

// Main function to carry out operations
func main() {
	var cmd []string
	var switches *sshlib.Switches = &sshlib.OSSwitcher
	cyan := color.Cyan.Render
	// had issues with pointers and booleans, so I found a article which gave this as a solution:
	toggleswitchtrue := true
	toggleswitchfalse := false
	switches.Updater = &toggleswitchfalse
	switches.UpdaterFull = &toggleswitchfalse
	switches.Install = &toggleswitchfalse
	switches.Uninstall = &toggleswitchfalse
	// create one variable and reference it, so essentially this would take less expensive, as it is pointing to 2 memory locations
	// instead of 4.
	app := cli.NewApp()
	app.Name = "GoSSH"
	app.Version = "1.5.0"
	app.Usage = "Open Source Go Infrastucture Automation Tool"
	app.UsageText = "GoSSH [global options] command [subcommand] [script or arguments...]"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:    "sequential",
			Aliases: []string{"s"},
			Usage:   "Run the command sequentially on all servers in your pool",
			Action: func(c *cli.Context) error {
				yamlparser.Rollcall()
				cmd = os.Args[2:]
				command := clioptions.GeneralCommandParse(cmd)
				sshlib.RunSequentially(&yamlparser.Config, &command)
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "run",
					Usage: "Run a bash script on the servers in your pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						cmd := os.Args[3]
						cmdargs := os.Args[4:]
						command := clioptions.BashScriptParse(cmd, cmdargs)
						sshlib.RunSequentially(&yamlparser.Config, &command)
						return nil
					},
				},
				{
					Name:  "update",
					Usage: "update all remote servers in pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						command := ""
						osSwitch := strings.Join(os.Args[3:], " ")
						if osSwitch == "os" || osSwitch == "OS" {
							switches.UpdaterFull = &toggleswitchtrue
						} else {
							switches.Updater = &toggleswitchtrue
						}
						sshlib.RunSequentially(&yamlparser.Config, &command)
						return nil
					},
				},
				{
					Name:  "install",
					Usage: "Install packages on all remote servers in pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						command := ""
						cmdargs := os.Args[3:]
						command = strings.Join(cmdargs, " ")
						switches.Install = &toggleswitchtrue
						sshlib.RunSequentially(&yamlparser.Config, &command)
						return nil
					},
				},
				{
					Name:  "uninstall",
					Usage: "Uninstall packages on all remote servers in pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						command := ""
						cmdargs := os.Args[3:]
						command = strings.Join(cmdargs, " ")
						switches.Uninstall = &toggleswitchtrue
						sshlib.RunSequentially(&yamlparser.Config, &command)
						return nil
					},
				},
				//-----------------placeholder--------------------
				// {
				// 	Name:  "remove",
				// 	Usage: "remove an existing template",
				// 	Action: func(c *cli.Context) error {
				// 		fmt.Println("removed task template: ", c.Args().First())
				// 		return nil
				// 	},
				// },
				//-----------------placeholder--------------------
			},
		},
		{
			Name:    "groups",
			Aliases: []string{"g"},
			Usage:   "Run the command on all servers, per group concurrently in your pool",
			Action: func(c *cli.Context) error {
				yamlparser.Rollcall()
				cmd = os.Args[2:]
				command := clioptions.GeneralCommandParse(cmd)
				if yamlparser.Grouptotal == 1 {
					fmt.Println("There is only 1 group in the pool, consider segmenting them")
					fmt.Println("or alternatively use 'all' instead")
					sshlib.RunAllServers(&yamlparser.Config, &command)
				} else {
					sshlib.RunGroups(&yamlparser.Config, &command)
				}
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "run",
					Usage: "Run a bash script on the servers in your pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						cmd := os.Args[3]
						cmdargs := os.Args[4:]
						command := clioptions.BashScriptParse(cmd, cmdargs)
						if yamlparser.Grouptotal == 1 {
							fmt.Println("There is only 1 group in the pool, consider segmenting them")
							fmt.Println("or alternatively use 'all' instead")
							sshlib.RunAllServers(&yamlparser.Config, &command)
						} else {
							sshlib.RunGroups(&yamlparser.Config, &command)
						}
						return nil
					},
				},
				{
					Name:  "update",
					Usage: "update all remote servers in pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						command := ""
						osSwitch := strings.Join(os.Args[3:], " ")
						if osSwitch == "os" || osSwitch == "OS" {
							switches.UpdaterFull = &toggleswitchtrue
						} else {
							switches.Updater = &toggleswitchtrue
						}
						if yamlparser.Grouptotal == 1 {
							fmt.Println("There is only 1 group in the pool, consider segmenting them")
							fmt.Println("or alternatively use 'all' instead")
							sshlib.RunAllServers(&yamlparser.Config, &command)
						} else {
							sshlib.RunGroups(&yamlparser.Config, &command)
						}
						return nil
					},
				},
				{
					Name:  "install",
					Usage: "Install packages on all remote servers in pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						command := ""
						cmdargs := os.Args[3:]
						command = strings.Join(cmdargs, " ")
						switches.Install = &toggleswitchtrue
						if yamlparser.Grouptotal == 1 {
							fmt.Println("There is only 1 group in the pool, consider segmenting them")
							fmt.Println("or alternatively use 'all' instead")
							sshlib.RunAllServers(&yamlparser.Config, &command)
						} else {
							sshlib.RunGroups(&yamlparser.Config, &command)
						}
						return nil
					},
				},
				{
					Name:  "uninstall",
					Usage: "Uninstall packages on all remote servers in pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						command := ""
						cmdargs := os.Args[3:]
						command = strings.Join(cmdargs, " ")
						switches.Uninstall = &toggleswitchtrue
						if yamlparser.Grouptotal == 1 {
							fmt.Println("There is only 1 group in the pool, consider segmenting them")
							fmt.Println("or alternatively use 'all' instead")
							sshlib.RunAllServers(&yamlparser.Config, &command)
						} else {
							sshlib.RunGroups(&yamlparser.Config, &command)
						}
						return nil
					},
				},
				//-----------------placeholder--------------------
				// {
				// 	Name:  "remove",
				// 	Usage: "remove an existing template",
				// 	Action: func(c *cli.Context) error {
				// 		fmt.Println("removed task template: ", c.Args().First())
				// 		return nil
				// 	},
				// },
				//-----------------placeholder--------------------
			},
		},
		{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "Run the command on all servers concurrently in your pool",
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
					Usage: "Run a bash script on the servers in your pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						cmd := os.Args[3]
						cmdargs := os.Args[4:]
						command := clioptions.BashScriptParse(cmd, cmdargs)
						sshlib.RunAllServers(&yamlparser.Config, &command)
						return nil
					},
				},
				{
					Name:  "update",
					Usage: "Update all remote servers in pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						command := ""
						osSwitch := strings.Join(os.Args[3:], " ")
						if osSwitch == "os" || osSwitch == "OS" {
							switches.UpdaterFull = &toggleswitchtrue
						} else {
							switches.Updater = &toggleswitchtrue
						}
						sshlib.RunAllServers(&yamlparser.Config, &command)
						return nil
					},
				},
				{
					Name:  "install",
					Usage: "Install packages on all remote servers in pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						command := ""
						cmdargs := os.Args[3:]
						command = strings.Join(cmdargs, " ")
						switches.Install = &toggleswitchtrue
						sshlib.RunAllServers(&yamlparser.Config, &command)
						return nil
					},
				},
				{
					Name:  "uninstall",
					Usage: "Uninstall packages on all remote servers in pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						command := ""
						cmdargs := os.Args[3:]
						command = strings.Join(cmdargs, " ")
						switches.Uninstall = &toggleswitchtrue
						sshlib.RunAllServers(&yamlparser.Config, &command)
						return nil
					},
				},
				//-----------------placeholder--------------------
				// {
				// 	Name:  "remove",
				// 	Usage: "remove an existing template",
				// 	Action: func(c *cli.Context) error {
				// 		fmt.Println("removed task template: ", c.Args().First())
				// 		return nil
				// 	},
				// },
				//-----------------placeholder--------------------
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"gen"},
			Usage:   "Generate a template file or print an example, if there are none",
			Action: func(c *cli.Context) error {
				fmt.Println("If you would like to specify 'pool' as your option")
				fmt.Println("more options will be available, as development continues.")
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "pool",
					Usage: "Generate or print an example of a pool.yml file. Use 'template' or 'example' as arguments",
					Action: func(c *cli.Context) error {
						cmd = os.Args[3:]
						command := strings.Join(cmd, " ")
						if command == "template" {
							err := examplegenerator.GeneratePool()
							if err != nil {
								fmt.Println("./config/pool.yml already exists. ")
								fmt.Println("If you would like to recreate this file,")
								fmt.Println("delete the existing pool.yml file and try again")
							} else {
								fmt.Println("Template has now been created")
							}
						} else if command == "example" {
							examplegenerator.PrintPoolExample()
						} else {
							fmt.Println("Please choose ", cyan("template"), " or ", cyan("example"))
							fmt.Println("to generate, or display a template")
						}
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
