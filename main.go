package main

import (
	"fmt"
	"log"
	"reflect"

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
	configs := ymlp.ParseYAML()
	v := reflect.ValueOf(configs)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	// configmap, err := ymlp.InterfaceToMap(configs)
	// check(err)
	fmt.Printf("%v \n  %T \n", values, values)

	// fmt.Printf("%v \n  %T \n", configmap, configmap)
}
