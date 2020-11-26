package sshlib

import (
	"io/ioutil"
	"time"

	"github.com/APoniatowski/GoSSH/loggerlib"
	"golang.org/x/crypto/ssh"
)

func (parseddata *ParsedPool) connectAndRunSeq(command *string, servername string) string {
	pp := parseddata
	derefcmd := *command
	authMethodCheck := []ssh.AuthMethod{}
	key, err := ioutil.ReadFile(pp.keypath.(string))
	if err != nil {
		authMethodCheck = append(authMethodCheck, ssh.Password(pp.password.(string)))
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
		User:            pp.username.(string),
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
			recoveries = recv
		}
	}()
	connection, err := ssh.Dial("tcp", pp.fqdn.(string)+":"+pp.port.(string), sshConfig)
	if err != nil {
		loggerlib.GeneralError(servername, "[ERROR: Connection Failed] ", err)
		validator = "NOK\n"
		return validator
	}
	defer connection.Close()
	derefcmd = OSSwitcher.Switcher(*pp, derefcmd)
	// fmt.Printf("%v: ", servername)
	return executeCommand(servername, derefcmd, pp.password.(string), connection)
}
