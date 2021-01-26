package main

import (
	"fmt"
	"os"
)

func GetCommits(modifiedFiles []string) []string {
	out := []string{}

	_, err := os.Stat(".commit")
	if !os.IsNotExist(err) {
		out = append(out, ReadFile(".commit")...)
	}

	// Parse Files
	for _, filename := range modifiedFiles {
		fi, err := os.Stat(filename)
		if err != nil {
			fmt.Println(err)
			return out
		}
		switch mode := fi.Mode(); {
		case mode.IsDir():
			// TODO: cd into the directory and read files
			// do directory stuff
			fmt.Println("directory")
		case mode.IsRegular():
			UnStageFile(filename)
			out = append(out, ReadFile(filename)...)
			StageFile(filename)
		}
	}

	return out
}
