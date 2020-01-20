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

//BashScriptParse bash script parser, to pass write the script on the server, run it and remove it.
// It also accepts args for the script. Dependency scripts will not work, as they are considered a separate script
func BashScriptParse(cmd string, cmdargs []string) string {
	scriptargs := strings.Join(cmdargs, " ")
	script, _ := os.Open(cmd)
	defer script.Close()
	scanner := bufio.NewScanner(script)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		lines = append(lines, "\n")
	}
	parsedcmd := strconv.Quote(strings.Join(lines, ""))
	parsedcmd = strings.Replace(parsedcmd, `$`, `\$`, -1)
	parsedlines := "set +H;echo -e " + parsedcmd + " > /tmp/gossh-script.sh;bash /tmp/gossh-script.sh " + scriptargs + ";rm /tmp/gossh-script.sh;set -H"
	return parsedlines
}
