package main

import (
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
	app.Version = "1.3.0"
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
						cmdargs := os.Args[4]
						if cmdargs == "os" || cmdargs == "OS" {
							yamlparser.UpdaterFull = true
						}
						yamlparser.Updater = true
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
						yamlparser.Install = true
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
			Usage:   "Run the command on all servers per group concurrently in your pool",
			Action: func(c *cli.Context) error {
				yamlparser.Rollcall()
				cmd = os.Args[2:]
				command := clioptions.GeneralCommandParse(cmd)
				sshlib.RunGroups(&yamlparser.Config, &command)
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
						sshlib.RunGroups(&yamlparser.Config, &command)
						return nil
					},
				},
				{
					Name:  "update",
					Usage: "update all remote servers in pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						command := ""
						cmdargs := os.Args[4]
						if cmdargs == "os" || cmdargs == "OS" {
							yamlparser.UpdaterFull = true
						}
						yamlparser.Updater = true
						sshlib.RunGroups(&yamlparser.Config, &command)
						return nil
					},
				},
				{
					Name:  "install",
					Usage: "Install packages on all remote servers in pool",
					Action: func(c *cli.Context) error {
						yamlparser.Rollcall()
						command := ""
						yamlparser.Install = true
						sshlib.RunGroups(&yamlparser.Config, &command)
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
					Usage: "Run a bash script on the defined servers",
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
						cmdargs := os.Args[4]
						if cmdargs == "os" || cmdargs == "OS" {
							yamlparser.UpdaterFull = true
						}
						yamlparser.Updater = true
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
						yamlparser.Install = true
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
