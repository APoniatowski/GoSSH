package main

import (
	"fmt"
	"log"

	"github.com/APoniatowski/GoSSH/yamlparser"

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

			for serverContentIndex, serverContentItem := range serverValue {
				fmt.Printf("\t\t%d %s: %v %T\n", serverContentIndex, serverContentItem.Key, serverContentItem.Value, serverContentItem.Value)

				// key, err := ioutil.ReadFile(serverItem.KeyPath)
				// generalError(err)
				// signer, err := ssh.ParsePrivateKey(key)
				// generalError(err)

				// sshConfig := &ssh.ClientConfig{
				// 	User: serverItem.Username,
				// 	Auth: []ssh.AuthMethod{
				// 		ssh.PublicKeys(signer),
				// 	},
				// 	// HostKeyCallback: ssh.FixedHostKey(),  //* will need to figure out how to use this for public use...
				// 	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				// }

				// connection, err := ssh.Dial("tcp", serverItem.FQDN+":"+serverItem.Port, sshConfig)
				// generalError(err)
				// defer connection.Close()
				// sshlib.ExecuteCommand("/usr/bin/uptime", connection)
			}
		}
	}

	//TODO /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//! fmt.Printf("Result: %v\n", config)
	//! fmt.Printf("Server22 is: %s\n", config["ServerGroup2"]["Server22"])

}
