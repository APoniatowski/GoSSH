package sshlib

import (
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

// PublicKeyHelper this helper function is to assist in reading and parsing the private key, in order to establish a ssh tunnel
func PublicKeyHelper(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
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
