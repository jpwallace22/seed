package parser

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jpwallace22/seed/internal/ctx"
)

type treeNodePool struct {
	node sync.Pool
}

func NewTreeNodePool() *treeNodePool {
	return &treeNodePool{
		node: sync.Pool{
			New: func() interface{} {
				return &TreeNode{
					children: make([]*TreeNode, 0, 8), // Preallocate with common size
				}
			},
		},
	}
}

type stringParser struct {
	ctx  *ctx.SeedContext
	pool *treeNodePool
}

func NewTreeParser(ctx *ctx.SeedContext) Parser {
	return &stringParser{
		ctx:  ctx,
		pool: NewTreeNodePool(),
	}
}

// converts a text representation of a directory tree into actual directories and files
func (p *stringParser) ParseTree(tree string) error {
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

	err = createFileSystem(root, "", p.ctx.Logger)
	if err != nil {
		return err
	}

	return nil
}

// converts the string lines into a tree structure
func (p *stringParser) buildTree(lines []string) (*TreeNode, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("no lines to parse")
	}

	// Preallocate lastNodes map with expected capacity
	lastNodes := make(map[int]*TreeNode, len(lines)/4)

	rootName := strings.TrimSpace(lines[0])
	if rootName == "" {
		return nil, fmt.Errorf("a root is required")
	}

	replacer := strings.NewReplacer(
		"├── ", "", "└── ", "", "│   ", "", "/", "", "\\", "", "    ", "",
		"│", "", "└", "", "├", "", "─", "",
	)

	root := p.pool.node.Get().(*TreeNode)
	root.name = rootName
	root.isFile = strings.Contains(rootName, ".")
	root.depth = 0
	root.children = root.children[:0] // Reset slice but keep capacity

	lastNodes[0] = root

	for i := 1; i < len(lines); i++ {
		// Need to normalize the line by changing all spaces with ASCII
		line := strings.ReplaceAll(lines[i], "\u00a0", " ")

		if strings.TrimSpace(line) == "" {
			continue
		}

		// Build the node
		depth := p.getDepth(line)
		name := strings.TrimSpace(replacer.Replace(strings.TrimSpace(line)))
		if name == "" {
			continue
		}

		node := p.pool.node.Get().(*TreeNode)
		node.name = name
		node.isFile = strings.Contains(name, ".")
		node.depth = depth
		node.children = node.children[:0]

		// Assign the node to a parent
		parentDepth := depth - 1
		parent := lastNodes[parentDepth]
		if parent == nil {
			p.cleanupTree(root)
			return nil, fmt.Errorf("invalid tree structure: missing parent at depth %d for node %s", parentDepth, name)
		}

		parent.children = append(parent.children, node)
		lastNodes[depth] = node
	}

	return root, nil
}

func (p *stringParser) getDepth(line string) int {
	depth := 0
	i := 0
	for i < len(line) {
		if strings.HasPrefix(line[i:], "│   ") ||
			strings.HasPrefix(line[i:], "    ") {
			depth++
			i += 4
		} else if strings.HasPrefix(line[i:], "├── ") ||
			strings.HasPrefix(line[i:], "└── ") {
			depth++
			i += 4
			break // These always appear at the end of the depth indicators and can break early
		} else {
			i++
		}
	}

	return depth
}

func (p *stringParser) cleanupTree(node *TreeNode) {
	if node == nil {
		return
	}
	for _, child := range node.children {
		p.cleanupTree(child)
	}
	node.children = nil
	p.pool.node.Put(node)
}
