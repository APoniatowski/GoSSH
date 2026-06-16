package sshlib

import (
	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"strings"
)

func serviceCommandBuilder(service, osOption *string, serviceOption string) string {
	serviceCommand := strings.Builder{}
	switch serviceOption {
	case "search":
		serviceCommand.WriteString(pkgmanlib.PkgSearch[*osOption] + *service)
	case "install":
		serviceCommand.WriteString(pkgmanlib.PkgInstallYes[*osOption] + *service)
	case "uninstall":
		serviceCommand.WriteString(pkgmanlib.PkgUninstallYes[*osOption] + *service)
	case "enable":
		serviceCommand.WriteString(pkgmanlib.OmniTools["systemctlenable"] + *service)
	case "disable":
		serviceCommand.WriteString(pkgmanlib.OmniTools["systemctldisable"] + *service)
	case "isactive":
		serviceCommand.WriteString(pkgmanlib.OmniTools["serviceisactive"] + *service)
	default:
		serviceCommand.WriteString("")
	}
	return serviceCommand.String()
}
