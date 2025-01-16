package runner

import (
	"fmt"

	"github.com/jpwallace22/seed/internal/ctx"
	"github.com/jpwallace22/seed/internal/parser"
	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

type RootFlags struct {
	FromClipboard bool
}

type RootRunner struct {
	clipboard clipboard.Clipboard
	parser    parser.Parser
	ctx       ctx.SeedContext
}

func NewRootRunner(config Config) Runner[RootFlags] {
	ctx := ctx.Build(config.Silent)
	return &RootRunner{
		ctx:       *ctx,
		clipboard: clipboard.New(),
		parser:    parser.New(*ctx),
	}
}

func (r *RootRunner) Run(flags RootFlags, args []string) error {
	logger := r.ctx.Logger
	if flags.FromClipboard {
		logger.Log("Planting from clipboard...")

		text, err := r.getClipboardContent()
		if err != nil {
			return fmt.Errorf("unable to access clipboard contents: %w", err)
		}

		if err := r.parser.ParseTreeString(text); err != nil {
			return fmt.Errorf("unable to parse the tree structure: %w", err)
		}
	} else if len(args) > 0 {
		logger.Log("Sprouting directories from seed: %s", args[0])

		if err := r.parser.ParseTreeString(args[0]); err != nil {
			return fmt.Errorf("unable to parse the tree structure: %w", err)
		}
	} else {
		return fmt.Errorf("no seeds provided: provide a path or use -c to source from your clipboard")
	}

	logger.Success("Your directory tree has grown successfully!")
	return nil
}

func (r *RootRunner) getClipboardContent() (string, error) {
	text, err := r.clipboard.PasteText()
	if err != nil {
		return "", fmt.Errorf("clipboard read error: %w", err)
	}
	return text, nil
}
