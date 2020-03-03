package sshlib

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"

	"github.com/APoniatowski/GoSSH/loggerlib"
)

// executeCommand function to run a command on remote servers. Arguments will run through this function and will take strings,
func executeCommand(servername string, cmd string, password string, connection *ssh.Client) string {
	// adding recover to avoid panics during a run. Logs are written, so no need to panic when it its one of
	//the errors below.
	defer func() {
		if recv := recover(); recv != nil {
			recoveries = recv
		}
	}()
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

func executeBaselines(servername string, cmd string, password string, connection *ssh.Client) string {
	// adding recover to avoid panics during a run. Logs are written, so no need to panic when it its one of
	// the errors below.
	defer func() {
		if recv := recover(); recv != nil {
			recoveries = recv
		}
	}()
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
