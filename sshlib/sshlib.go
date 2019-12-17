package sshlib

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"sync"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	// "golang.org/x/crypto/ssh/knownhosts"

	// "github.com/APoniatowski/GoSSH/yamlparser"
)

// Error checking function
func generalError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// ExecuteCommand function to run a command on remote servers. Arguments will run through this function and will take strings,
// or it will be converted to string if it is any other value. Booleans will most likely not be convertable
func ExecuteCommand(cmd string, connection *ssh.Client) {
	session, err := connection.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	sessStdOut, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}
	// go io.Copy(os.Stdout, sessStdOut)

	sessStderr, err := session.StderrPipe()
	if err != nil {
		panic(err)
	}
	// go io.Copy(os.Stderr, sessStderr)

	if sessStderr == nil {
		fmt.Println(sessStdOut)
	}

	err = session.Run(cmd)
	if err != nil {
		panic(err)
	}
}

// ConnectAndRun Establish a connection and run command(s), will add CLI args in the near future
func ConnectAndRun(fqdn string, username string, password string, keypath string, port string) {
	key, err := ioutil.ReadFile(keypath)
	generalError(err)
	signer, err := ssh.ParsePrivateKey(key)
	generalError(err)
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// HostKeyCallback: ssh.FixedHostKey(),  //* will need to figure out how to use this for public use...
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	connection, err := ssh.Dial("tcp", fqdn+":"+port, sshConfig)
	generalError(err)
	defer connection.Close()
	ExecuteCommand("/usr/bin/uptime", connection) // add CLI arg here
}

//=============================== sequential and concurrent functions listed below =============================

// RunSequentially Function for running everything sequentially, this will be the default behaviour
func RunSequentially(configs *yaml.MapSlice) {
	for groupIndex, groupItem := range *configs {
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
			password := serverValue[2].Value
			keypath := serverValue[3].Value
			port := serverValue[4].Value

			if username == nil {
				username = "root"
				fmt.Println("No username specified in config.yml, defaulting to 'root'...")
			}
			if password == nil {
				password = ""
				fmt.Println("No password specified in config.yml, defaulting to SSH key based authentication...")
			}
			if keypath == nil {
				keypath = ""
				fmt.Println("No username specified in config.yml, defaulting to password based authentication...")
			}
			if port == nil {
				port = 22
				port = strconv.Itoa(port.(int))
				fmt.Println("No port specified in config.yml, defaulting to port 22...")
			} else {
				port = strconv.Itoa(serverValue[4].Value.(int))
			}

			ConnectAndRun(fqdn.(string), username.(string), password.(string), keypath.(string), port.(string))
		}
	}
}

var wg sync.WaitGroup

// RunServersConcurrently As the function implies, this will run servers concurrently and groups sequentially
func RunServersConcurrently(configs *yaml.MapSlice) {
	for groupIndex, groupItem := range *configs {
		fmt.Printf("%d %s %T:\n", groupIndex, groupItem.Key, groupItem.Key)
		groupValue, ok := groupItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
		}
		for serverIndex, serverItem := range groupValue {
			// wg.Add(yamlparser.ServersPerGroup[serverIndex+1])
			go func(serverIndex int, serverItem yaml.MapItem) {
				// defer wg.Done()
				fmt.Printf("\t%d %s %T:\n", serverIndex, serverItem.Key, serverItem.Key)
				serverValue, ok := serverItem.Value.(yaml.MapSlice)
				if !ok {
					panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
				}

				fqdn := serverValue[0].Value
				username := serverValue[1].Value
				password := serverValue[2].Value
				keypath := serverValue[3].Value
				port := serverValue[4].Value

				if username == nil {
					username = "root"
					fmt.Println("No username specified in config.yml, defaulting to 'root'...")
				}
				if password == nil {
					password = ""
					fmt.Println("No password specified in config.yml, defaulting to SSH key based authentication...")
				}
				if keypath == nil {
					keypath = ""
					fmt.Println("No username specified in config.yml, defaulting to password based authentication...")
				}
				if port == nil {
					port = 22
					port = strconv.Itoa(port.(int))
					fmt.Println("No port specified in config.yml, defaulting to port 22...")
				} else {
					port = strconv.Itoa(serverValue[4].Value.(int))
				}

				ConnectAndRun(fqdn.(string), username.(string), password.(string), keypath.(string), port.(string))
			}(serverIndex, serverItem)
			// wg.Wait()
		}
	}
}
