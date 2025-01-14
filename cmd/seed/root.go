package main

import (
	"fmt"
	"os"

	"github.com/jpwallace22/seed/internal/ctx"
	"github.com/jpwallace22/seed/internal/parser"
	"github.com/spf13/cobra"
	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

type Flags struct {
	FromClipboard bool
	Silent        bool
}

type Runner struct {
	clipboard clipboard.Clipboard
	ctx       ctx.SeedContext
}

func NewRunner(flags Flags) *Runner {
	return &Runner{
		ctx:       *ctx.Build(flags.Silent),
		clipboard: clipboard.New(),
	}
}

var rootCmd = &cobra.Command{
	Version: "0.1.0",
	Use:     "seed [path/to/tree]",
	Short:   "Plant the seeds of your directory tree 🌱.",
	Long:    "Seed is a CLI tool that helps you grow directory structures from a tree representation provided via string or clipboard.",
	Args:    cobra.MaximumNArgs(1),
	RunE:    runCommand,
}

var flags Flags

func init() {
	rootCmd.Flags().BoolVarP(&flags.FromClipboard, "clipboard", "c", false, "Use tree structure from clipboard.")
	rootCmd.Flags().BoolVarP(&flags.Silent, "silent", "s", false, "If true, suppresses all non-essential console output.")
}

func runCommand(cmd *cobra.Command, args []string) error {
	runner := NewRunner(flags)
	return runner.Run(flags.FromClipboard, args)
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
		logger.Log("Sprouting directories from seed: %s", args[0])

		if err := parser.ParseTreeString(args[0]); err != nil {
			return fmt.Errorf("unable to parse the tree structure: %w", err)
		}
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
