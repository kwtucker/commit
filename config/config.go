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
}

type Config struct {
	Commit       *Commit  `json:"commit,omitempty"`
	Clean        *Clean   `json:"clean,omitempty"`
	DryRun       bool     `json:"-"`
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
	cfg := &Config{
		DryRun: flags.DryRun,
		Commit: &Commit{
			CopyToClipboard: flags.CopyToClipboard,
			RemoveText:      flags.RemoveText,
			Output: &Output{
				Prefix: "*",
			},
		},
	}
	cfg.FillEnvs(".")

	return cfg
}

func (c *Config) FillEnvs(dir string) {
	_ = godotenv.Overload(fmt.Sprintf("%s/.env", dir))
	prefix := os.Getenv("COMMIT_PREFIX")
	if prefix != "" {
		c.Commit.Output.Prefix = prefix
	}
}
