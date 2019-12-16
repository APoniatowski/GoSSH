package sshlib

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	// "golang.org/x/crypto/ssh/knownhosts"
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
	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := session.StderrPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stderr, sessStderr)
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
