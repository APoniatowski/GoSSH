package yamlparser

import (
	"log"
)

// Error checking function
func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}
