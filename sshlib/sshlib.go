package sshlib

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/superhawk610/bar"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"

	// "golang.org/x/crypto/ssh/knownhosts"
	"github.com/APoniatowski/GoSSH/loggerlib"
	"github.com/APoniatowski/GoSSH/yamlparser"
)

// Error checking function
func generalError(e error) {
	if e != nil {
		// log.Fatal(e)
		log.Println(e)
	}
}

// channelReaderAll Function to read channel until it is closed (all servers only)
func channelReaderAll(channel <-chan string, wg *sync.WaitGroup) {
	successcount := 0
	barp := bar.New(yamlparser.Waittotal)
	for i := 0; i < yamlparser.Waittotal; i++ {
		for message := range channel {
			if message == "Ok\n" {
				barp.Tick()
				successcount++
			} else {
				barp.Tick()
			}
		}
	}
	defer fmt.Printf("%d/%d Succeeded\n", successcount, yamlparser.Waittotal)
	defer barp.Done()
}

// channelReaderGroups Function to read channel until it is closed (groups only)
func channelReaderGroups(channel <-chan string, wg *sync.WaitGroup) {
	loopcountval := len(yamlparser.ServersPerGroup) - 1
	var totalsuccesscount int
	for i := 0; i < loopcountval; i++ {
		successcount := 0
		barp := bar.New(yamlparser.ServersPerGroup[i])
		for im := 0; im < yamlparser.ServersPerGroup[i]; im++ {
			for message := range channel {
				if message == "Ok\n" {
					barp.Tick()
					successcount++
					totalsuccesscount++
				} else {
					barp.Tick()
				}
			}
		}
		barp.Done()
		fmt.Printf("%d/%d Succeeded\n", successcount, yamlparser.ServersPerGroup[i])
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
	if err := session.RequestPty("xterm", 50, 100, modes); err != nil {
		session.Close()
		log.Fatal(err)
	}

	var validator string
	// shellErr := session.Shell()
	// if shellErr != nil {
	// 	log.Fatal(shellErr)
	// }
	currentDate := time.Now()
	dateFormatted := currentDate.Format("2006-01-02")
	terminaloutput, err := session.CombinedOutput(cmd)
	if err != nil {
		validator = "Failed\n"
		// path, _ := filepath.Abs("./logs/errors/")
		// err := os.MkdirAll(path, os.ModePerm)
		// if err == nil || os.IsExist(err) {
		// 	errFile, err := os.OpenFile(path+"/"+dateFormatted+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	defer errFile.Close()
		// 	logger := log.New(errFile, "[INFO: Failed] ", log.LstdFlags)
		// 	logger.Print(servername + ": " + string(terminaloutput))
		// } else {
		// 	log.Println(err)
		// }
		loggerlib.ErrorLogger(servername, terminaloutput)
	} else {
		validator = "Ok\n"
		// path, _ := filepath.Abs("./logs/output/")
		// err := os.MkdirAll(path, os.ModePerm)
		// if err == nil || os.IsExist(err) {
		// 	okFile, err := os.OpenFile(path+"/"+dateFormatted+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	defer okFile.Close()
		// 	logger := log.New(okFile, "[INFO: Succeeded] ", log.LstdFlags)
		// 	logger.Print(servername + ": " + string(terminaloutput))
		// } else {
		// 	log.Println(err)
		// }
		loggerlib.OutputLogger(servername, terminaloutput)
	}

	return validator
}

// connectAndRun Establish a connection and run command(s), will add CLI args in the near future
func connectAndRun(command *string, servername string, fqdn string, username string, password string, keypath string, port string, output chan<- string, wg *sync.WaitGroup) {
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
	// cmd := *command
	defer connection.Close()
	defer wg.Done()
	output <- executeCommand(servername, *command, connection)
}

func connectAndRunSeq(command *string, servername string, fqdn string, username string, password string, keypath string, port string) string {
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
	// cmd := *command
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
		channelReaderGroups(output, &wg)
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
	channelReaderAll(output, &wg)
}
