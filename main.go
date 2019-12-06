package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

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

// Main function to carry out operations
func main() {
	var config yaml.MapSlice
	data := yamlparser.ParseServersList()
	err := yaml.Unmarshal([]byte(data), &config)
	generalError(err)
	//TODO 4 functions will be made of this in ssh-lib, each using anonymous functions and goroutines, which will be chosen via cmdline args
	//TODO first option: run groups sequentially and servers concurrently
	//TODO second option: run groups concurrently and servers sequentially
	//TODO third option: run groups and servers sequentially
	//TODO fourth option: apocalypse mode, run groups and servers concurrently, execute all in other words
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
			// password := serverValue[2].Value
			keypath := serverValue[3].Value
			port := strconv.Itoa(serverValue[4].Value.(int))
			// fmt.Printf("%v %T\n", port, port)

			key, err := ioutil.ReadFile(keypath.(string))
			generalError(err)
			signer, err := ssh.ParsePrivateKey(key)
			generalError(err)

			sshConfig := &ssh.ClientConfig{
				User: username.(string),
				Auth: []ssh.AuthMethod{
					ssh.PublicKeys(signer),
				},
				// HostKeyCallback: ssh.FixedHostKey(),  //* will need to figure out how to use this for public use...
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}

			connection, err := ssh.Dial("tcp", fqdn.(string)+":"+port, sshConfig)
			generalError(err)
			defer connection.Close()
			sshlib.ExecuteCommand("/usr/bin/uptime", connection)
		}
	}

	//TODO /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//! fmt.Printf("Result: %v\n", config)
	//! fmt.Printf("Server22 is: %s\n", config["ServerGroup2"]["Server22"])

}
