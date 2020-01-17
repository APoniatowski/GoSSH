package sshlib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"

	"github.com/APoniatowski/GoSSH/channelreaderlib"
	"github.com/APoniatowski/GoSSH/loggerlib"
	knownhosts "golang.org/x/crypto/ssh/knownhosts"
)

// executeCommand function to run a command on remote servers. Arguments will run through this function and will take strings,
func executeCommand(servername string, cmd string, connection *ssh.Client) string {
	session, err := connection.NewSession()
	loggerlib.GeneralError(err)
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 50, 100, modes); err != nil {
		session.Close()
		log.Fatal(err)
	}

	var validator string
	// shellErr := session.Shell()
	// if shellErr != nil {
	// 	log.Fatal(shellErr)
	// }
	terminaloutput, err := session.CombinedOutput(cmd)
	if err != nil {
		validator = "Failed\n"
		loggerlib.ErrorLogger(servername, terminaloutput)
	} else {
		validator = "Ok\n"
		loggerlib.OutputLogger(servername, terminaloutput)
	}

	return validator
}

// connectAndRun Establish a connection and run command(s), will add CLI args in the near future
func connectAndRun(command *string, servername string, fqdn string, username string, password string, keypath string, port string, output chan<- string, wg *sync.WaitGroup) {
	authMethodCheck := []ssh.AuthMethod{}
	key, err := ioutil.ReadFile(keypath)
	if err != nil {
		authMethodCheck = append(authMethodCheck, ssh.Password(password))
	} else {
		signer, err := ssh.ParsePrivateKey(key)
		loggerlib.GeneralError(err)
		authMethodCheck = append(authMethodCheck, ssh.PublicKeys(signer))
	}
	filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	hostKeyCallback, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}
	sshConfig := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethodCheck,
		HostKeyCallback: hostKeyCallback,
		Timeout:         5 * time.Second,
	}
	connection, err := ssh.Dial("tcp", fqdn+":"+port, sshConfig)
	loggerlib.GeneralError(err)
	defer connection.Close()
	defer wg.Done()
	output <- executeCommand(servername, *command, connection)
}

func connectAndRunSeq(command *string, servername string, fqdn string, username string, password string, keypath string, port string) string {
	authMethodCheck := []ssh.AuthMethod{}
	key, err := ioutil.ReadFile(keypath)
	if err != nil {
		authMethodCheck = append(authMethodCheck, ssh.Password(password))
	} else {
		signer, err := ssh.ParsePrivateKey(key)
		loggerlib.GeneralError(err)
		authMethodCheck = append(authMethodCheck, ssh.PublicKeys(signer))
	}
	filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	hostKeyCallback, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}
	sshConfig := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethodCheck,
		HostKeyCallback: hostKeyCallback,
		Timeout:         5 * time.Second,
	}
	connection, err := ssh.Dial("tcp", fqdn+":"+port, sshConfig)
	loggerlib.GeneralError(err)
	defer connection.Close()
	return servername + ": " + executeCommand(servername, *command, connection)
}

//=============================== sequential and concurrent functions listed below =============================

// RunSequentially Function for running everything sequentially, this will be the default behaviour
func RunSequentially(configs *yaml.MapSlice, command *string) {
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
			output := connectAndRunSeq(command, servername.(string), fqdn.(string), username.(string), password.(string), keypath.(string), port.(string))
			fmt.Print(output)
		}
	}
}

// RunGroups This will run servers concurrently and groups sequentially
func RunGroups(configs *yaml.MapSlice, command *string) {
	for _, groupItem := range *configs {
		output := make(chan string)
		var wg sync.WaitGroup
		fmt.Printf("Processing %s:\n", groupItem.Key)
		groupValue, ok := groupItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
		}
		// wg.Add(yamlparser.ServersPerGroup[groupIndex])
		for _, serverItem := range groupValue {
			wg.Add(1)
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
			go connectAndRun(command, servername.(string), fqdn.(string), username.(string), password.(string), keypath.(string), port.(string), output, &wg)
		}
		// Lesson learned with go routines... when waiting for waitgroup to decrement inside the loop will wait forever
		// when reading from the channel, defer wg.Done() inside the function run in a goroutine, as it needs to tell the waitgroup
		// to decrement the waitgroup amount, as the channel never closes below, when calling wg.Done() in the loop
		go func() {
			wg.Wait()
			close(output)
		}()
		channelreaderlib.ChannelReaderGroups(output, &wg)
	}

}

// RunAllServers As the function implies, this will run all servers concurrently
func RunAllServers(configs *yaml.MapSlice, command *string) {
	var allServers yaml.MapSlice
	output := make(chan string)
	var wg sync.WaitGroup

	// Concatenates the groups to create a single group
	for _, groupItem := range *configs {
		groupValue, ok := groupItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
		}
		for _, serverItem := range groupValue {
			allServers = append(allServers, serverItem)
		}
	}
	// wg.Add(yamlparser.Waittotal)
	for _, serverItem := range allServers {
		wg.Add(1)
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
		go connectAndRun(command, servername.(string), fqdn.(string), username.(string), password.(string), keypath.(string), port.(string), output, &wg)
	}
	// Lesson learned with go routines... when waiting for waitgroup to decrement inside the loop will wait forever
	// when reading from the channel, defer wg.Done() inside the function run in a goroutine, as it needs to tell the waitgroup
	// to decrement the waitgroup amount, as the channel never closes below, when calling wg.Done() in the loop

	// this resolved the stuck go routines. I needed to close the channel, as the channelreader gets stuck at an open and empty channel
	go func() {
		wg.Wait()
		close(output)
	}()
	channelreaderlib.ChannelReaderAll(output, &wg)
}
