package yamlparser

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"

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

// public var to avoid the "expression" issues

// ParseYAML   function to return the parsed yaml file as an map/dictionary/vector
func ParseYAML() map[interface{}]interface{} {
	yamlLocation, _ := filepath.Abs("./config/config.yml")
	configYaml, err := ioutil.ReadFile(yamlLocation)
	check(err)
	var configs Config
	err = yaml.Unmarshal([]byte(configYaml), &configs)
	// mapped := map[string]string{}
	// for k, v := range configs.(map[string]interface{}) {
	// 	mapped[k] = v.(string)
	// }
	// fmt.Println(reflect.TypeOf(configs), configs)
	return configs.(map[interface{}]interface{})
}

// InterfaceToMap  to convert an interface to a map to extract the values... the go way of doing it allegedly   {mainly for testing, steep learning curve for this kind of thing}
func InterfaceToMap(i interface{}) (interface{}, error) {
	t := reflect.TypeOf(i)
	switch t.Kind() {
	case reflect.Map:
		v := reflect.ValueOf(i)
		it := reflect.TypeOf((*interface{})(nil)).Elem()
		m := reflect.MakeMap(reflect.MapOf(t.Key(), it))
		for _, mk := range v.MapKeys() {
			m.SetMapIndex(mk, v.MapIndex(mk))
		}
		return m.Interface(), nil
	}
	return nil, errors.New("Unsupported type")
}
