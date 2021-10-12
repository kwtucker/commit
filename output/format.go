package output

import (
	"regexp"
	"strings"
)

func FormatFinalCommit(out []string) string {
	for i, commit := range out {
		space := regexp.MustCompile(`\s+`)
		str := space.ReplaceAllString(commit, " ")
		out[i] = str
	}

	return strings.Join(out, "\n\n")
}
