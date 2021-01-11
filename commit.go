package main

import "os"

func GetCommits(modifiedFiles []string) []string {
	out := []string{}

	_, err := os.Stat(".commit")
	if !os.IsNotExist(err) {
		out = append(out, ReadFile(".commit")...)
	}

	// Parse Files
	for _, filename := range modifiedFiles {
		UnStageFile(filename)
		out = append(out, ReadFile(filename)...)
		StageFile(filename)
	}

	return out
}
