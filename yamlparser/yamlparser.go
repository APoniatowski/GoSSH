package yamlparser

import (
	"log"
	"os"
	"bufio"
	"path/filepath"
	"strings"


)

// Error checking function
func generalError(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

// ParseServersList server list parser, parses it to a map of structs in main function
func ParseServersList() string {
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

	return parse
}