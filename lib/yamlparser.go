package yamlParser

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

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
// parse function to return the parsed yaml file as an map/dictionary/vector
func parseYaml() map[string]string {
	yamlLocation, _ := filepath.Abs("./config/config.yml")
	configYaml, err := ioutil.ReadFile(yamlLocation)
	check(err)
	var configs Config
	err = yaml.Unmarshal(configYaml, &configs)
	check(err)
	parsedYaml := configs
	return configs
}