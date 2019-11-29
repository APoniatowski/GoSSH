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
	Key_Path string `yaml:"Key_Path"`
	Port     string `yaml:"Port"`
}

// Main function to carry out operations
func main() {
	var config map[string]map[string]Config

	yamlLocation, _ := filepath.Abs("./config/config.yml")
	bufRead, err := os.Open(yamlLocation)
	generalError(err)
	defer bufRead.Close()

	scanner := bufio.NewScanner(bufRead)
	var configYaml []string

	for scanner.Scan() {
		configYaml = append(configYaml, scanner.Text())
	}
	// configYaml, err := ioutil.ReadFile(yamlLocation)
	parsed := strings.Join(configYaml, "\n")

	err = yaml.Unmarshal([]byte(parsed), &config)

	for groupKey, groupValue := range config {
		fmt.Printf("ServerGroup name: %v\n", groupKey)
		for serverKey, serverValue := range groupValue {
			fmt.Printf("\tServer name: %v\n", serverKey)
			fmt.Printf("\t\t%v\n", serverValue.FQDN)
			fmt.Printf("\t\t%v\n", serverValue.Username)
			fmt.Printf("\t\t%v\n", serverValue.Password)
			fmt.Printf("\t\t%v\n", serverValue.Key_Path)
			fmt.Printf("\t\t%v\n", serverValue.Port)
		}
	}

	// fmt.Printf("Result: %v\n", config)
	// fmt.Printf("Server22 is: %s\n", config["ServerGroup2"]["Server22"])

}
