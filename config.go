package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Commit       *Commit  `json:"commit,omitempty"`
	Clean        *Clean   `json:"clean,omitempty"`
	IgnoredFiles []string `json:"ignoredFiles"`
}

type Clean struct {
	Identifier string `json:"identifier"`
}

type Commit struct {
	TitlePrompt     bool    `json:"titlePrompt"`
	CopyToClipboard bool    `json:"copyToClipboard"`
	RemoveText      bool    `json:"removeText"`
	Output          *Output `json:"output,omitempty"`
}

type Output struct {
	Prefix      string `json:"prefix"`
	TitlePrefix string `json:"titlePrefix"`
}

// flag for config path
var config = flag.String("config", "forgit.json", "Path to json forgit file.")

// Parse reads config file and parses cli flags into c by calling flag.Parse()
func Parse(c interface{}) error {
	flag.Parse() // Parse first time to get config path.
	err := parseJSON(*config, c)
	if err != nil {
		return err
	}
	flag.Parse() // Call again to overwrite json values with flags.
	return nil
}

func parseJSON(configPath string, c interface{}) error {
	var err error
	if *config == "" {
		return err
	}
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("commit: *** Could not %+v, please create one and run command in that directory.", err.Error())
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		return fmt.Errorf("%+v", err.Error())
	}
	return err
}
