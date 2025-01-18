package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jpwallace22/seed/cmd/flags"
	"github.com/jpwallace22/seed/internal/ctx"
	"github.com/jpwallace22/seed/pkg/logger"
)

type Parser interface {
	ParseTree(string) error
}

type TreeNode struct {
	name     string
	children []*TreeNode
	isFile   bool
	depth    int
}

type Option func(*config)

type config struct {
	format flags.Format
}

func NewParser(ctx *ctx.SeedContext, opts ...Option) (Parser, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	cfg := &config{
		format: flags.Formats.Tree,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	switch cfg.format {
	case flags.Formats.JSON:
		return NewJSONParser(ctx), nil
	case flags.Formats.Tree:
		return NewTreeParser(ctx), nil
	default:
		return nil, fmt.Errorf("unsupported parser format: %s", cfg.format)
	}
}

func WithFormat(format flags.Format) Option {
	return func(c *config) {
		c.format = format
	}
}

// This should move to a new module for filesystems, its breaking single job pattern
func createFileSystem(node *TreeNode, parentPath string, logger logger.Logger) error {
	if node == nil {
		return nil
	}
	permissions := os.FileMode(0755)

	currentPath := parentPath
	if node.name != "." {
		currentPath = filepath.Join(parentPath, node.name)
	}

	// create current node unless it's the "." root
	if node.name != "." {
		if node.isFile {
			// ensure parent directory exists
			parentDir := filepath.Dir(currentPath)
			if err := os.MkdirAll(parentDir, permissions); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", parentDir, err)
			}

			// create file
			f, err := os.Create(currentPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", currentPath, err)
			}

			f.Close()
			logger.Info("Planted file: " + currentPath)
		} else {
			if err := os.MkdirAll(currentPath, permissions); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", currentPath, err)
			}
			logger.Info("Planted directory: " + currentPath)
		}
	}

	// loop through children with the correct parent path
	for _, child := range node.children {
		if err := createFileSystem(child, currentPath, logger); err != nil {
			return err
		}
	}

	return nil
}
