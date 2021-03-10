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
			finalCommand.WriteString("cat /proc/loadavg;")
		case "memory":
			finalCommand.WriteString("cat /proc/meminfo;")
		case "storage":
			finalCommand.WriteString("df -T;")
		case "io":
			finalCommand.WriteString("iostat;")
		case "network":
			finalCommand.WriteString("cat /proc/net/dev;")
			finalCommand.WriteString("cat /proc/net/netlink;")

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

