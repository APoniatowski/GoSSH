package main

import (
	"fmt"
	"log"

	ymlp "github.com/APoniatowski/GoSSH/yamlparser"
)

// Error checking function
func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

// Main function to carry out operations
func main() {
	test := ymlp.ParseYAML()
	fmt.Printf("%v      %T\n", test, test)
	// check(err)
	// will add the funcs from the lib, once I have it setup with the proper args... when I have the time
}
