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
			fetchURLCommand.WriteString(prereqCleanup(&parsedDirName))
		case "gitlab":
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
			fetchURLCommand.WriteString(" " + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(&parsedDirName))
		case "bitbucket":
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
			fetchURLCommand.WriteString(" " + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(&parsedDirName))
		case "gerrit":
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
			fetchURLCommand.WriteString(" " + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(&parsedDirName))
		case "git":
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["git"])
			fetchURLCommand.WriteString(" " + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(&parsedDirName))
		case "svn":
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["svn"])
			fetchURLCommand.WriteString(*url)
			fetchURLCommand.WriteString(" " + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(&parsedDirName))
		default:
			// cleanup: download to /tmp (so it can be cleaned up later) and
			// symlink it into the home dir for use during the run.
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["curl"])
			fetchURLCommand.WriteString(*url + " -o /tmp/" + parsedDirName)
			fetchURLCommand.WriteString(" || ")
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["wget"])
			fetchURLCommand.WriteString(*url)
			fetchURLCommand.WriteString(" -O /tmp/" + parsedDirName + ";")
			fetchURLCommand.WriteString(prereqCleanup(&parsedDirName))
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
			// no cleanup: save the file into the current (home) dir by basename.
			// curl needs -O to write a file; wget saves by basename by default.
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["curl"])
			fetchURLCommand.WriteString(*url + " -O")
			fetchURLCommand.WriteString(" || ")
			fetchURLCommand.WriteString(pkgmanlib.OmniTools["wget"])
			fetchURLCommand.WriteString(*url)
		}
	}
	return fetchURLCommand.String()
}

func prereqCleanup(dirName *string) string {
	// Symlink the /tmp copy into the home dir: `ln -sfn /tmp/<name> ~/<name>`.
	// The space before ~/ is required — without it this collapses into a single
	// malformed argument and the symlink clobbers the real download.
	symlinker := strings.Builder{}
	symlinker.WriteString("ln -sfn /tmp/")
	symlinker.WriteString(*dirName)
	symlinker.WriteString(" ~/")
	symlinker.WriteString(*dirName)
	return symlinker.String()
}

func prereqRemoteCleanup(dirName, fileName *string, cleanup *bool) string {
	symlinker := strings.Builder{}
	var mountedPath string
	if *fileName == "" || strings.ToLower(*fileName) == "all" {
		mountedPath = *dirName
	} else {
		mountedPath = *dirName + *fileName
	}
	if *cleanup {
		symlinker.WriteString(";ln -sfn ")
		symlinker.WriteString(mountedPath)
	} else {
		symlinker.WriteString(";cp -a ")
		symlinker.WriteString(mountedPath)
	}
	symlinker.WriteString(" ~/.")
	return symlinker.String()
}
