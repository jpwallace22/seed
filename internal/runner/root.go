package runner

import (
	"fmt"

	"github.com/jpwallace22/seed/internal/ctx"
	"github.com/jpwallace22/seed/internal/parser"
	"github.com/spf13/cobra"
	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

type RootFlags struct {
	FilePath      string
	FromClipboard bool
}

type RootRunner struct {
	clipboard clipboard.Clipboard
	parser    parser.Parser
	ctx       ctx.SeedContext
}

func NewRootRunner(cobra *cobra.Command, silent bool) Runner[RootFlags] {
	ctx := ctx.Build(cobra, silent)
	return &RootRunner{
		ctx:       *ctx,
		clipboard: clipboard.New(),
		parser:    parser.New(*ctx),
	}
}

func (r *RootRunner) Run(flags RootFlags, args []string) error {
	logger := r.ctx.Logger
	success := false

	switch {

	case flags.FromClipboard:
		err := r.parseFromClipboard()
		if err != nil {
			return fmt.Errorf("unable to parse from clipboard: %w", err)
		}
		success = true

	case flags.FilePath != "":
		logger.Log("Sowing the seeds of " + flags.FilePath)
		success = true

	case len(args) > 0:
		logger.Log("Sprouting directories from seed: %s", args[0])

		if err := r.parser.ParseTreeString(args[0]); err != nil {
			return fmt.Errorf("unable to parse the tree structure: %w", err)
		}
	}

	if success {
		logger.Success("Your directory tree has grown successfully!")
	} else {
		r.ctx.Cobra.Help()
	}
	return nil
}

func (r *RootRunner) parseFromClipboard() error {
	text, err := r.clipboard.PasteText()
	if err != nil {
		return fmt.Errorf("clipboard read error: %w", err)
	}

	r.ctx.Logger.Log("Planting from clipboard...")

	if err := r.parser.ParseTreeString(text); err != nil {
		return fmt.Errorf("unable to parse the tree structure: %w", err)
	}
	return nil
}
