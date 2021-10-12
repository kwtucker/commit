package git

import (
	"os/exec"
)

func UnStageFile(file string) {
	if file == "" {
		return
	}

	gsArgs := []string{"reset", "--", file}
	gitCommand := exec.Command("git", gsArgs...)

	err := gitCommand.Run()
	if err != nil {
		return
	}
}

func StageFile(file string) {
	if file == "" {
		return
	}

	gsArgs := []string{"add", file}

	gitCommand := exec.Command("git", gsArgs...)
	err := gitCommand.Run()
	if err != nil {
		return
	}
}
