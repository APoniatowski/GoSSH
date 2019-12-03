package main

import (
	"fmt"
	"log"
	"io/ioutil"

	"github.com/APoniatowski/GoSSH/yamlparser"
	"github.com/APoniatowski/GoSSH/sshlib"

	"golang.org/x/crypto/ssh"
	// "golang.org/x/crypto/ssh/knownhosts"
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
	data := yamlparser.ParseServersList()
	_ = yaml.Unmarshal([]byte(data), &config)
	// generalError(err)
	//TODO 4 functions will be made of this in ssh-lib, each using anonymous functions and goroutines, which will be chosen via cmdline args
	//TODO first option: run groups sequentially and servers concurrently
	//TODO second option: run groups concurrently and servers sequentially
	//TODO third option: run groups and servers sequentially
	//TODO fourth option: apocalypse mode, run groups and servers concurrently, execute all in other words
	for groupKey, groupValue := range config {
		fmt.Printf("ServerGroup name: %v\n", groupKey)
		for serverKey, serverValue := range groupValue {
			fmt.Printf("\tServer name: %v\n", serverKey)
			//! fmt.Printf("\t\t%v\n", serverValue.FQDN)
			//! fmt.Printf("\t\t%v\n", serverValue.Username)
			//! fmt.Printf("\t\t%v\n", serverValue.Password)
			//! fmt.Printf("\t\t%v\n", serverValue.KeyPath)
			//! fmt.Printf("\t\t%v\n", serverValue.Port)
			//! if statement will be needed to make empty ports default to port 22 and passwords to default to ssh keys
			// FqdnplusPort := serverValue.FQDN + ":" + serverValue.Port

			key, err := ioutil.ReadFile(serverValue.KeyPath)
			generalError(err)
			signer, err := ssh.ParsePrivateKey(key)
			generalError(err)

			sshConfig := &ssh.ClientConfig{
				User: serverValue.Username,
				Auth: []ssh.AuthMethod{
					ssh.PublicKeys(signer),
				},
				// HostKeyCallback: ssh.FixedHostKey(),  //* will need to figure out how to use this for public use...
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}

			connection, err := ssh.Dial("tcp", serverValue.FQDN + ":" + serverValue.Port, sshConfig)
			generalError(err)
			defer connection.Close()
			sshlib.ExecuteCommand("/usr/bin/uptime", connection)
		}
	}
	//TODO /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//! fmt.Printf("Result: %v\n", config)
	//! fmt.Printf("Server22 is: %s\n", config["ServerGroup2"]["Server22"])

}
