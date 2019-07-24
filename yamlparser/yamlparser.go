package yamlparser

import (
	"io/ioutil"
	"log"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Error checking function
func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

// Config   Yaml structs goes here
type Config struct {
	ServerList map[string]Servers `yaml:"ServerList"`
}

// Servers   Follow up on structs
type Servers struct {
	FQDN     string `yaml:"FQDN"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Key      string `yaml:"Key"`
}

// ParseYAML   function to return the parsed yaml file as an map/dictionary/vector
func ParseYAML() interface{} {
	yamlLocation, _ := filepath.Abs("./config/config.yml")
	configYaml, err := ioutil.ReadFile(yamlLocation)
	check(err)
	var configs Config
	err = yaml.Unmarshal(configYaml, &configs)
	return configs
}
