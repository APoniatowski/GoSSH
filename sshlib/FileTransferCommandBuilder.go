package sshlib

import (
	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"strings"
)

func (remoteMount *filesremote) remoteFilesCommandBuilder(file *string, chosenOption string) string {
	remoteMountCommand := strings.Builder{}
	switch chosenOption {
	case "check":
		remoteMountCommand.WriteString("test")
		remoteMountCommand.WriteString(remoteMount.dest)
		remoteMountCommand.WriteString(*file)
		remoteMountCommand.WriteString(" && echo \"1\"")
	case "apply":
		remoteMountCommand.WriteString(pkgmanlib.OmniTools["mount"] + "-F ")
		remoteMountCommand.WriteString(remoteMount.mounttype)
		if remoteMount.username != "" && remoteMount.pwd != "" {
			remoteMountCommand.WriteString(" -o user=")
			remoteMountCommand.WriteString(remoteMount.username)
			remoteMountCommand.WriteString(",pass=")
			remoteMountCommand.WriteString(remoteMount.pwd)
		}
		remoteMountCommand.WriteString(" ")
		remoteMountCommand.WriteString(remoteMount.address)
		remoteMountCommand.WriteString(":")
		remoteMountCommand.WriteString(remoteMount.src)
		remoteMountCommand.WriteString(" ")
		remoteMountCommand.WriteString(remoteMount.dest)
	default:
	}
	return remoteMountCommand.String()
}
