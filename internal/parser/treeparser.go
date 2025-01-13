package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jpwallace22/seed/internal/ctx"
)

type Parser struct {
	ctx ctx.SeedContext
}

type TreeNode struct {
	name     string
	children []*TreeNode
	isFile   bool
	depth    int
}

func New(ctx ctx.SeedContext) *Parser {
	return &Parser{
		ctx: ctx,
	}
}

// converts a text representation of a directory tree into actual directories and files
func (p *Parser) ParseTreeString(tree string) error {
	lines := strings.Split(strings.TrimSpace(tree), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no tree provided")
	}

	if strings.TrimSpace(lines[0]) == "tree" {
		lines = lines[1:]
	}

	root, err := p.buildTree(lines)
	if err != nil {
		return fmt.Errorf("failed to parse tree: %w", err)
	}

	err = p.createFileSystem(root, "")
	if err != nil {
		return err
	}

	return nil
}

// converts the string lines into a tree structure
func (p *Parser) buildTree(lines []string) (*TreeNode, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("no lines to parse")
	}

	rootName := strings.TrimSpace(lines[0])
	if rootName == "" {
		return nil, fmt.Errorf("a root is required")
	}

	root := &TreeNode{
		name:     rootName,
		isFile:   strings.Contains(rootName, ".") && rootName != ".",
		children: make([]*TreeNode, 0),
		depth:    0,
	}

	// Keep track of last nodes at each depth level
	lastNodes := make(map[int]*TreeNode)
	lastNodes[0] = root

	for i := 1; i < len(lines); i++ {
		// Need to normalize the line by changing all spaces with ASCII
		line := strings.ReplaceAll(lines[i], "\u00a0", " ")

		if strings.TrimSpace(line) == "" {
			continue
		}

		// Build the node
		depth := p.getDepth(line)
		name := p.extractName(line)
		if name == "" {
			continue
		}

		node := &TreeNode{
			name:     name,
			isFile:   strings.Contains(name, "."),
			children: make([]*TreeNode, 0),
			depth:    depth,
		}

		// Assign the node to a parent
		parentDepth := depth - 1
		parent := lastNodes[parentDepth]
		if parent == nil {
			return nil, fmt.Errorf("invalid tree structure: missing parent at depth %d for node %s", parentDepth, name)
		}

		parent.children = append(parent.children, node)
		lastNodes[depth] = node
	}

	return root, nil
}

func (p *Parser) getDepth(line string) int {
	depth := 0
	for i := 0; i < len(line); {
		if strings.HasPrefix(line[i:], "│   ") {
			depth++
			i += 4
		} else if strings.HasPrefix(line[i:], "    ") {
			depth++
			i += 4
		} else if strings.HasPrefix(line[i:], "├── ") {
			depth++
			i += 4
		} else if strings.HasPrefix(line[i:], "└── ") {
			depth++
			i += 4
		} else {
			i++
		}
	}
	return depth
}

func (p *Parser) extractName(line string) string {
	line = strings.TrimSpace(line)
	unwantedChars := []string{
		"├── ", "└── ", "/", "\\",
		"│   ", "    ",
		"│", "└", "├", "─",
	}
	for _, char := range unwantedChars {
		line = strings.ReplaceAll(line, char, "")
	}
	return strings.TrimSpace(line)
}

func (p *Parser) createFileSystem(node *TreeNode, parentPath string) error {
	if node == nil {
		return nil
	}
	logger := p.ctx.Logger
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
		if err := p.createFileSystem(child, currentPath); err != nil {
			return err
		}
	}

	return nil
}
