package cmd

import (
	"fmt"
	"os"

	"github.com/jpwallace22/seed/pkg/ctx"
	"github.com/jpwallace22/seed/pkg/parser"
	"github.com/spf13/cobra"
	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

type Runner struct {
	ctx       ctx.SeedContext
	clipboard clipboard.Clipboard
}

func NewRunner(silent bool) *Runner {
	return &Runner{
		ctx:       *ctx.Build(silent),
		clipboard: clipboard.New(),
	}
}

var (
	// filePath      string //TODO: implement this feature
	fromClipboard bool
	silent        bool
)

func init() {
	rootCmd.Flags().BoolVarP(&fromClipboard, "clipboard", "c", false, "Use tree structure from clipboard.")
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "If true, suppresses all non-essential console output.")
}

var rootCmd = &cobra.Command{
	Use:   "seed [path/to/tree]",
	Short: "Plant the seeds of your directory tree 🌱.",
	Long:  "Seed is a CLI tool that helps you grow directory structures from a tree representation provided via file or clipboard.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCommand,
}

func runCommand(cmd *cobra.Command, args []string) error {
	runner := NewRunner(silent)
	return runner.Run(fromClipboard, args)
}

func (r *Runner) Run(fromClipboard bool, args []string) error {
	parser := parser.New(r.ctx)
	logger := r.ctx.Logger
	if fromClipboard {
		logger.Log("Planting from clipboard...")
		text, err := r.getClipboardContent()
		if err != nil {
			return fmt.Errorf("unable to access clipboard contents: %w", err)
		}

		if err := parser.ParseTreeString(text); err != nil {
			return fmt.Errorf("unable to parse the tree structure: %w", err)
		}
	} else if len(args) > 0 {
		logger.Log("Sprouting directories from file: %s", args[0])
		// Logic to read and process the file
	} else {
		return fmt.Errorf("no seeds provided: provide a path or use -c to source from your clipboard")
	}

	logger.Success("Your directory tree has grown successfully!")
	return nil
}

func (r *Runner) getClipboardContent() (string, error) {
	text, err := r.clipboard.PasteText()
	if err != nil {
		return "", fmt.Errorf("clipboard read error: %w", err)
	}
	return text, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
