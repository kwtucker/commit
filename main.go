package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	norepository             = "No repository in current directory"
	ignoredDefaultConfigFile = "forgit.json"
)

var Configuration *Config

func main() {

	Configuration = &Config{
		Commit:       &Commit{},
		Clean:        &Clean{},
		IgnoredFiles: []string{},
	}

	Parse(Configuration)

	modifiedFiles, err := GitStatus()
	if err != nil {
		return
	}

	// Remove ignored files
	if Configuration.IgnoredFiles != nil {
		modifiedFiles = RemoveIgnoredFiles(modifiedFiles, Configuration.IgnoredFiles)
	}

	commits := GetCommits(modifiedFiles)

	final := FormatFinalCommit(commits)
	if final == "" {
		return
	}

	if Configuration.Commit.CopyToClipboard {
		toClipboard([]byte(final))
	}
	fmt.Println(final)
}

func GetTitle(commitConfig *Commit) string {

	titlePrompt := "\n>"
	var titlePrefix string

	if commitConfig.Output != nil {
		output := commitConfig.Output
		if output.TitlePrefix != "" {
			titlePrefix = output.TitlePrefix
			titlePrompt = fmt.Sprintf("\n> %s ", titlePrefix)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	title := titlePrefix

	for {
		fmt.Print(titlePrompt)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("", text) == 0 {
			fmt.Println("Please input a title for your commit.")
		} else {
			title += " " + text
			break
		}
	}

	return title
}

func RemoveIgnoredFiles(fileList, ignoredFiles []string) []string {
	for i := 0; i < len(fileList); i++ {
		for _, file := range ignoredFiles {
			if strings.Contains(fileList[i], file) || strings.Contains(fileList[i], ignoredDefaultConfigFile) {
				fileList = append(fileList[:i], fileList[i+1:]...)
			}
		}
	}
	return fileList
}
