package clioptions

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// GeneralCommandParse CLI parser for general commands to linux servers
func GeneralCommandParse(cmd []string) string {
	command := strconv.Quote(strings.Join(cmd, " "))
	command = "sh -c " + command + " 2>&1"
	return command
}

//BashScriptParse not yet worked on
func BashScriptParse(cmd string, cmdargs []string) string {
	scriptargs := strings.Join(cmdargs, " ")
	script, _ := os.Open(cmd)
	defer script.Close()
	scanner := bufio.NewScanner(script)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		// need to do some conditional logic to find $ and add a \ before it, to make echo work
		lines = append(lines, scanner.Text())
		lines = append(lines, "\n")
	}
	parsedcmd := strconv.Quote(strings.Join(lines, ""))
	parsedcmd = strings.Replace(parsedcmd, `$`, `\$`, -1)
	parsedlines := "set +H;echo -e " + parsedcmd + " > /tmp/gossh-script.sh;bash /tmp/gossh-script.sh " + scriptargs + ";rm /tmp/gossh-script.sh;set -H"
	return parsedlines
}
