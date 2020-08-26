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
		fwCommand.WriteString("'/"+*port+"/")
		fwCommand.WriteString(" && ")
		fwCommand.WriteString("'/"+*protocol+"/'")
	case "apply":

	default:
		fwCommand.WriteString("")
	}

	return fwCommand.String()
}

func (mountDetails *mountdetails) mountCommandBuilder() string {
	mountCommand := strings.Builder{}
	mountCommand.WriteString(pkgmanlib.OmniTools["mkdir"])
	mountCommand.WriteString(mountDetails.dest + " && ")
	mountCommand.WriteString("echo '")
	mountCommand.WriteString(mountDetails.address + ":" + mountDetails.src + " " + mountDetails.dest + " ")
	mountCommand.WriteString(mountDetails.mounttype)
	mountCommand.WriteString(" defaults 0 0") // Default mounting details
	mountCommand.WriteString("' >> /etc/fstab;")
	mountCommand.WriteString(pkgmanlib.OmniTools["mount"] + mountDetails.dest)
	return mountCommand.String()
}