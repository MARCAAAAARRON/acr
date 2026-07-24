package gitctx

import (
	"os/exec"
	"strings"
)

// ModifiedFiles returns the list of files that are modified, staged, or
// untracked in repoPath, using `git status --porcelain`. If repoPath isn't a
// git repo, or git isn't available, it returns nil and no error — callers
// should treat this as "no git awareness available," not a failure.
func ModifiedFiles(repoPath string) []string {
	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	rawLines := strings.Split(string(output), "\n")
	var files []string
	for _, line := range rawLines {
		line = strings.TrimRight(line, "\r") // strip Windows line-ending remnants only
		if len(line) > 3 {
			files = append(files, strings.TrimSpace(line[3:]))
		}
	}
	return files
}