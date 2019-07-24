package main

import (
	"log"
	ymlp "yamlparser"
)

// Error checking function
func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

// Main function to carry out operations
func main() {
	ymlp.ParseYAML()
	// check(err)
	// will add the funcs from the lib, once I have it setup with the proper args... when I have the time
}
