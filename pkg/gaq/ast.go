package gaq

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/tamayika/gaq/pkg/gaq/query"
)

// Node represents traversible ast.Node
type Node struct {
	Type     string  `json:"type"`
	Pos      int     `json:"pos"`
	End      int     `json:"end"`
	Children []*Node `json:"children,omitempty"`

	Node ast.Node `json:"-"`
	Name string   `json:"-"`
}

// Parse parses source and returns ast
func Parse(source string) (*Node, error) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	w := &walker{fset: fset}
	ast.Walk(w, f)
	if w.err != nil {
		return nil, err
	}
	return w.node, nil
}

// MustParse parses source and returns ast
// If failed to parse, fatal occurs
func MustParse(source string) *Node {
	n, err := Parse(source)
	if err != nil {
		log.Fatalf("Cannot parse source. %v", err)
	}
	return n
}

type walker struct {
	node *Node
	fset *token.FileSet
	err  error
}

func (w *walker) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	child := buildNode(node)
	if w.node == nil {
		// for file
		w.node = child
	} else {
		w.node.Children = append(w.node.Children, child)
	}
	return &walker{node: child}
}

func buildNode(n ast.Node) *Node {
	nodeType := nodeType(n)
	splits := strings.Split(nodeType, ".")
	node := &Node{
		Type:     nodeType,
		Pos:      int(n.Pos()),
		End:      int(n.End()),
		Children: []*Node{},
		Node:     n,
		Name:     splits[len(splits)-1],
	}
	return node
}

func nodeType(n interface{}) string {
	var bf bytes.Buffer
	fmt.Fprintf(&bf, "%T", n)
	nt := string(bf.Bytes())
	nt = strings.Replace(nt, "*", "", -1)
	nt = strings.Replace(nt, "[]", "", -1)
	return nt
}

// QuerySelector queries to node and return first matched node
func (n *Node) QuerySelector(q *query.Query) ast.Node {
	var firstNode ast.Node
	n.apply(q, 0, 0, -1, func(n *Node) bool {
		firstNode = n.Node
		return false
	})
	return firstNode
}

// QuerySelectorAll queries to node and return all matched nodes
func (n *Node) QuerySelectorAll(q *query.Query) []ast.Node {
	nodes := []ast.Node{}
	n.apply(q, 0, 0, -1, func(n *Node) bool {
		nodes = append(nodes, n.Node)
		return true
	})
	return nodes
}

type callback func(n *Node) bool

func (n *Node) apply(q *query.Query, entryIndex int, nodeDepth int, lastMatchedNodeDepth int, cb callback) bool {
	if entryIndex >= len(q.Entries) {
		return true
	}
	entry := q.Entries[entryIndex]
	mustBeChild := entry.Combinator == ">"
	if mustBeChild && lastMatchedNodeDepth >= 0 && nodeDepth-lastMatchedNodeDepth > 1 {
		return true
	}
	if entry.Name == n.Name {
		if entryIndex+1 == len(q.Entries) {
			continues := cb(n)
			if !continues {
				return false
			}
		}
		return n.applyChildren(q, entryIndex+1, nodeDepth+1, nodeDepth, cb)
	}
	return n.applyChildren(q, entryIndex, nodeDepth+1, lastMatchedNodeDepth, cb)
}

func (n *Node) applyChildren(q *query.Query, entryIndex, nodeDepth int, lastMatchedNodeDepth int, cb callback) bool {
	for _, childNode := range n.Children {
		continues := childNode.apply(q, entryIndex, nodeDepth, lastMatchedNodeDepth, cb)
		if !continues {
			return false
		}
	}
	return true
}
