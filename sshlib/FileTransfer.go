package sshlib

import (
	"encoding/base64"
	"os"
	"path"
)

// fileToCommand reads a local file and returns a single shell command that
// recreates that file at remoteDest on the target. The file content travels
// inside the command string as base64, so no separate transport (scp/sftp) and
// no new dependency is required: the normal command channel carries everything.
//
// The command is a single line (echo '<b64>' | base64 -d | tee) so callers can
// safely chain more steps with `&&` (e.g. Final Scripts appends chmod + exec).
// An earlier heredoc form broke that: appended text landed on the heredoc
// delimiter line, so the delimiter never matched and chmod/exec were swallowed.
// The redirect uses tee (not `>`) so that, when sudo is requested, the
// privileged write happens in the sudo'd process rather than the calling shell.
// All paths are single-quoted to survive spaces and shell metacharacters.
func fileToCommand(localSrc, remoteDest string, sudo bool) (string, error) {
	data, err := os.ReadFile(localSrc)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	dir := path.Dir(remoteDest)

	pfx := ""
	if sudo {
		pfx = "sudo "
	}

	return pfx + "mkdir -p '" + dir + "' && echo '" + encoded + "' | base64 -d | " +
		pfx + "tee '" + remoteDest + "' >/dev/null", nil
}
