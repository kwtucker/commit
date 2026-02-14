package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kwtucker/commit/config"
	"github.com/kwtucker/commit/git"
	"github.com/kwtucker/commit/internal/teaui"
	"github.com/kwtucker/commit/output"
	"github.com/spf13/cobra"
)

var CopyToClipboard bool

var RootCmd = &cobra.Command{
	Use:   "commit",
	Short: "Constructs formatted commit messages",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig(config.Flags{
			CopyToClipboard: CopyToClipboard,
		})

		status, err := git.GetStagedFiles()
		if err != nil || len(status) == 0 {
			fmt.Println("no staged changes", err)
			os.Exit(1)
			return
		}

		m := teaui.New(cfg.Commit.Format.BodyPrefix)
		p := tea.NewProgram(m)

		final, err := p.Run()
		if err != nil {
			fmt.Println("commit failed:", err)
			os.Exit(1)
		}

		// Cast to Model to access Result()
		commit := final.(teaui.Model).Result()
		fmt.Println(commit)

		if cfg.CopyToClipboard {
			output.ToClipboard(cfg, []byte(commit))
			return
		}

		if err := git.Commit(commit); err != nil {
			fmt.Println("failed to commit:", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&CopyToClipboard, "copy", "c", false, "copy commit message to clipboard")
}
