package sshlib

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/APoniatowski/GoSSH/loggerlib"
	"golang.org/x/crypto/ssh"
)

// statusOK / statusNOK are the per-command outcome strings, kept identical to the
// values the legacy command path emits so existing reporting stays compatible.
const (
	statusOK  = "OK\n"
	statusNOK = "NOK\n"
)

// Command is one unit of work fanned out to a single host worker.
type Command struct {
	Host string // fqdn; the worker's identity (matches the sshList key)
	Cmd  string // OS-resolved, sudo-prefixed shell string ("" == no-op)
}

// Result is one outcome returned by a worker for one Command.
type Result struct {
	Host   string
	Status string // statusOK / statusNOK
	Output []byte // accumulated terminal output (for collect steps)
	Err    error
}

// baselineExecutor is the seam between the orchestrator and the transport.
// Production wraps *ssh.Client; tests provide an in-memory fake. This is what
// makes the fan-out/gather logic unit-testable without real SSH.
type baselineExecutor interface {
	// Run executes one command over the persistent connection and reports a
	// status string (statusOK/statusNOK) plus the accumulated terminal output.
	// A transport-level failure is returned as err; a command that ran but failed
	// is reported via statusNOK.
	Run(cmd string) (status string, output []byte, err error)
	Close() error
}

// sshExecutor is the production baselineExecutor. It owns one persistent client.
type sshExecutor struct {
	servername string
	password   string
	conn       *ssh.Client
}

// dialExecutor establishes a persistent SSH connection for one host and returns
// a ready executor. The auth/handshake mirrors the proven path in concurrent.go;
// it returns an error instead of writing to a channel.
func (pp *ParsedPool) dialExecutor(servername string) (baselineExecutor, error) {
	authMethods := []ssh.AuthMethod{}
	key, err := os.ReadFile(pp.KeyPath)
	if err != nil {
		authMethods = append(authMethods, ssh.Password(pp.Password))
	} else {
		signer, perr := ssh.ParsePrivateKey(key)
		if perr != nil {
			loggerlib.GeneralError(servername, "[INFO: Failed To Parse Private Key] ", perr)
			authMethods = append(authMethods, ssh.Password(pp.Password))
		} else {
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}
	}
	sshConfig := &ssh.ClientConfig{
		User:            pp.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
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
	conn, err := ssh.Dial("tcp", pp.FQDN+":"+strconv.Itoa(pp.Port), sshConfig)
	if err != nil {
		return nil, err
	}
	return &sshExecutor{servername: servername, password: pp.Password, conn: conn}, nil
}

// Run executes a single command on a fresh session over the persistent
// connection. It feeds the sudo password when prompted. No shared global state,
// so it is safe to call concurrently across distinct executors.
func (e *sshExecutor) Run(cmd string) (string, []byte, error) {
	session, err := e.conn.NewSession()
	if err != nil {
		loggerlib.GeneralError(e.servername, "[ERROR: Failed To Create Session] ", err)
		return statusNOK, nil, err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 50, 100, modes); err != nil {
		loggerlib.GeneralError(e.servername, "[ERROR: Pty Request Failed] ", err)
		return statusNOK, nil, err
	}

	in, err := session.StdinPipe()
	if err != nil {
		loggerlib.GeneralError(e.servername, "[ERROR: Stdin Error] ", err)
		return statusNOK, nil, err
	}
	out, err := session.StdoutPipe()
	if err != nil {
		loggerlib.GeneralError(e.servername, "[ERROR: Stdout Error] ", err)
		return statusNOK, nil, err
	}

	var terminalOutput []byte
	var waitOutput sync.WaitGroup
	waitOutput.Add(1)
	go func(in io.WriteCloser, out io.Reader, terminalOutput *[]byte) {
		defer waitOutput.Done()
		var line string
		read := bufio.NewReader(out)
		for {
			buffer, rerr := read.ReadByte()
			if rerr != nil {
				break
			}
			*terminalOutput = append(*terminalOutput, buffer)
			if buffer == byte('\n') {
				line = ""
				continue
			}
			line += string(buffer)
			if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
				if _, werr := in.Write([]byte(e.password + "\n")); werr != nil {
					break
				}
			}
		}
	}(in, out, &terminalOutput)

	runErr := session.Run(cmd)
	waitOutput.Wait()
	if runErr != nil {
		loggerlib.ErrorLogger(e.servername, "[INFO: Failed] ", terminalOutput)
		// A non-zero exit means the command RAN and reported failure
		// (non-compliant / unsuccessful) — the connection is still healthy, so
		// report NOK without an error so the fleet keeps the host. Only a
		// transport/session failure (anything that is not an ExitError) is a
		// real error that drops the host from the fleet.
		if _, ok := runErr.(*ssh.ExitError); ok {
			return statusNOK, nil, nil
		}
		return statusNOK, terminalOutput, runErr
	}
	loggerlib.OutputLogger(e.servername, "[INFO: Success] ", terminalOutput)
	return statusOK, terminalOutput, nil
}

func (e *sshExecutor) Close() error {
	return e.conn.Close()
}

// baselineWorker owns exactly one host. It dials once (caller supplies the live
// executor), then serves commands off in until in is closed, emitting one Result
// per command. close(in) is the only shutdown signal; no shared boolean channels.
type baselineWorker struct {
	host    string
	exec    baselineExecutor
	in      chan Command
	results chan<- Result
}

func (w *baselineWorker) run() {
	for cmd := range w.in {
		if cmd.Cmd == "" { // no-op step for this host
			w.results <- Result{Host: w.host, Status: statusOK}
			continue
		}
		status, output, err := w.exec.Run(cmd.Cmd)
		w.results <- Result{Host: w.host, Status: status, Output: output, Err: err}
	}
	w.exec.Close()
}
