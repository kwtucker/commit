package main

import (
	"regexp"
	"strings"
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
