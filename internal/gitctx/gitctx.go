package gitctx

import (
	"os/exec"
	"strings"
)

// ModifiedFiles returns the list of files that are modified, staged, or
// untracked in repoPath, using `git status --porcelain`. If repoPath isn't a
// git repo, or git isn't available, it returns nil and no error.
func ModifiedFiles(repoPath string) []string {
	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}
	return parsePorcelain(string(output))
}

// parsePorcelain parses `git status --porcelain` output into a list of file
// paths. Split out from ModifiedFiles so the parsing logic can be tested
// directly with hand-written input, without needing a real git repo.
func parsePorcelain(output string) []string {
	rawLines := strings.Split(output, "\n")
	var files []string
	for _, line := range rawLines {
		line = strings.TrimRight(line, "\r")
		if len(line) > 3 {
			files = append(files, strings.TrimSpace(line[3:]))
		}
	}
	return files
}