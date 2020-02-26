package sshlib

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"

	"github.com/APoniatowski/GoSSH/loggerlib"
	"github.com/APoniatowski/GoSSH/pkgmanlib"
)

// Switches For checking what CLI option was used and run the appropriate functions
type Switches struct {
	Updater, UpdaterFull, Install, Uninstall *bool
}

//Switcher Method to check the switches set for each respective action
func (S *Switches) Switcher(pd ParsedData, command string) string {
	rtncommand := ""

	if *S.Updater {
		rtncommand = pkgmanlib.Update(pd.username.(string), pd.os.(string))
	}
	if *S.UpdaterFull {
		rtncommand = pkgmanlib.UpdateOS(pd.username.(string), pd.os.(string))
	}
	if *S.Install {
		rtncommand = pkgmanlib.Install(pd.username.(string), pd.os.(string)) + command + " -y 2>&1"
	}
	if *S.Uninstall {
		rtncommand = pkgmanlib.Uninstall(pd.username.(string), pd.os.(string)) + command + " -y 2>&1"
	}

	return rtncommand
}

// OSSwitcher a much needed var between main and sshlib
var OSSwitcher Switches

//ParsedData Parsing data to struct to cleanup some code
type ParsedData struct {
	fqdn     interface{}
	username interface{}
	password interface{}
	keypath  interface{}
	port     interface{}
	os       interface{}
}

//defaulter defaults all empty fields in yaml file and to abort if too many values are missing, eg password and key_path
func defaulter(pd *ParsedData) {
	if pd.password == nil && pd.keypath == nil {
		panic(fmt.Sprintf("Both 'Password' and 'Key_Path' fields are empty... Aborting.\n"))
	}
	if pd.username == nil {
		pd.username = "root"
	}
	if pd.password == nil {
		pd.password = ""
	}
	if pd.keypath == nil {
		pd.keypath = ""
	}
	if pd.port == nil {
		pd.port = 22
		pd.port = strconv.Itoa(pd.port.(int))
	} else {
		pd.port = strconv.Itoa(pd.port.(int))
	}
}

// executeCommand function to run a command on remote servers. Arguments will run through this function and will take strings,
func executeCommand(servername string, cmd string, password string, connection *ssh.Client) string {
	session, err := connection.NewSession()
	if err != nil {
		loggerlib.GeneralError(servername, "[ERROR: Failed To Create Session] ", err)
	}
	defer session.Close()
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 50, 100, modes); err != nil {
		session.Close()
		loggerlib.GeneralError(servername, "[ERROR: Pty Request Failed] ", err)
	}
	in, err := session.StdinPipe()
	if err != nil {
		loggerlib.GeneralError(servername, "[ERROR: Stdin Error] ", err)
	}
	out, err := session.StdoutPipe()
	if err != nil {
		loggerlib.GeneralError(servername, "[ERROR: Stdout Error] ", err)
	}
	var validator string
	var terminaloutput []byte
	var waitoutput sync.WaitGroup
	// it does not wait for output on some machines that are taking too long to respond. I'd like to avoid using Rlocks/Runlocks for this
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
