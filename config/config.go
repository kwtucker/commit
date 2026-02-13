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
}

type Config struct {
	CopyToClipboard bool    `json:"copyToClipboard"`
	Commit          *Commit `json:"commit,omitempty"`
}

type Commit struct {
	Format *Format `json:"output,omitempty"`
}

type Format struct {
	BodyPrefix string `json:"bodyPrefix"`
}

func LoadConfig(flags Flags) *Config {
	_ = godotenv.Load(".envrc")
	cfg := &Config{
		CopyToClipboard: flags.CopyToClipboard,
		Commit: &Commit{
			Format: &Format{
				BodyPrefix: "*",
			},
		},
	}

	cfg.FillEnvs(".")

	return cfg
}

func (c *Config) FillEnvs(dir string) {
	_ = godotenv.Load(fmt.Sprintf("%s/.envrc", dir))

	prefix := os.Getenv("COMMIT_BODY_PREFIX")
	if prefix != "" {
		c.Commit.Format.BodyPrefix = prefix
	}
}
