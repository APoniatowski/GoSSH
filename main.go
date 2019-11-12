package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Error checking function
func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

type Config struct {
	FQDN     string `yaml:"FQDN"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Key_Path string `yaml:"Key_Path"`
}

// Main function to carry out operations
func main() {
	var config map[string]map[string]Config
	yamlLocation, _ := filepath.Abs("./config/config.yml")
	configYaml, err := ioutil.ReadFile(yamlLocation)
	check(err)

	err = yaml.Unmarshal([]byte(configYaml), &config)

	for k, v := range config {
		fmt.Printf("KEY: %v\n", k)
		for key, val := range v {
			fmt.Printf("KEY2: %v\n", key)
			fmt.Printf("VALUE: %v\n", val)
		}
	}
	// fmt.Printf("Result: %v\n", config)
	// fmt.Printf("Server22 is: %s\n", config["ServerGroup2"]["Server22"])

}
