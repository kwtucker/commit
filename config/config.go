package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

type Flags struct {
	DryRun          bool
	CopyToClipboard bool
	RemoveText      bool
	Title           string
}

type Config struct {
	Commit *Commit `json:"commit,omitempty"`
	Clean  *Clean  `json:"clean,omitempty"`
	DryRun bool    `json:"-"`
}

type Clean struct {
	Identifier string `json:"identifier"`
}

type Commit struct {
	CopyToClipboard bool    `json:"copyToClipboard"`
	RemoveText      bool    `json:"removeText"`
	Output          *Output `json:"output,omitempty"`
}

type Output struct {
	Prefix      string `json:"prefix"`
	TitlePrefix string `json:"titlePrefix"`
}

func LoadConfig(flags Flags) *Config {
	_ = godotenv.Load(".envrc")
	cfg := &Config{
		DryRun: flags.DryRun,
		Commit: &Commit{
			CopyToClipboard: flags.CopyToClipboard,
			RemoveText:      flags.RemoveText,
			Output: &Output{
				Prefix:      "*",
				TitlePrefix: os.Getenv("COMMIT_TITLE_PREFIX"),
			},
		},
	}

	cfg.FillEnvs(".")

	return cfg
}

func (c *Config) FillEnvs(dir string) {
	_ = godotenv.Load(fmt.Sprintf("%s/.envrc", dir))

	prefix := os.Getenv("COMMIT_PREFIX")
	if prefix != "" {
		c.Commit.Output.Prefix = prefix
	}
}
