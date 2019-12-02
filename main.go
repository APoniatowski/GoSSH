package main

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
		log.Fatal(e)
	}
}

// Config to structure data for querying and/or running the main function of this tool
type Config struct {
	FQDN     string `yaml:"FQDN"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	KeyPath  string `yaml:"Key_Path"`
	Port     string `yaml:"Port"`
}

// Main function to carry out operations
func main() {
	var config map[string]map[string]Config

	//TODO This part needs to be in a function... struggling with local/global variables does not function like in rust/python /////////////////////
	yamlLocation, _ := filepath.Abs("./config/config.yml")
	bufRead, err := os.Open(yamlLocation)
	generalError(err)
	defer bufRead.Close()

	scanner := bufio.NewScanner(bufRead)
	var configYaml []string

	for scanner.Scan() {
		configYaml = append(configYaml, scanner.Text())
	}
	parse := strings.Join(configYaml, "\n")
	// TODO ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	err = yaml.Unmarshal([]byte(parse), &config)

	//TODO 4 functions will be made of this in ssh-lib, each using anonymous functions and goroutines, which will be chosen via cmdline args
	//TODO first option: run groups sequentially and servers concurrently
	//TODO second option: run groups concurrently and servers sequentially
	//TODO third option: run groups and servers sequentially
	//TODO fourth option: apocalypse mode, run groups and servers concurrently, execute all in other words
	for groupKey, groupValue := range config {
		fmt.Printf("ServerGroup name: %v\n", groupKey)
		for serverKey, serverValue := range groupValue {
			fmt.Printf("\tServer name: %v\n", serverKey)
			fmt.Printf("\t\t%v\n", serverValue.FQDN)
			fmt.Printf("\t\t%v\n", serverValue.Username)
			fmt.Printf("\t\t%v\n", serverValue.Password)
			fmt.Printf("\t\t%v\n", serverValue.KeyPath)
			fmt.Printf("\t\t%v\n", serverValue.Port)
			//! if statement will be needed to make empty ports default to port 22 and passwords to default to ssh keys
			test := serverValue.FQDN + ":" + serverValue.Port
			fmt.Printf("%v %T\n", test, test)
		}
	}
	//TODO /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//! fmt.Printf("Result: %v\n", config)
	//! fmt.Printf("Server22 is: %s\n", config["ServerGroup2"]["Server22"])

}
