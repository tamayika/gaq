package parser

import "github.com/tamayika/gaq/pkg/gaq"

// Parser represents query parser
type Parser struct {
	node *gaq.Node
}

// NewParser creates Parser and returns it
func NewParser(node *gaq.Node) *Parser {
	return &Parser{node}
}
