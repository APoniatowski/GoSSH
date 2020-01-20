package cli-options

import (
	"bufio"
	"fmt"
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
func BashScriptParse(cmd []string) {
	// readfile then readlines, when readline reads { or } then it should add / and when its not, add a ;
	inFile, _ := os.Open(cmd)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	fmt.Println(lines)
	// return lines
}
