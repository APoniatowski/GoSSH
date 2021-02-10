package sshlib

import (
	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"math/rand"
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

func randomStringGenerator(strLength int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+-=!@#$%^&*")

	str := make([]rune, strLength)
	for i := range str {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}

