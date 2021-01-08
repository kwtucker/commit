package main

import (
	"os/exec"
	"strings"
)

// GitStatus gets the git status of the repos directory
func GitStatus() ([]string, error) {
	var (
		err      error
		modfiles []string
	)

	// Git status command
	gsArgs := []string{"diff", "--staged", "--name-status"}
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
		// If the file was deleted there is no reason to read that file.
		if strings.Contains(s[0], "D") {
			continue
		}

		// Only get the file path.
		modfiles = append(modfiles, s[len(s)-1])
	}

	return modfiles, err
}
