package main

import (
	"fmt"
	"os/exec"
)

func UnStageFile(file string) {
	if file == "" {
		return
	}
	gsArgs := []string{"reset", "HEAD", file}

	gitCommand := exec.Command("git", gsArgs...)
	err := gitCommand.Run()
	if err != nil {
		fmt.Println(fmt.Sprintf("could not unstage %s"))
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
		fmt.Println(fmt.Sprintf("could not stage %s"))
		return
	}
}
