package runner

import (
	"fmt"
	"os"

	"github.com/jpwallace22/seed/internal/ctx"
	"github.com/jpwallace22/seed/internal/parser"
	"github.com/spf13/cobra"
	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

const (
	msgSuccess = "Your directory tree has grown successfully!"
)

type Format string

const (
	Tree Format = "tree"
	JSON Format = "json"
	YAML Format = "yaml"
)

func (f Format) String() string {
	return string(f)
}

func (f *Format) Set(value string) error {
	switch Format(value) {
	case Tree, JSON, YAML:
		*f = Format(value)
		return nil
	default:
		return fmt.Errorf("invalid format %q, must be one of: tree, json, yaml", value)
	}
}

func (f Format) Type() string {
	return "format"
}

type RootFlags struct {
	FilePath      string
	Format        Format
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

	switch {
	case flags.FromClipboard:
		if err := r.parseFromClipboard(); err != nil {
			return fmt.Errorf("unable to parse from clipboard: %w", err)
		}
		logger.Success(msgSuccess)
		return nil

	case flags.FilePath != "":
		logger.Log("Sowing the seeds of " + flags.FilePath)
		if err := r.parseFromFile(flags.FilePath); err != nil {
			return fmt.Errorf("unable to parse from file: %w", err)
		}
		logger.Success(msgSuccess)
		return nil

	case len(args) > 0:
		logger.Log("Sprouting directories from seed: %s", args[0])
		if err := r.parser.ParseTreeString(args[0]); err != nil {
			return fmt.Errorf("unable to parse the tree structure: %w", err)
		}
		logger.Success(msgSuccess)
		return nil
	}

	return r.ctx.Cobra.Help()
}

func (r *RootRunner) parseFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("file read error: %w", err)
	}

	r.ctx.Logger.Log("Planting from file...")
	if err := r.parser.ParseTreeString(string(data)); err != nil {
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

	if err := r.parser.ParseTreeString(text); err != nil {
		return fmt.Errorf("unable to parse the tree structure: %w", err)
	}
	return nil
}
