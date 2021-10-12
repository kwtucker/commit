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
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "d", false, "Dry run to inspect the result.")
	rootCmd.PersistentFlags().BoolVarP(&CopyToClipboard, "copy", "c", true, "Copy commit message to clipboard.")
	rootCmd.PersistentFlags().BoolVarP(&RemoveText, "rm-text", "r", false, "Remove text from the file after reading message.")
}

var rootCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit will help construct a commit message.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig(config.Flags{
			DryRun:          DryRun,
			CopyToClipboard: CopyToClipboard,
			RemoveText:      RemoveText,
		})

		status, err := git.GitStatus()
		if err != nil {
			return
		}

		commits := git.GetCommits(cfg, status)

		final := output.FormatFinalCommit(commits)
		if final == "" {
			return
		}

		if !cfg.DryRun {
			output.ToClipboard(cfg, []byte(final))
		}

		fmt.Println(final)
	},
}
