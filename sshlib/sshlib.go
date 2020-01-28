package sshlib

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"

	"github.com/APoniatowski/GoSSH/channelreaderlib"
	"github.com/APoniatowski/GoSSH/loggerlib"
)

//ParsedData Parsing data to struct to cleanup some code
type ParsedData struct {
	fqdn     interface{}
	username interface{}
	password interface{}
	keypath  interface{}
	port     interface{}
}

//defaulter defaults all empty fields in yaml file and to abort if too many values are missing, eg password and key_path
func defaulter(pd *ParsedData) {
	if pd.password == nil && pd.keypath == nil {
		panic(fmt.Sprintf("Both 'Password' and 'Key_Path' fields are empty... Aborting.\n"))
	}
	if pd.username == nil {
		pd.username = "root"
		// fmt.Println("No username specified in config.yml, defaulting to 'root'...")
	}
	if pd.password == nil {
		pd.password = ""
		// fmt.Println("No password specified in config.yml, defaulting to SSH key based authentication...")
	}
	if pd.keypath == nil {
		pd.keypath = ""
		// fmt.Println("No username specified in config.yml, defaulting to password based authentication...")
	}
	if pd.port == nil {
		pd.port = 22
		pd.port = strconv.Itoa(pd.port.(int))
		// fmt.Println("No port specified in config.yml, defaulting to port 22...")
	} else {
		pd.port = strconv.Itoa(pd.port.(int))
	}
}

// executeCommand function to run a command on remote servers. Arguments will run through this function and will take strings,
func executeCommand(servername string, cmd string, password string, connection *ssh.Client) string {
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
	in, err := session.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	out, err := session.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	var validator string
	var terminaloutput []byte
	var waitoutput sync.WaitGroup
	// it does not wait for output on some machines that are taking too long to respond
	waitoutput.Add(1)
	go func(in io.WriteCloser, out io.Reader, terminaloutput *[]byte) {
		var (
			line string
			read = bufio.NewReader(out)
		)
		for {
			buffer, err := read.ReadByte()
			if err != nil {
				break
			}
			*terminaloutput = append(*terminaloutput, buffer)
			if buffer == byte('\n') {
				line = ""
				continue
			}
			line += string(buffer)
			if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
				_, err = in.Write([]byte(password + "\n"))
				if err != nil {
					break
				}
			}
		}
		waitoutput.Done()
	}(in, out, &terminaloutput)
	_, err = session.Output(cmd)
	waitoutput.Wait()
	if err != nil {
		validator = "NOK\n"
		loggerlib.ErrorLogger(servername, "[INFO: Failed] ", terminaloutput)
	} else {
		validator = "OK\n"
		loggerlib.OutputLogger(servername, "[INFO: Success] ", terminaloutput)
	}
	return validator
}

// connectAndRun Establish a connection and run command(s), will add CLI args in the near future
func connectAndRun(command *string, servername string, parseddata *ParsedData, output chan<- string, wg *sync.WaitGroup) {
	pd := *parseddata
	authMethodCheck := []ssh.AuthMethod{}
	key, err := ioutil.ReadFile(pd.keypath.(string))
	if err != nil {
		authMethodCheck = append(authMethodCheck, ssh.Password(pd.password.(string)))
	} else {
		signer, err := ssh.ParsePrivateKey(key)
		loggerlib.GeneralError(err)
		authMethodCheck = append(authMethodCheck, ssh.PublicKeys(signer))
	}
	// hostKeyCallback, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	// if err != nil {
	hostKeyCallback := ssh.InsecureIgnoreHostKey()
	// }
	sshConfig := &ssh.ClientConfig{
		User:            pd.username.(string),
		Auth:            authMethodCheck,
		HostKeyCallback: hostKeyCallback,
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
		},
		Timeout: 15 * time.Second,
	}
	connection, err := ssh.Dial("tcp", pd.fqdn.(string)+":"+pd.port.(string), sshConfig)
	loggerlib.GeneralError(err)
	defer connection.Close()
	defer wg.Done()
	output <- executeCommand(servername, *command, pd.password.(string), connection)
}

func connectAndRunSeq(command *string, servername string, parseddata *ParsedData) string {
	pd := parseddata
	authMethodCheck := []ssh.AuthMethod{}
	key, err := ioutil.ReadFile(pd.keypath.(string))
	if err != nil {
		authMethodCheck = append(authMethodCheck, ssh.Password(pd.password.(string)))
	} else {
		signer, err := ssh.ParsePrivateKey(key)
		loggerlib.GeneralError(err)
		authMethodCheck = append(authMethodCheck, ssh.PublicKeys(signer))
	}
	// hostKeyCallback, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	// if err != nil {
	hostKeyCallback := ssh.InsecureIgnoreHostKey()
	// }
	sshConfig := &ssh.ClientConfig{
		User:            pd.username.(string),
		Auth:            authMethodCheck,
		HostKeyCallback: hostKeyCallback,
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
		},
		Timeout: 15 * time.Second,
	}
	connection, err := ssh.Dial("tcp", pd.fqdn.(string)+":"+pd.port.(string), sshConfig)
	loggerlib.GeneralError(err)
	defer connection.Close()
	return servername + ": " + executeCommand(servername, *command, pd.password.(string), connection)
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
			var pd ParsedData
			pd.fqdn = serverValue[0].Value
			pd.username = serverValue[1].Value
			pd.password = serverValue[2].Value
			pd.keypath = serverValue[3].Value
			pd.port = serverValue[4].Value
			defaulter(&pd)
			output := connectAndRunSeq(command, servername.(string), &pd)
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
		for _, serverItem := range groupValue {
			wg.Add(1)
			servername := serverItem.Key
			serverValue, ok := serverItem.Value.(yaml.MapSlice)
			if !ok {
				panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
			}
			var pd ParsedData
			pd.fqdn = serverValue[0].Value
			pd.username = serverValue[1].Value
			pd.password = serverValue[2].Value
			pd.keypath = serverValue[3].Value
			pd.port = serverValue[4].Value
			defaulter(&pd)
			go connectAndRun(command, servername.(string), &pd, output, &wg)
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
	for _, serverItem := range allServers {
		wg.Add(1)
		servername := serverItem.Key
		serverValue, ok := serverItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
		}
		var pd ParsedData
		pd.fqdn = serverValue[0].Value
		pd.username = serverValue[1].Value
		pd.password = serverValue[2].Value
		pd.keypath = serverValue[3].Value
		pd.port = serverValue[4].Value
		defaulter(&pd)
		go connectAndRun(command, servername.(string), &pd, output, &wg)
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
