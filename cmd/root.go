package cmd

import (
	"fmt"

	"github.com/kwtucker/commit/config"
	"github.com/kwtucker/commit/git"
	"github.com/kwtucker/commit/output"
	"github.com/spf13/cobra"
)

var (
	DryRun          bool
	CopyToClipboard bool
	RemoveText      bool
	Title           string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "d", false, "dry run to inspect the result")
	rootCmd.PersistentFlags().BoolVarP(&CopyToClipboard, "copy", "c", false, "copy commit message to clipboard")
	rootCmd.PersistentFlags().BoolVarP(&RemoveText, "rm-text", "r", false, "remove text within the delimiters from the file after reading message")
	rootCmd.PersistentFlags().StringVarP(&Title, "title-msg", "t", "", "quoted title of the commit message")
}

var rootCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit will help construct a commit message.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig(config.Flags{
			DryRun:          DryRun,
			CopyToClipboard: CopyToClipboard,
			RemoveText:      RemoveText,
			Title:           Title,
		})

		status, err := git.GitStatus()
		if err != nil || len(status) == 0 {
			return
		}

		var title string
		if Title != "" {
			title = fmt.Sprintf("%s", Title)
			if cfg.Commit.Output != nil {
				if titlePrefix := cfg.Commit.Output.TitlePrefix; titlePrefix != "" {
					title = fmt.Sprintf("%s %s", titlePrefix, Title)
				}
			}
		}

		commits := git.GetCommits(cfg, status)

		final := output.FormatFinalCommit(title, commits)
		if final == "" {
			return
		}

		output.ToClipboard(cfg, []byte(final))

		fmt.Println(final)
	},
}
