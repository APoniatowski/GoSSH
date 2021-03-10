package sshlib

import (
	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"strings"
)

func prereqURLFetch(url *string) string {
	fetchURLCommand := strings.Builder{}
	stripSlashURL := strings.Split(*url, "/")
	parsedURL := strings.Split(stripSlashURL[2], ".")
	var checkURL string
	if parsedURL[0] == "www" {
		checkURL = parsedURL[1]
	} else {
		checkURL = parsedURL[0]
	}
	switch checkURL {
	case "github":
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
		fetchURLCommand.WriteString(*url)
	case "gitlab":
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
		fetchURLCommand.WriteString(*url)
	case "bitbucket":
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
		fetchURLCommand.WriteString(*url)
	case "gerrit":
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
		fetchURLCommand.WriteString(*url)
	case "git":
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
		fetchURLCommand.WriteString(*url)
	case "svn":
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["svn"])
		fetchURLCommand.WriteString(*url)
	default:
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["curl"])
		fetchURLCommand.WriteString(*url)
		fetchURLCommand.WriteString(" || ")
		fetchURLCommand.WriteString(pkgmanlib.OmniTools["wget"])
		fetchURLCommand.WriteString(*url)
	}
	return fetchURLCommand.String()
}

