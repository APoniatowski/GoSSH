package sshlib

import (
	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"strings"
)

func firewallCommandBuilder(port, protocol, zone *string, chosenOption string) string {
	// TODO chang awk to grep and add another parameter for open/deny/closed/etc and add another option for open and closed rules
	fwCommand := strings.Builder{}
	protocolSlice := strings.Split(*protocol, " ")
	const orOperator string = " || "
	switch chosenOption {
	case "check":
		for i := range protocolSlice {
			fwCommand.WriteString(pkgmanlib.Firewalld["list"])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString(pkgmanlib.Ufw["list"])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString(pkgmanlib.Iptables["list"])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString(pkgmanlib.Nftables["list"])
			//fwCommand.WriteString(orOperator)
			//fwCommand.WriteString(pkgmanlib.PfFirewall["list"])
			fwCommand.WriteString(" > ")
			fwCommand.WriteString(pkgmanlib.OmniTools["awk"])
			fwCommand.WriteString("'/" + *port + "/' && '/" + protocolSlice[i] + "/'")
			fwCommand.WriteString(";")
		}

	case "apply-open":
		for i := range protocolSlice {
			fwCommand.WriteString("firewall-cmd --zone=")
			if *zone == ""{
				fwCommand.WriteString("$(firewall-cmd --get-default-zone)")
			} else{
				fwCommand.WriteString(*zone)
			}
			// if to check protocol, if both then udp and tcp and none, default to tcp
			fwCommand.WriteString(" --add-port=" + *port + "/" + protocolSlice[i])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("ufw allow " + *port + "/" + protocolSlice[i])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("iptables -A INPUT -p " + protocolSlice[i] + " --dport " + *port + " -j ACCEPT")
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("nft add rule ip filter input " + protocolSlice[i] + " dport " + *port + " ACCEPT;")
			//fwCommand.WriteString(orOperator)
			// pf/ipfw too complex for simple commands
			// I will need to add OS specific checks to add a script to add rules, due to rule number/order
		}
		fwCommand.WriteString("iptables-save")

	case "apply-closed":
		for i := range protocolSlice {
			fwCommand.WriteString("firewall-cmd --zone=")
			if *zone == ""{
				fwCommand.WriteString("$(firewall-cmd --get-default-zone)")
			} else{
				fwCommand.WriteString(*zone)
			}
			// if to check protocol, if both then udp and tcp and none, default to tcp
			fwCommand.WriteString(" --remove-port=" + *port + "/" + protocolSlice[i])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("ufw deny " + *port + "/" + protocolSlice[i])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("iptables -A INPUT -p " + protocolSlice[i] + " --dport " + *port + " -j DROP")
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("nft add rule ip filter input " + protocolSlice[i] + " dport " + *port + " DROP;")
			//fwCommand.WriteString(orOperator)
			// pf/ipfw too complex for simple commands
			// I will need to add OS specific checks to add a script to add rules, due to rule number/order
		}
		fwCommand.WriteString("iptables-save")

	case "remove-open":
		for i := range protocolSlice {
			fwCommand.WriteString("firewall-cmd --zone=")
			if *zone == ""{
				fwCommand.WriteString("$(firewall-cmd --get-default-zone)")
			} else{
				fwCommand.WriteString(*zone)
			}
			// if to check protocol, if both then udp and tcp and none, default to tcp
			fwCommand.WriteString(" --remove-port=" + *port + "/" + protocolSlice[i])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("ufw deny " + *port + "/" + protocolSlice[i])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("iptables -A INPUT -p " + protocolSlice[i] + " --dport " + *port + " -j DROP")
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("nft add rule ip filter input " + protocolSlice[i] + " dport " + *port + " DROP;")
			//fwCommand.WriteString(orOperator)
			// pf/ipfw too complex for simple commands
			// I will need to add OS specific checks to add a script to add rules, due to rule number/order
		}
		fwCommand.WriteString("iptables-save")

	case "remove-closed":
		for i := range protocolSlice {
			fwCommand.WriteString("firewall-cmd --zone=")
			if *zone == ""{
				fwCommand.WriteString("$(firewall-cmd --get-default-zone)")
			} else{
				fwCommand.WriteString(*zone)
			}
			// if to check protocol, if both then udp and tcp and none, default to tcp
			fwCommand.WriteString(" --add-port=" + *port + "/" + protocolSlice[i] + " --permanent")
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("ufw allow " + *port + "/" + protocolSlice[i])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("iptables -A INPUT -p " + protocolSlice[i] + " --dport " + *port + " -j ACCEPT")
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("nft add rule ip filter input " + protocolSlice[i] + " dport " + *port + " ACCEPT;")
			//fwCommand.WriteString(orOperator)
			// pf/ipfw too complex for simple commands
			// I will need to add OS specific checks to add a script to add rules, due to rule number/order
		}
		fwCommand.WriteString("iptables-save")

	default:
		fwCommand.WriteString("")
	}

	return fwCommand.String()
}

