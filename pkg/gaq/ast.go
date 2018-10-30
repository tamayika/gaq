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
	Parent   *Node   `json:"-"`

	Node ast.Node `json:"-"`
	Name string   `json:"-"`
}

// Parse parses source and returns *Node
func Parse(source string) (*Node, error) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return ParseNode(f)
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

// ParseNode parses ast.Node and returns *Node
func ParseNode(n ast.Node) (*Node, error) {
	w := &walker{}
	ast.Walk(w, n)
	if w.err != nil {
		return nil, w.err
	}
	return w.node, nil
}

// MustParseNode parses ast.Node and returns *Node
func MustParseNode(n ast.Node) *Node {
	node, err := ParseNode(n)
	if err != nil {
		log.Fatalf("Cannot parse node. %v", err)
	}
	return node
}

type walker struct {
	node *Node
	err  error
}

func (w *walker) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	child := buildNode(node)
	child.Parent = w.node
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

// NextSibiling returns next sibling if exists
func (n *Node) NextSibiling() *Node {
	if n.Parent == nil {
		return nil
	}
	index := -1
	for i, child := range n.Parent.Children {
		if child == n {
			index = i
		}
	}
	nextIndex := index + 1
	if index < 0 || nextIndex >= len(n.Parent.Children) {
		return nil
	}
	return n.Parent.Children[nextIndex]
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
	entry := q.Entries[entryIndex]
	mustBeChild := entry.Combinator == ">"
	mustBeDecendant := mustBeChild || entry.Combinator == ""
	if mustBeChild && lastMatchedNodeDepth >= 0 && nodeDepth-lastMatchedNodeDepth > 1 {
		return true
	}
	if entry.Name == n.Name {
		if entryIndex+1 == len(q.Entries) {
			return cb(n)
		}
		nextEntry := q.Entries[entryIndex+1]
		switch nextEntry.Combinator {
		case ">", "":
			return n.applyChildren(q, entryIndex+1, nodeDepth+1, nodeDepth, cb)
		case "+":
			nextSibling := n.NextSibiling()
			if nextSibling == nil {
				return true
			}
			return nextSibling.apply(q, entryIndex+1, nodeDepth, nodeDepth, cb)
		default:
			return true
		}
	}
	if mustBeDecendant {
		return n.applyChildren(q, entryIndex, nodeDepth+1, lastMatchedNodeDepth, cb)
	}
	return true
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
