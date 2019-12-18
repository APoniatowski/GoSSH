package sshlib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	// "golang.org/x/crypto/ssh/knownhosts"
	"github.com/APoniatowski/GoSSH/yamlparser"
)

// Error checking function
func generalError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// channelReader Function to read channel until it is closed
func channelReader(channel <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for message := range channel {
		fmt.Print(message)
	}
}

// executeCommand function to run a command on remote servers. Arguments will run through this function and will take strings,
func executeCommand(servername string, cmd string, connection *ssh.Client) string {
	session, err := connection.NewSession()
	if err != nil {
		log.Fatal("Failed to establish a session: ", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		fmt.Errorf("PTY Request Failed: %s", err)
	}

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(cmd)
	terminaloutput := servername + ": " + stdoutBuf.String()
	return terminaloutput
}

// connectAndRun Establish a connection and run command(s), will add CLI args in the near future
func connectAndRun(servername string, fqdn string, username string, password string, keypath string, port string, output chan<- string, wg *sync.WaitGroup) {
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
		Timeout:         5 * time.Second,
	}
	connection, err := ssh.Dial("tcp", fqdn+":"+port, sshConfig)
	generalError(err)
	defer connection.Close()
	defer wg.Done()
	output <- executeCommand(servername, "hostname", connection)
}

func connectAndRunSeq(servername string, fqdn string, username string, password string, keypath string, port string) string {
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
		Timeout:         5 * time.Second,
	}
	connection, err := ssh.Dial("tcp", fqdn+":"+port, sshConfig)
	generalError(err)
	defer connection.Close()
	return executeCommand(servername, "hostname", connection) // add CLI arg here
}

//=============================== sequential and concurrent functions listed below =============================

// RunSequentially Function for running everything sequentially, this will be the default behaviour
func RunSequentially(configs *yaml.MapSlice) {
	for _, groupItem := range *configs {
		fmt.Printf("Processing %s:\n", groupItem.Key)
		groupValue, ok := groupItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
		}
		for _, serverItem := range groupValue {
			servername := serverItem.Key
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
				// fmt.Println("No username specified in config.yml, defaulting to 'root'...")
			}
			if password == nil {
				password = ""
				// fmt.Println("No password specified in config.yml, defaulting to SSH key based authentication...")
			}
			if keypath == nil {
				keypath = ""
				// fmt.Println("No username specified in config.yml, defaulting to password based authentication...")
			}
			if port == nil {
				port = 22
				port = strconv.Itoa(port.(int))
				// fmt.Println("No port specified in config.yml, defaulting to port 22...")
			} else {
				port = strconv.Itoa(serverValue[4].Value.(int))
			}
			if password == nil && keypath == nil {
				panic(fmt.Sprintf("Both 'Password' and 'Key_Path' fields are empty... Aborting.\n"))
			}
			output := connectAndRunSeq(servername.(string), fqdn.(string), username.(string), password.(string), keypath.(string), port.(string))
			fmt.Print(output)
		}
	}
}

// RunServersConcurrently As the function implies, this will run servers concurrently and groups sequentially
func RunServersConcurrently(configs *yaml.MapSlice) {
	for groupIndex, groupItem := range *configs {
		output := make(chan string)
		var wg sync.WaitGroup
		fmt.Printf("Processing %s:\n", groupItem.Key)
		groupValue, ok := groupItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
		}
		wg.Add(yamlparser.ServersPerGroup[groupIndex])
		for _, serverItem := range groupValue {
			// fmt.Printf("%s:\n", serverItem.Key)
			servername := serverItem.Key
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
				// fmt.Println("No username specified in config.yml, defaulting to 'root'...")
			}
			if password == nil {
				password = ""
				// fmt.Println("No password specified in config.yml, defaulting to SSH key based authentication...")
			}
			if keypath == nil {
				keypath = ""
				// fmt.Println("No username specified in config.yml, defaulting to password based authentication...")
			}
			if port == nil {
				port = 22
				port = strconv.Itoa(port.(int))
				// fmt.Println("No port specified in config.yml, defaulting to port 22...")
			} else {
				port = strconv.Itoa(serverValue[4].Value.(int))
			}
			if password == nil && keypath == nil {
				panic(fmt.Sprintf("Both 'Password' and 'Key_Path' fields are empty... Aborting.\n"))
			}
			go connectAndRun(servername.(string), fqdn.(string), username.(string), password.(string), keypath.(string), port.(string), output, &wg)
		}
		// Lesson learned with go routines... when waiting for waitgroup to decrement inside the loop will wait forever
		// when reading from the channel, defer wg.Done() inside the function run in a goroutine, as it needs to tell the waitgroup
		// to decrement the waitgroup amount, as the channel never closes below, when calling wg.Done() in the loop
		go channelReader(output, &wg)
		wg.Wait()

	}

}
