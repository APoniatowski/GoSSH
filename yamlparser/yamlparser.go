package yamlparser

// Global command-path progress counters, set by ParsePool and read by
// channelreaderlib and main.go.
//
// TODO: these remain package globals for now; de-globaling them is a separate
// later step (they feed the command-path progress bars).
var (
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
