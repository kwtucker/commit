package git

import (
	"os"

	"github.com/kwtucker/commit/config"
	"github.com/kwtucker/commit/parser"
)

func GetCommits(cfg *config.Config, modifiedFiles []string) []string {
	out := []string{}

	// If a ".commit" file exists in the directory commit is executed
	// it will be parsed.
	_, err := os.Stat(".commit")
	if !os.IsNotExist(err) {
		out = append(out, parser.ReadFile(cfg, ".commit")...)
	}

	// Parse Files
	for _, filename := range modifiedFiles {
		UnStageFile(filename)
		out = append(out, parser.ReadFile(cfg, filename)...)
		StageFile(filename)
	}

	return out
}
