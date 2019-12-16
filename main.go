package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/APoniatowski/GoSSH/sshlib"
	"github.com/APoniatowski/GoSSH/yamlparser"

	
	"gopkg.in/yaml.v2"
)

// Error checking function
func generalError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var config yaml.MapSlice
var waittotal int
var grouptotal int
var serversPerGroup []int

func init() {
	fmt.Println("Parsing data...")
	data := yamlparser.ParseServersList()
	err := yaml.Unmarshal([]byte(data), &config)
	generalError(err)
	fmt.Println("Data parsed, no errors encountered...")
	waittotal = yamlparser.TotalServercount(config)
	serversPerGroup = yamlparser.ServersPerGroupcount(config)
	grouptotal = len(config)
}

// Main function to carry out operations
func main() {
	//TODO 4 functions will be made of this in ssh-lib, each using anonymous functions and goroutines, which will be chosen via cmdline args
	//TODO first option: run groups sequentially and servers concurrently
	//TODO second option: run groups concurrently and servers sequentially
	//TODO third option: run groups and servers sequentially
	//TODO fourth option: apocalypse mode, run groups and servers concurrently, execute all in other words

	fmt.Printf("Total number of servers: %d\n", waittotal)
	fmt.Printf("Total number of servers per group: ")
	for _, totalItem := range serversPerGroup {
		fmt.Printf(" %d ", totalItem)
	}
	fmt.Printf("\n")
	fmt.Printf("Total groups of servers: %d\n", grouptotal)
	for groupIndex, groupItem := range config {
		fmt.Printf("%d %s %T:\n", groupIndex, groupItem.Key, groupItem.Key)
		groupValue, ok := groupItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
		}
		for serverIndex, serverItem := range groupValue {
			fmt.Printf("\t%d %s %T:\n", serverIndex, serverItem.Key, serverItem.Key)
			serverValue, ok := serverItem.Value.(yaml.MapSlice)
			if !ok {
				panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
			}

			fqdn := serverValue[0].Value
			username := serverValue[1].Value
			password := serverValue[2].Value
			keypath := serverValue[3].Value
			port := serverValue[4].Value

			if username == nil {
				username = "root"
				fmt.Println("No username specified in config.yml, defaulting to 'root'...")
			}
			if password == nil {
				password = ""
				fmt.Println("No password specified in config.yml, defaulting to SSH key based authentication...")
			}
			if keypath == nil {
				keypath = ""
				fmt.Println("No username specified in config.yml, defaulting to password based authentication...")
			}
			if port == nil {
				port = 22
				port = strconv.Itoa(port.(int))
				fmt.Println("No port specified in config.yml, defaulting to port 22...")
			} else {
				port = strconv.Itoa(serverValue[4].Value.(int))
			}

			sshlib.ConnectAndRun(fqdn.(string), username.(string), password.(string), keypath.(string), port.(string))
		}
	}
	//TODO /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}
