package sshlib

//ParsedPool Parsing data to struct. Mandatory data
type ParsedPool struct {
	fqdn     interface{}
	username interface{}
	password interface{}
	keypath  interface{}
	port     interface{}
	os       interface{}
}
