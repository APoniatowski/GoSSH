package sshlib

import (
	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"strings"
)

func (mountDetails *mountdetails) mountCommandBuilder(chosenOption string) string {
	mountCommand := strings.Builder{}
	switch chosenOption {
	case "check":
		//TODO add check here
	case "apply":
		mountCommand.WriteString(pkgmanlib.OmniTools["mkdir"])
		mountCommand.WriteString(mountDetails.dest + " && ")
		mountCommand.WriteString("echo '")
		mountCommand.WriteString(mountDetails.address + ":" + mountDetails.src + " " + mountDetails.dest + " ")
		mountCommand.WriteString(mountDetails.mounttype)
		mountCommand.WriteString(" defaults 0 0") // Default mounting details
		mountCommand.WriteString("' >> /etc/fstab;")
		mountCommand.WriteString(pkgmanlib.OmniTools["mount"] + mountDetails.dest)
		mountCommand.WriteString(" && ")
		mountCommand.WriteString(pkgmanlib.OmniTools["mount"] + "-F ")
		mountCommand.WriteString(mountDetails.mounttype)
		if mountDetails.username != "" && mountDetails.pwd != "" {
			mountCommand.WriteString(" -o user=")
			mountCommand.WriteString(mountDetails.username)
			mountCommand.WriteString(",pass=")
			mountCommand.WriteString(mountDetails.pwd)
		}
		mountCommand.WriteString(" ")
		mountCommand.WriteString(mountDetails.address)
		mountCommand.WriteString(":")
		mountCommand.WriteString(mountDetails.src)
		mountCommand.WriteString(" ")
		mountCommand.WriteString(mountDetails.dest)
	default:
		mountCommand.WriteString("")
	}

	return mountCommand.String()
}

