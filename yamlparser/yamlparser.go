package yamlparser

import "gopkg.in/yaml.v2"

// Global vars to be able to access it in all packages
var (
	Baseline        yaml.MapSlice
	Pool            yaml.MapSlice
	Waittotal       int
	Grouptotal      int
	ServersPerGroup []int
)

/*
structure:

yamlparser.go
	poolparser.go
	baselineparser.go
	counters.go
*/
