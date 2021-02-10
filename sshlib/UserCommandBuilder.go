package sshlib

import (
	"fmt"
	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"io/ioutil"
	"strings"
)

func (userDetails *musthaveusersstruct) userManagementCommandBuilder(user *string, chosenOption string) string {
	userCommand := strings.Builder{}
	switch chosenOption { // TODO add "check" here later
	case "check":
		userCommand.WriteString(pkgmanlib.OmniTools["userinfo"] + *user)
	case "add":
		userCommand.WriteString(pkgmanlib.OmniTools["useradd"])
		userCommand.WriteString(" -g users ")
		if len(userDetails.groups) != 0 {
			userCommand.WriteString(" -G ")
			for comma, group := range userDetails.groups {
				userCommand.WriteString(group)
				if comma != len(userDetails.groups)-1 {
					userCommand.WriteString(",")
				}
			}
			if userDetails.sudoer != false {
				userCommand.WriteString(",wheel")
			} else {
			}
		}
		if userDetails.home != "" {
			userCommand.WriteString(" -d " + userDetails.home)
		}
		if userDetails.shell != "" {
			userCommand.WriteString(" -s " + userDetails.shell)
		}
		password := randomStringGenerator(8)
		userCommand.WriteString(" -p " + "\"" + password + "\" " + *user)
		err := ioutil.WriteFile("./config/passwords/"+*user, []byte(password), 0644)
		if err != nil {
			fmt.Printf("Error writing generated password to file, outputting now -> %s\n", password)
		}
	case "remove":
		userCommand.WriteString("killall -u " + *user + ";")
		userCommand.WriteString(pkgmanlib.OmniTools["userdel"] + *user)
	default:
		userCommand.WriteString("")
	}
	return userCommand.String()
}
