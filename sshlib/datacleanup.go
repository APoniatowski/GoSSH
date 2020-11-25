package sshlib

import (
	"strconv"
)

//defaulter defaults all empty fields in yaml file and to abort if too many values are missing, eg password and key_path
func (pp *ParsedPool) defaulter() {
	if pp.password == nil && pp.keypath == nil {
		panic("Both 'Password' and 'Key_Path' fields are empty... Aborting.\n")
	}
	if pp.username == nil {
		pp.username = "root"
	}
	if pp.password == nil {
		pp.password = ""
	}
	if pp.keypath == nil {
		pp.keypath = ""
	}
	if pp.port == nil {
		pp.port = 22
		pp.port = strconv.Itoa(pp.port.(int))
	} else {
		pp.port = strconv.Itoa(pp.port.(int))
	}
}
