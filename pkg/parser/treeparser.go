package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jpwallace22/seed/pkg/logger"
)

// represents a single node in the directory tree
type TreeNode struct {
	name     string
	children []*TreeNode
	isFile   bool
}

// converts a text representation of a directory tree into actual directories and files
func ParseTreeString(tree string) error {
	lines := strings.Split(strings.TrimSpace(tree), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no tree provided")
	}

	if strings.TrimSpace(lines[0]) == "tree" {
		lines = lines[1:]
	}

	root, err := buildTree(lines)
	if err != nil {
		return fmt.Errorf("failed to parse tree: %w", err)
	}

	if err := createFileSystem(root, ""); err != nil {
		return err
	}

	return nil
}

// converts the string lines into a tree structure
func buildTree(lines []string) (*TreeNode, error) {
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

	currentLevel := []*TreeNode{root}
	previousDepth := 0

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		depth := getDepth(lines[i])
		name := extractName(lines[i])
		if name == "" {
			continue
		}

		node := &TreeNode{
			name:     name,
			isFile:   strings.Contains(name, "."),
			children: make([]*TreeNode, 0),
		}

		if depth > previousDepth {
			currentLevel = append(currentLevel, currentLevel[len(currentLevel)-1])
		} else if depth < previousDepth {
			currentLevel = currentLevel[:depth+1]
		}

		parent := currentLevel[len(currentLevel)-1]
		parent.children = append(parent.children, node)

		if !node.isFile {
			currentLevel[len(currentLevel)-1] = node
		}
		previousDepth = depth
	}

	return root, nil
}

// calculates the depth of a line based on tree characters
func getDepth(line string) int {
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
func extractName(line string) string {
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
func createFileSystem(node *TreeNode, parentPath string) error {
	if node == nil {
		return nil
	}
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
		if err := createFileSystem(child, currentPath); err != nil {
			return err
		}
	}

	return nil
}
