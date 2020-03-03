package sshlib

// top level struct
type ParsedBaseline struct {
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
	tools []string
	files prereqsfiles
	vcs   prereqvcs
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
}

type prereqvcs struct {
	urls    []string
	execute []string
	script  string
	cleanup bool
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
	services interface{} // map[string]string - string values
}
type musthaveusers struct {
	users interface{} // map[string]string - string values
}
type musthavepolicies struct {
	polstatus string
	polimport string
	polreboot bool
}
type musthaverules struct {
	fwtype   string
	fwimport string
	fwopen   []string
	fwclose  []string
	fwzones  []string
}
type musthavemounts struct {
	mountname mountdetails
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
	fwtype  string
	fwopen  []string
	fwclose []string
	fwzones []string
}

// finals struct
type finals struct {
	scripts  []string
	commands []string
	collect  []string
	restart  bool
}
