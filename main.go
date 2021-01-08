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
	noDetectedCommit         = "No Commit Messages Detected"
)

var Configuration *Config

func main() {

	Configuration = &Config{}
	err := Parse(Configuration)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	modifiedFiles, err := GitStatus()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(modifiedFiles) == 0 {
		fmt.Println("**** No files are staged yet.")
		return
	}

	// Remove ignored files
	if Configuration.IgnoredFiles != nil {
		modifiedFiles = RemoveIgnoredFiles(modifiedFiles, Configuration.IgnoredFiles)
	}

	commit := Configuration.Commit
	if commit != nil {

		title := GetTitle(commit)
		commits := GetCommits(modifiedFiles)

		final := FormatFinalCommit(title, commits)
		if final == "" {
			fmt.Println(noDetectedCommit)
		}

		if Configuration.Commit.CopyToClipboard {
			toClipboard([]byte("\"" + final + "\""))
		} else {
			fmt.Println(final)
		}
	}

	clean := Configuration.Clean
	if clean != nil {
		err := CleanFiles()
		if err != nil {
			fmt.Println("\nClean Error:")
			fmt.Println(err.Error())
		}
	}

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

func CleanFiles() error {
	return nil
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
