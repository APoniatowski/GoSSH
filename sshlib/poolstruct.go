package sshlib

import (
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// ParsedPool holds one server's connection parameters in typed form. The
// pool.yml on-disk format is still loaded as a yaml.MapSlice; parseServer adapts
// that MapSlice into this typed struct so the rest of the code can use plain
// fields instead of interface{} assertions.
type ParsedPool struct {
	FQDN     string
	Username string
	Password string
	KeyPath  string
	Port     int
	OS       string
}

// normalizeKey lower-cases a key and strips underscores so that "Key_Path",
// "keypath" and "KeyPath" all collapse to the same canonical form.
func normalizeKey(key interface{}) string {
	ks, ok := key.(string)
	if !ok {
		return ""
	}
	return strings.ToLower(strings.ReplaceAll(ks, "_", ""))
}

// parseServer converts one server's yaml.MapSlice into a ParsedPool by matching
// key names case-insensitively and ignoring underscores. Field order is
// irrelevant and unmatched/missing keys are simply ignored, so a pool entry that
// omits optional fields (or reorders them) no longer index-panics. Port accepts
// either a YAML int or a numeric string; an absent/unparseable port is left 0 so
// defaulter() can fill it in. Callers are expected to invoke pp.defaulter()
// afterwards to fill in defaults.
func parseServer(sv yaml.MapSlice) ParsedPool {
	var pp ParsedPool
	for _, item := range sv {
		switch normalizeKey(item.Key) {
		case "fqdn":
			pp.FQDN, _ = item.Value.(string)
		case "username":
			pp.Username, _ = item.Value.(string)
		case "password":
			pp.Password, _ = item.Value.(string)
		case "keypath":
			pp.KeyPath, _ = item.Value.(string)
		case "port":
			pp.Port = parsePort(item.Value)
		case "os":
			pp.OS, _ = item.Value.(string)
		}
	}
	return pp
}

// parsePort normalizes a YAML port value to an int. Unquoted ints come through
// as int; quoted/string values are parsed best-effort; anything else (nil,
// unparseable) yields 0 so defaulter() applies the default.
func parsePort(v interface{}) int {
	switch p := v.(type) {
	case int:
		return p
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return 0
		}
		return n
	default:
		return 0
	}
}
