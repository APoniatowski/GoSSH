package main

import (
	"log"

	"github.com/APoniatowski/GoSSH/sshlib"
	"github.com/APoniatowski/GoSSH/yamlparser"
)

// Error checking function
func generalError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// Main function to carry out operations
func main() {
	yamlparser.Rollcall()
	newConfig := yamlparser.Config
	// sshlib.RunSequentially(&newConfig)
	sshlib.RunServersConcurrently(&newConfig)
	//TODO /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}
