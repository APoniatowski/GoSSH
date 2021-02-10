package sshlib

import (
	"strings"
)

func finalCommandBuilder(command *string, chosenOption string) string {
	finalCommand := strings.Builder{}
	switch chosenOption {
	case "script":
		finalCommand.WriteString("chmod +x /tmp/" + *command)
		finalCommand.WriteString("/tmp/" + *command)
	case "command":
		finalCommand.WriteString(*command)
	case "logs":
		finalCommand.WriteString("journalctl -u " + *command + " -S today --no-tail > /tmp/" + *command + ".log")
	case "stats":
		switch strings.ToLower(*command) {
		case "cpu":
			finalCommand.WriteString("")
		case "memory":
			finalCommand.WriteString("")
		case "storage":
			finalCommand.WriteString("")
		default:
			finalCommand.WriteString("")
		}
	case "files":
		finalCommand.WriteString("")
	default:
		finalCommand.WriteString("")
	}
	return finalCommand.String()
}

