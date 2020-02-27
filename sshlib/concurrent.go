package sshlib

import (
	"io/ioutil"
	"sync"
	"time"

	"github.com/APoniatowski/GoSSH/loggerlib"
	"golang.org/x/crypto/ssh"
)

// connectAndRun Establish a connection and run command(s), will add CLI args in the near future
func connectAndRun(command *string, servername string, parseddata *ParsedData, output chan<- string, wg *sync.WaitGroup) {
	pd := parseddata
	derefcmd := *command
	authMethodCheck := []ssh.AuthMethod{}
	key, err := ioutil.ReadFile(pd.keypath.(string))
	if err != nil {
		authMethodCheck = append(authMethodCheck, ssh.Password(pd.password.(string)))
	} else {
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			loggerlib.GeneralError(servername, "[INFO: Failed To Parse PrivKey] ", err)
		}
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
	defer func() {
		if recv := recover(); recv != nil {
		}
	}()
	connection, err := ssh.Dial("tcp", pd.fqdn.(string)+":"+pd.port.(string), sshConfig)
	if err != nil {
		loggerlib.GeneralError(servername, "[ERROR: Connection Failed] ", err)
		validator = "NOK\n"
		output <- validator
		wg.Done()
	} else {
		defer connection.Close()
		defer wg.Done()
		derefcmd = OSSwitcher.Switcher(*pd, derefcmd)
		output <- executeCommand(servername, derefcmd, pd.password.(string), connection)
	}
}
