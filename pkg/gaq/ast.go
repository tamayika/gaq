package gaq

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"strings"

	"github.com/tamayika/gaq/pkg/gaq/query"
)

// Node represents traversible ast.Node
type Node struct {
	Type     string  `json:"type"`
	Pos      int     `json:"pos"`
	End      int     `json:"end"`
	Children []*Node `json:"children,omitempty"`

	Parent *Node    `json:"-"`
	Index  int      `json:"-"`
	Node   ast.Node `json:"-"`
	Name   string   `json:"-"`
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
	if child.Parent != nil {
		child.Index = len(child.Parent.Children)
	}
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
		Index:    -1,
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

// SameTypeChildren returns the same type of target Node among children
func (n *Node) SameTypeChildren(target *Node) []*Node {
	ret := []*Node{}
	for _, child := range n.Children {
		if child.Type == target.Type {
			ret = append(ret, child)
		}
	}
	return ret
}

// NextSibiling returns next sibling if exists
func (n *Node) NextSibiling() *Node {
	if n.Parent == nil || n.Index < 0 {
		return nil
	}
	nextIndex := n.Index + 1
	if nextIndex >= len(n.Parent.Children) {
		return nil
	}
	return n.Parent.Children[nextIndex]
}

// NextSibilings returns next siblings if exists
func (n *Node) NextSibilings() []*Node {
	if n.Parent == nil {
		return nil
	}
	nextIndex := n.Index + 1
	if nextIndex >= len(n.Parent.Children) {
		return nil
	}
	return n.Parent.Children[nextIndex:]
}

// QuerySelector queries to node and return first matched node
func (n *Node) QuerySelector(q *query.Query) ast.Node {
	var firstNode ast.Node
	for _, selector := range q.Selectors {
		n.apply(selector, 0, 0, -1, func(n *Node) bool {
			firstNode = n.Node
			return false
		})
		if firstNode != nil {
			break
		}
	}
	return firstNode
}

// QuerySelectorAll queries to node and return all matched nodes
func (n *Node) QuerySelectorAll(q *query.Query) []ast.Node {
	nodes := []ast.Node{}
	addedNodes := map[*Node]bool{}
	for _, selector := range q.Selectors {
		n.apply(selector, 0, 0, -1, func(n *Node) bool {
			if _, ok := addedNodes[n]; !ok {
				nodes = append(nodes, n.Node)
				addedNodes[n] = true
			}
			return true
		})
	}
	return nodes
}

type callback func(n *Node) bool

func (n *Node) apply(s *query.Selector, selectorIndex int, nodeDepth int, lastMatchedNodeDepth int, cb callback) bool {
	ss := s.SimpleSelectors[selectorIndex]
	mustBeChild := ss.Combinator == ">"
	mustBeDecendant := mustBeChild || ss.Combinator == ""
	if mustBeChild && lastMatchedNodeDepth >= 0 && nodeDepth-lastMatchedNodeDepth > 1 {
		return true
	}
	if n.isMatchSimpleSelector(ss) {
		if selectorIndex+1 == len(s.SimpleSelectors) {
			continues := cb(n)
			if !continues {
				return false
			}
			// if selector length is 1, we must continue to query
			if len(s.SimpleSelectors) > 1 {
				return continues
			}
		} else {
			nextEntry := s.SimpleSelectors[selectorIndex+1]
			switch nextEntry.Combinator {
			case ">", "":
				return n.applyChildren(s, selectorIndex+1, nodeDepth+1, nodeDepth, cb)
			case "+":
				nextSibling := n.NextSibiling()
				if nextSibling == nil {
					return true
				}
				return nextSibling.apply(s, selectorIndex+1, nodeDepth, nodeDepth, cb)
			case "~":
				nextSibilings := n.NextSibilings()
				if nextSibilings == nil {
					return true
				}
				for _, nextSibling := range nextSibilings {
					continues := nextSibling.apply(s, selectorIndex+1, nodeDepth, nodeDepth, cb)
					if !continues {
						return false
					}
				}
				return true
			default:
				return true
			}
		}
	}
	if mustBeDecendant {
		return n.applyChildren(s, selectorIndex, nodeDepth+1, lastMatchedNodeDepth, cb)
	}
	return true
}

func (n *Node) isMatchSimpleSelector(ss *query.SimpleSelector) bool {
	return (ss.Name == n.Name || ss.Name == "*" || ss.Name == "") && n.isMatchOptions(ss.Options)
}

func (n *Node) applyChildren(s *query.Selector, selectorIndex int, nodeDepth int, lastMatchedNodeDepth int, cb callback) bool {
	for _, childNode := range n.Children {
		continues := childNode.apply(s, selectorIndex, nodeDepth, lastMatchedNodeDepth, cb)
		if !continues {
			return false
		}
	}
	return true
}

func (n *Node) isMatchOptions(opts []*query.SimpleSelectorOption) bool {
	if opts == nil {
		return true
	}
	for _, opt := range opts {
		if !n.isMatchOption(opt) {
			return false
		}
	}
	return true
}

func (n *Node) isMatchOption(opt *query.SimpleSelectorOption) bool {
	return n.isMatchOptionAttribute(opt.Attribute) && n.isMatchOptionPseudo(opt.Pseudo)
}

func (n *Node) isMatchOptionAttribute(oa *query.Attribute) bool {
	if oa == nil {
		return true
	}
	nodeValue := reflect.ValueOf(n.Node).Elem()
	field := nodeValue.FieldByName(oa.Name)
	if !field.IsValid() {
		return false
	}
	v := field.Interface()
	value, ok := v.(string)
	if !ok {
		return false
	}
	switch oa.Operator {
	case "=":
		return oa.Value == value
	case "~=":
		splitted := strings.Split(value, " ")
		found := false
		for _, s := range splitted {
			if s == oa.Value {
				found = true
				break
			}
		}
		return found
	case "|=":
		return strings.Contains(value, fmt.Sprintf("%s-", oa.Value))
	case "^=":
		return strings.HasPrefix(value, oa.Value)
	case "$=":
		return strings.HasSuffix(value, oa.Value)
	case "*=":
		return strings.Contains(value, oa.Value)
	}
	return false
}

func (n *Node) isMatchOptionPseudo(op *query.Pseudo) bool {
	if op == nil {
		return true
	}

	if op.Empty != nil {
		if len(n.Children) == 0 {
			return true
		}
		for _, child := range n.Children {
			if _, ok := child.Node.(*ast.Comment); ok {
				continue
			}
			if _, ok := child.Node.(*ast.CommentGroup); ok {
				continue
			}
			return false
		}
		return true
	} else if op.FirstChild != nil {
		return n.Index == 0
	} else if op.FirstOfType != nil {
		if n.Parent != nil {
			for i, child := range n.Parent.SameTypeChildren(n) {
				if child == n {
					return i == 0
				}
			}
		}
	} else if op.Has != nil {
		for _, selector := range op.Has.Selectors {
			found := false
			n.apply(selector, 0, 1, 0, func(n *Node) bool {
				found = true
				return true
			})
			if found {
				return true
			}
		}
	} else if op.Is != nil {
		for _, selector := range op.Is.Selectors {
			if n.isMatchSimpleSelector(selector.SimpleSelectors[0]) {
				return true
			}
		}
	} else if op.LastChild != nil {
		if n.Parent != nil {
			return n.Index == len(n.Parent.Children)-1
		}
	} else if op.LastOfType != nil {
		if n.Parent != nil {
			children := n.Parent.SameTypeChildren(n)
			for i, child := range children {
				if child == n {
					return i == len(children)-1
				}
			}
		}
	} else if op.Not != nil {
		for _, selector := range op.Not.Selectors {
			if n.isMatchSimpleSelector(selector.SimpleSelectors[0]) {
				return false
			}
		}
		return true
	} else if op.Root != nil {
		return n.Parent == nil
	}
	return false
}
