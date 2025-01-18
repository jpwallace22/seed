package runner

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jpwallace22/seed/cmd/flags"
	"github.com/jpwallace22/seed/internal/ctx"
	"github.com/jpwallace22/seed/internal/parser"
	"github.com/spf13/cobra"
	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

const (
	msgSuccess = "Your directory tree has grown successfully!"
)

type RootRunner struct {
	clipboard clipboard.Clipboard
	parser    parser.Parser
	ctx       ctx.SeedContext
}

func NewRootRunner(cobra *cobra.Command, ctx *ctx.SeedContext) Runner[flags.RootFlags] {
	return &RootRunner{
		ctx:       *ctx,
		clipboard: clipboard.New(),
		parser:    parser.NewTreeParser(*ctx),
	}
}

func (r *RootRunner) Run(args []string) error {
	logger := r.ctx.Logger
	flags := r.ctx.Flags.Root

	switch {
	case flags.FromClipboard:
		if err := r.parseFromClipboard(); err != nil {
			return fmt.Errorf("unable to parse from clipboard: %w", err)
		}
		logger.Success(msgSuccess)
		return nil

	case flags.FilePath != "":
		if err := r.parseFromFile(flags.FilePath); err != nil {
			return fmt.Errorf("unable to parse from file: %w", err)
		}
		logger.Success(msgSuccess)
		return nil

	case len(args) > 0:
		logger.Log("Sprouting directories from seed: %s", args[0])
		if err := r.parser.ParseTree(args[0]); err != nil {
			return fmt.Errorf("unable to parse the tree structure: %w", err)
		}
		logger.Success(msgSuccess)
		return nil
	}

	return r.ctx.Cobra.Help()
}

func (r *RootRunner) parseFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("file read error: %w", err)
	}

	r.ctx.Logger.Log("Sowing the seeds of " + filepath.Base(path) + "...")
	if err := r.parser.ParseTree(string(data)); err != nil {
		return fmt.Errorf("unable to parse the tree structure: %w", err)
	}
	return nil
}

func (r *RootRunner) parseFromClipboard() error {
	text, err := r.clipboard.PasteText()
	if err != nil {
		return fmt.Errorf("clipboard read error: %w", err)
	}

	r.ctx.Logger.Log("Planting from clipboard...")

	if err := r.parser.ParseTree(text); err != nil {
		return fmt.Errorf("unable to parse the tree structure: %w", err)
	}
	return nil
}
