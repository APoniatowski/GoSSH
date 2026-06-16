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
		// For each requested protocol, run whichever firewall lister is
		// available (suppressing stderr so missing tools don't fail the
		// pipeline) and grep its output for the port/protocol. grep exits 0
		// on a match and nonzero otherwise, so GoSSH reports OK when the rule
		// is present and NOK when it is absent. Per-protocol checks are
		// chained with && so the whole command only succeeds when every
		// requested protocol is found.
		nftlist := pkgmanlib.Nftables["list"]
		if nftlist == "" {
			nftlist = "nft list ruleset"
		}
		for i := range protocolSlice {
			if i > 0 {
				fwCommand.WriteString(" && ")
			}
			fwCommand.WriteString("{ ")
			fwCommand.WriteString(pkgmanlib.Firewalld["list"] + " 2>/dev/null")
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString(pkgmanlib.Ufw["list"] + " 2>/dev/null")
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString(pkgmanlib.Iptables["list"] + " 2>/dev/null")
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString(nftlist + " 2>/dev/null")
			//fwCommand.WriteString(orOperator)
			//fwCommand.WriteString(pkgmanlib.PfFirewall["list"] + " 2>/dev/null")
			fwCommand.WriteString("; } | ")
			fwCommand.WriteString(pkgmanlib.OmniTools["grep"])
			fwCommand.WriteString("-E '" + *port + "(/" + protocolSlice[i] + "|.*" + protocolSlice[i] + ")'")
		}
		fwCommand.WriteString(";")

	case "apply-open":
		for i := range protocolSlice {
			fwCommand.WriteString("firewall-cmd --zone=")
			if *zone == "" {
				fwCommand.WriteString("$(firewall-cmd --get-default-zone)")
			} else {
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
			if *zone == "" {
				fwCommand.WriteString("$(firewall-cmd --get-default-zone)")
			} else {
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
		// Delete an existing OPEN/allow rule (do NOT add a deny).
		for i := range protocolSlice {
			fwCommand.WriteString("firewall-cmd --zone=")
			if *zone == "" {
				fwCommand.WriteString("$(firewall-cmd --get-default-zone)")
			} else {
				fwCommand.WriteString(*zone)
			}
			// if to check protocol, if both then udp and tcp and none, default to tcp
			fwCommand.WriteString(" --remove-port=" + *port + "/" + protocolSlice[i])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("ufw delete allow " + *port + "/" + protocolSlice[i])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("iptables -D INPUT -p " + protocolSlice[i] + " --dport " + *port + " -j ACCEPT")
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("nft delete rule ip filter input " + protocolSlice[i] + " dport " + *port + " ACCEPT;")
			//fwCommand.WriteString(orOperator)
			// pf/ipfw too complex for simple commands
			// I will need to add OS specific checks to add a script to add rules, due to rule number/order
		}
		fwCommand.WriteString("iptables-save")

	case "remove-closed":
		// Delete an existing CLOSED/deny rule (do NOT add an allow).
		for i := range protocolSlice {
			fwCommand.WriteString("firewall-cmd --zone=")
			if *zone == "" {
				fwCommand.WriteString("$(firewall-cmd --get-default-zone)")
			} else {
				fwCommand.WriteString(*zone)
			}
			// if to check protocol, if both then udp and tcp and none, default to tcp
			fwCommand.WriteString(" --remove-port=" + *port + "/" + protocolSlice[i])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("ufw delete deny " + *port + "/" + protocolSlice[i])
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("iptables -D INPUT -p " + protocolSlice[i] + " --dport " + *port + " -j DROP")
			fwCommand.WriteString(orOperator)
			fwCommand.WriteString("nft delete rule ip filter input " + protocolSlice[i] + " dport " + *port + " DROP;")
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
