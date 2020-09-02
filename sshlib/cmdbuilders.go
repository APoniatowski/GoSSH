package sshlib

import (
	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"strings"
)

// Switches For checking what CLI option was used and run the appropriate functions
type Switches struct {
	Updater, UpdaterFull, Install, Uninstall *bool
}

// OSSwitcher a much needed var between main and sshlib
var OSSwitcher Switches

// validator This needs to be outside of the function for extra error handling
var validator string

// recoveries keeping count of recoveries. Might be useful later
var recoveries interface{}

//Switcher Method to check the switches set for each respective action (update/install/uninstall)
func (S *Switches) Switcher(pp ParsedPool, command string) (rtncommand string) {
	if *S.Updater {
		rtncommand = pkgmanlib.Update(pp.username.(string), pp.os.(string))
	}
	if *S.UpdaterFull {
		rtncommand = pkgmanlib.UpdateOS(pp.username.(string), pp.os.(string))
	}
	if *S.Install {
		rtncommand = pkgmanlib.Install(pp.username.(string), pp.os.(string)) + command + " -y 2>&1"
	}
	if *S.Uninstall {
		rtncommand = pkgmanlib.Uninstall(pp.username.(string), pp.os.(string)) + command + " -y 2>&1"
	}
	return
}

func prereqURLFetch(url string) string {
	fetchURLCommand := strings.Builder{}
	stripSlashURL := strings.Split(url, "/")
	parsedURL := strings.Split(stripSlashURL[2], ".")
	var checkURL string
	if parsedURL[0] == "www" {
		checkURL = parsedURL[1]
	} else {
		checkURL = parsedURL[0]
	}
	switch checkURL {
	case "github":
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
		fetchURLCommand.WriteString(url)
	case "gitlab":
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
		fetchURLCommand.WriteString(url)
	case "bitbucket":
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
		fetchURLCommand.WriteString(url)
	case "gerrit":
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
		fetchURLCommand.WriteString(url)
	case "git":
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
		fetchURLCommand.WriteString(url)
	case "svn":
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["svn"])
		fetchURLCommand.WriteString(url)
	default:
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["curl"])
		fetchURLCommand.WriteString(url)
		fetchURLCommand.WriteString(" || ")
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["wget"])
		fetchURLCommand.WriteString(url)
	}
	return fetchURLCommand.String()
}

func (remoteMount *filesremote) remoteFilesCommandBuilder(chosenOption,file string) string {
	remoteMountCommand := strings.Builder{}
	switch chosenOption {
	case "check":
		remoteMountCommand.WriteString("test")
		remoteMountCommand.WriteString(remoteMount.dest)
		remoteMountCommand.WriteString(file)
		remoteMountCommand.WriteString(" && echo \"1\"")
	case "apply":
		remoteMountCommand.WriteString("mount -F ")
		remoteMountCommand.WriteString(remoteMount.mounttype)
		remoteMountCommand.WriteString(" -o user=")
		remoteMountCommand.WriteString(remoteMount.username)
		remoteMountCommand.WriteString(",pass=")
		remoteMountCommand.WriteString(remoteMount.pwd)
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

func serviceCommandBuilder(service *string, chosenOption,serviceOption,osOption string) string {
	serviceCommand := strings.Builder{}
	switch chosenOption {
	case "check":
		switch serviceOption {
		case "install":
			serviceCommand.WriteString(pkgmanlib.PkgInstall[osOption] + *service)
		case "uninstall":

		case "enabled":

		case "disabled":

		default:
		}
	case "apply":
		switch serviceOption {
		case "install":
			serviceCommand.WriteString("")
		case "uninstall":

		case "enabled":

		case "disabled":

		default:
		}
	default:
		serviceCommand.WriteString("")
	}
	return serviceCommand.String()
}

func firewallCommandBuilder(port, protocol *string, chosenOption string) string {
	// TODO chang awk to grep and add another parameter for open/deny/closed/etc
	fwCommand := strings.Builder{}
	switch chosenOption {
	case "check":
		fwCommand.WriteString("{ ")
		fwCommand.WriteString(pkgmanlib.Firewalld["list"])
		fwCommand.WriteString(" || ")
		fwCommand.WriteString(pkgmanlib.Ufw["list"])
		fwCommand.WriteString(" || ")
		fwCommand.WriteString(pkgmanlib.Iptables["list"])
		fwCommand.WriteString(" || ")
		fwCommand.WriteString(pkgmanlib.Nftables["list"])
		fwCommand.WriteString(" || ")
		fwCommand.WriteString(pkgmanlib.PfFirewall["list"])
		fwCommand.WriteString(" } ")
		fwCommand.WriteString(" > ")
		fwCommand.WriteString(pkgmanlib.OmniTools["awk"])
		fwCommand.WriteString("'/" + *port + "/")
		fwCommand.WriteString(" && ")
		fwCommand.WriteString("'/" + *protocol + "/'")
	case "apply":

	default:
		fwCommand.WriteString("")
	}

	return fwCommand.String()
}

func (mountDetails *mountdetails) mountCommandBuilder(chosenOption string) string {
	mountCommand := strings.Builder{}
	switch chosenOption {
	case "check":

	case "apply":
		mountCommand.WriteString(pkgmanlib.OmniTools["mkdir"])
		mountCommand.WriteString(mountDetails.dest + " && ")
		mountCommand.WriteString("echo '")
		mountCommand.WriteString(mountDetails.address + ":" + mountDetails.src + " " + mountDetails.dest + " ")
		mountCommand.WriteString(mountDetails.mounttype)
		mountCommand.WriteString(" defaults 0 0") // Default mounting details
		mountCommand.WriteString("' >> /etc/fstab;")
		mountCommand.WriteString(pkgmanlib.OmniTools["mount"] + mountDetails.dest)
	default:
		mountCommand.WriteString("")
	}

	return mountCommand.String()
}
