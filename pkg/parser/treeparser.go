package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jpwallace22/seed/pkg/ctx"
)

type Parser struct {
	ctx ctx.SeedContext
}

type TreeNode struct {
	name     string
	children []*TreeNode
	isFile   bool
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

	if err := p.createFileSystem(root, ""); err != nil {
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
		return nil, fmt.Errorf("invalid root name")
	}

	root := &TreeNode{
		name:     rootName,
		isFile:   strings.Contains(rootName, ".") && rootName != ".",
		children: make([]*TreeNode, 0),
	}

	// Stack to keep track of the current path in the tree
	stack := []*TreeNode{root}
	previousDepth := 0

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		depth := p.getDepth(lines[i])
		name := p.extractName(lines[i])
		if name == "" {
			continue
		}

		node := &TreeNode{
			name:     name,
			isFile:   strings.Contains(name, "."),
			children: make([]*TreeNode, 0),
		}

		// Adjust stack based on depth
		if depth > previousDepth {
			// Going deeper - keep current path
		} else if depth < previousDepth {
			// Going back up - pop from stack until we're at the right level
			stack = stack[:depth+1]
		} else if depth == previousDepth && len(stack) > depth+1 {
			// Same level - ensure stack has correct length
			stack = stack[:depth+1]
		}

		// Add node to its parent
		parent := stack[depth]
		parent.children = append(parent.children, node)

		// If this is a directory, add it to the stack for potential children
		if !node.isFile {
			if depth+1 < len(stack) {
				stack[depth+1] = node
			} else {
				stack = append(stack, node)
			}
		}

		previousDepth = depth
	}

	return root, nil
}

// calculates the depth of a line based on tree characters
func (p *Parser) getDepth(line string) int {
	depth := 0
	for i := 0; i < len(line); {
		if strings.HasPrefix(line[i:], "│   ") || strings.HasPrefix(line[i:], "    ") {
			depth++
			i += 4
		} else {
			break
		}
	}

	return depth
}

// gets the actual name from a tree line by removing tree characters
func (p *Parser) extractName(line string) string {
	line = strings.TrimSpace(line)
	treeChars := []string{
		"├── ", "└── ",
		"│   ", "    ",
		"│", "└", "├", "─",
	}

	for _, char := range treeChars {
		line = strings.ReplaceAll(line, char, "")
	}

	return strings.TrimSpace(line)
}

// recursivly creates the actual directory structure from the parsed tree
func (p *Parser) createFileSystem(node *TreeNode, parentPath string) error {
	if node == nil {
		return nil
	}
	logger := p.ctx.Logger
	permissions := os.FileMode(0755)

	// handle root node specially
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

	// loop it
	for _, child := range node.children {
		if err := p.createFileSystem(child, currentPath); err != nil {
			return err
		}
	}

	return nil
}
