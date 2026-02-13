package git

import (
	"fmt"
	"os"
	"os/exec"
)

// CommitToGit writes the commit string to a temp file and runs `git commit -F`
func Commit(commit string) error {
	tmp, err := os.CreateTemp("", "gitcommit*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(commit); err != nil {
		return fmt.Errorf("failed to write commit message: %w", err)
	}
	tmp.Close()

	cmd := exec.Command("git", "commit", "-F", tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}
	return nil
}
