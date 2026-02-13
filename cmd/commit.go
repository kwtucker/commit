package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kwtucker/commit/config"
	"github.com/kwtucker/commit/git"
	"github.com/kwtucker/commit/internal/teaui"
	"github.com/kwtucker/commit/output"
	"github.com/mattn/go-isatty"
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
			// return
		}

		if cmd != cmd.Root() {
			fmt.Println("--interactive can only be used on the root command")
			os.Exit(1)
		}

		if !isatty.IsTerminal(os.Stdout.Fd()) {
			fmt.Println("interactive mode requires a TTY")
			os.Exit(1)
		}

		m := teaui.New(cfg.Commit.Format.BodyPrefix)
		p := tea.NewProgram(m)

		final, err := p.Run()
		if err != nil {
			fmt.Println("interactive mode failed:", err)
			os.Exit(1)
		}

		// Cast to Model to access Result()
		commit := final.(teaui.Model).Result()
		fmt.Println(commit)

		if err := git.Commit(commit); err != nil {
			fmt.Println("failed to commit:", err)
			os.Exit(1)
		}

		if cfg.CopyToClipboard {
			output.ToClipboard(cfg, []byte(commit))
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&CopyToClipboard, "copy", "c", false, "copy commit message to clipboard")
}
