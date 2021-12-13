package sshlib

import (
	"github.com/APoniatowski/GoSSH/pkgmanlib"
	"strings"
)

func prereqURLFetch(url *string, cleanup *bool) string {
	fetchURLCommand := strings.Builder{}
	stripSlashURL := strings.Split(*url, "/")
	stripDirName := strings.Split(*url, "/")
	parsedURL := strings.Split(stripSlashURL[2], ".")
	parsedDirName := stripDirName[len(stripDirName)-1]
	var checkURL string
	if parsedURL[0] == "www" {
		checkURL = parsedURL[1]
	} else {
		checkURL = parsedURL[0]
	}
	if *cleanup {
		switch checkURL {
		case "github":
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
			fetchURLCommand.WriteString(" " + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(parsedDirName))
		case "gitlab":
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
			fetchURLCommand.WriteString(" " + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(parsedDirName))
		case "bitbucket":
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
			fetchURLCommand.WriteString(" " + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(parsedDirName))
		case "gerrit":
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
			fetchURLCommand.WriteString(" " + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(parsedDirName))
		case "git":
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
			fetchURLCommand.WriteString(" " + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(parsedDirName))
		case "svn":
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["svn"])
			fetchURLCommand.WriteString(*url)
			fetchURLCommand.WriteString(" " + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(parsedDirName))
		default:
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["curl"])
			fetchURLCommand.WriteString(*url + " > ")
			fetchURLCommand.WriteString(parsedDirName)
			fetchURLCommand.WriteString(" || ")
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["wget"])
			fetchURLCommand.WriteString(*url)
			fetchURLCommand.WriteString(" -O " + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(parsedDirName))
		}
	} else {
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
	}
	return fetchURLCommand.String()
}

func prereqCleanup(dirName string) string {
	symlinker := strings.Builder{}
	symlinker.WriteString("ln -sfn /tmp/")
	symlinker.WriteString(dirName)
	symlinker.WriteString("~/")
	symlinker.WriteString(dirName)
	return symlinker.String()
}
