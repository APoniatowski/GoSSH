package main

import (
	"log"
	"src/github.com/aponiatowski/gossh/lib/yamlparser" // need to look into this, as I would like to have separate libs, this is easier in python and rust
)

// Error checking function
func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

// Main function to carry out operations
func main() {
	infotoProcess, err := yamlParser.ParseYAML()
	check(err)
	// will add the funcs from the lib, once I have it setup with the proper args... when I have the time
}
