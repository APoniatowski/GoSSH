package main

import (
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

// Yaml structs goes here
type Config struct {
	ServerList map[string]Servers `yaml:"ServerList"`
}

type Servers struct {
	FQDN     string `yaml:"FQDN"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Key      string `yaml:"Key"`
}

// Main function to carry out operations
func main() {
	yamlLocation, _ := filepath.Abs("./config/config.yml")
	configYaml, err := ioutil.ReadFile(yamlLocation)
	check(err)
	var configs Config
	err = yaml.Unmarshal(configYaml, &configs)
	check(err)
	// will add the funcs from the lib, once I have it setup with the proper args... when I have the time
}
