package git

import (
	"os/exec"
	"strings"
	"unicode"
)

// GitStatus gets the git status of the repos directory
func GitStatus() ([]string, error) {
	var (
		err      error
		modfiles []string
	)

	// Git status command
	// gsArgs := []string{"diff", "--staged", "--name-status"}
	// gitStatus := exec.Command("git", gsArgs...)
	gsArgs := []string{"status", "--short"}
	gitStatus := exec.Command("git", gsArgs...)

	// Get stdout and trim the empty last index
	fileStatus, err := gitStatus.CombinedOutput()
	if err != nil {
		return nil, err
	}

	fsSplit := strings.Split(string(fileStatus), "\n")

	for _, status := range fsSplit {

		s := strings.Fields(status)
		if len(s) == 0 {
			continue
		}

		// With the "git status --short" command the staged files are at 0 index.
		// Unstaged: " M filename" has a space before the M.
		// Staged: "M filename" has no space before the M.
		//
		// If there is a space we do not want to include it in the commit.
		if unicode.IsSpace(rune(status[0])) {
			continue
		}

		// With the "git status --short" command the staged files are at 0 index.
		// Unstaged: "?? directory/"
		//
		// If there is a question mark that means that it is untracked and unstaged.
		if unicode.IsPunct(rune(status[0])) {
			continue
		}

		// If the file was deleted there is no reason to read that file.
		if strings.Contains(s[0], "D") {
			continue
		}

		// Only get the file path.
		modfiles = append(modfiles, s[len(s)-1])
	}

	return modfiles, err
}
