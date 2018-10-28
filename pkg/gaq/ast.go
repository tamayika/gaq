package gaq

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

type Node struct {
	Type     string  `json:"type"`
	Pos      int     `json:"pos"`
	End      int     `json:"end"`
	Children []*Node `json:"children"`
}

// ParseSource parses source and returns ast
func ParseSource(source string) (*Node, error) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	w := &walker{}
	ast.Walk(w, f)
	return w.node, nil
}

type walker struct {
	node *Node
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
	node := &Node{Type: nodeType(n), Children: []*Node{}, Pos: int(n.Pos()), End: int(n.End())}
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
