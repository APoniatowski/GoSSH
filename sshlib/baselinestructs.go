package sshlib

// ParsedBaseline - top level struct
type ParsedBaseline struct {
	// servergroups []string
	exclude     exclusion
	prereq      prereqs
	musthave    musthaves
	mustnothave mustnothaves
	final       finals
}

// exclusion struct
type exclusion struct {
	osExcl      []string
	serversExcl []string
}

// prerequisite struct
type prereqs struct {
	tools   []string
	files   prereqsfiles
	vcs     prereqvcs
	script  string
	cleanup bool
}

type prereqsfiles struct {
	urls   []string
	local  fileslocal
	remote filesremote
}

type fileslocal struct {
	src  string
	dest string
}

type filesremote struct {
	mounttype string
	address   string
	username  string
	pwd       string
	src       string
	dest      string
	files     []string
}

type prereqvcs struct {
	urls    []string
	execute []string
}

// musthaves struct
type musthaves struct {
	installed  []string
	enabled    []string
	disabled   []string
	configured musthaveconfigured
	users      musthaveusers
	policies   musthavepolicies
	rules      musthaverules
	mounts     musthavemounts
}

type musthaveconfigured struct {
	services map[string]musthaveconfiguredservices
}

type musthaveconfiguredservices struct {
	source      []string
	destination []string
}

type musthaveusers struct {
	users map[string]musthaveusersstruct
}

type musthaveusersstruct struct {
	groups []string
	shell  string
	home   string
	sudoer bool
}

type musthavepolicies struct {
	polstatus string
	polimport string
	polreboot bool
}
type musthaverules struct {
	fwopen   musthaverulesopen
	fwclosed musthaverulesclosed
	fwzones  []string
}

type musthaverulesopen struct {
	ports     []string
	protocols []string
}

type musthaverulesclosed struct {
	ports     []string
	protocols []string
}

type musthavemounts struct {
	mountname map[string]mountdetails
}

type mountdetails struct {
	mounttype string
	address   string
	username  string
	pwd       string
	src       string
	dest      string
}

// mustnothaves struct
type mustnothaves struct {
	installed []string
	enabled   []string
	disabled  []string
	users     []string
	rules     mustnothaverules
	mounts    []string
}

type mustnothaverules struct {
	fwopen   mustnothaverulesopen
	fwclosed mustnothaverulesclosed
	fwzones  []string
}

type mustnothaverulesopen struct {
	ports     []string
	protocols []string
}

type mustnothaverulesclosed struct {
	ports     []string
	protocols []string
}

// finals struct
type finals struct {
	scripts  []string
	commands []string
	collect  collections
	restart  restarts
}

type collections struct {
	logs  []string
	stats []string
	files []string
	users bool
}

type restarts struct {
	services bool
	servers  bool
}
