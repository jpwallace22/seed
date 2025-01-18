package parser

import (
	"encoding/json"
	"fmt"

	"github.com/jpwallace22/seed/internal/ctx"
)

type FileNode struct {
	Type     string     `json:"type"`
	Name     string     `json:"name"`
	Contents []FileNode `json:"contents,omitempty"`
}

type Report struct {
	Type        string `json:"type"`
	Directories int    `json:"directories"`
	Files       int    `json:"files"`
}

type jsonParser struct {
	ctx *ctx.SeedContext
}

func NewJSONParser(ctx *ctx.SeedContext) Parser {
	return &jsonParser{ctx: ctx}
}

func (p *jsonParser) ParseTree(jsonStr string) error {
	if jsonStr == "" {
		return fmt.Errorf("no tree provided")
	}

	var nodes []json.RawMessage
	if err := json.Unmarshal([]byte(jsonStr), &nodes); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if len(nodes) == 0 {
		return fmt.Errorf("empty JSON array")
	}

	// First validate the root node structure
	if err := p.validateNode(nodes[0]); err != nil {
		return fmt.Errorf("failed to parse tree: %w", err)
	}

	// Then parse into FileNode
	var rootFileNode FileNode
	if err := json.Unmarshal(nodes[0], &rootFileNode); err != nil {
		return fmt.Errorf("failed to parse root node: %w", err)
	}

	// Convert to TreeNode
	rootTreeNode := p.fileNodeToTreeNode(&rootFileNode)

	// Create filesystem using the existing function
	if err := createFileSystem(rootTreeNode, "", p.ctx.Logger); err != nil {
		return fmt.Errorf("failed to create filesystem: %w", err)
	}

	// Verify report if present
	if len(nodes) > 1 {
		var report Report
		if err := json.Unmarshal(nodes[1], &report); err != nil {
			return fmt.Errorf("failed to parse report: %w", err)
		}

		dirs, files := p.countTreeNodes(rootTreeNode)
		if dirs != report.Directories || files != report.Files {
			return fmt.Errorf("file system count mismatch - expected: %d directories and %d files, got: %d directories and %d files",
				report.Directories, report.Files, dirs, files)
		}
	}

	return nil
}

func (p *jsonParser) validateNode(raw json.RawMessage) error {
	var node struct {
		Type     string            `json:"type"`
		Name     string            `json:"name"`
		Contents []json.RawMessage `json:"contents,omitempty"`
	}

	if err := json.Unmarshal(raw, &node); err != nil {
		return fmt.Errorf("invalid node: %w", err)
	}

	if node.Type == "" {
		return fmt.Errorf("missing type field")
	}

	if node.Name == "" {
		return fmt.Errorf("missing name field")
	}

	// Recursively validate contents
	for _, content := range node.Contents {
		if err := p.validateNode(content); err != nil {
			return err
		}
	}

	return nil
}

func (p *jsonParser) fileNodeToTreeNode(node *FileNode) *TreeNode {
	treeNode := &TreeNode{
		name:     node.Name,
		isFile:   node.Type == "file",
		children: make([]*TreeNode, 0),
	}

	for i := range node.Contents {
		childNode := p.fileNodeToTreeNode(&node.Contents[i])
		treeNode.children = append(treeNode.children, childNode)
	}

	return treeNode
}

func (p *jsonParser) countTreeNodes(node *TreeNode) (directories int, files int) {
	if node == nil {
		return 0, 0
	}

	if node.isFile {
		files = 1
	} else {
		directories = 1
	}

	for _, child := range node.children {
		d, f := p.countTreeNodes(child)
		directories += d
		files += f
	}

	return directories, files
}
