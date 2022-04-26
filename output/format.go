package output

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kwtucker/commit/config"
)

func FormatFinalCommit(title string, out []string) string {
	if title != "" {
		out = append([]string{title}, out...)
	}

	for i, commit := range out {
		space := regexp.MustCompile(`\s+`)
		str := space.ReplaceAllString(commit, " ")
		out[i] = str
	}

	return strings.Join(out, "\n\n")
}

func GetTitle(cfg *config.Config) string {
	titlePrompt := "\n>"
	var titlePrefix string
	if cfg.Commit.Output != nil {
		if titlePrefix = cfg.Commit.Output.TitlePrefix; titlePrefix != "" {
			titlePrompt = fmt.Sprintf("\n> %s ", titlePrefix)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	title := titlePrefix

	for {
		fmt.Print(titlePrompt)
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		if strings.Compare("", text) == 0 {
			fmt.Println("Please input a title for your commit.")
		} else {
			title += " " + text
			break
		}
	}

	return title
}
