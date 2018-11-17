package query

import (
	"log"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
)

// Query represents parsed query ast
type Query struct {
	Pos lexer.Position

	Selectors []*Selector `parser:"[ @@ { ',' @@ } ]"`
}

// Selector represents multiple selectors node delimited by comma
type Selector struct {
	Pos lexer.Position

	SimpleSelectors []*SimpleSelector `parser:"@@ { @@ }"`
}

// SimpleSelector represents node selector with combinator and options
type SimpleSelector struct {
	Pos lexer.Position

	Combinator string                  `parser:"[ @('>' | '+' | '~') ]"`
	Name       string                  `parser:"[ @(Ident | '*') ]"`
	Options    []*SimpleSelectorOption `parser:"{ @@ }"`
}

// SimpleSelectorOption represents the option for SimpleSelector
type SimpleSelectorOption struct {
	Pos lexer.Position

	Attribute *Attribute `parser:"'[' @@ ']'"`
	Pseudo    *Pseudo    `parser:"| ':' @@"`
}

// Attribute represents the attribute option for SimpleSelector
type Attribute struct {
	Pos lexer.Position

	Name     string `parser:"@Ident"`
	Operator string `parser:"[ @('=' | ('~' '=') | ('|' '=') | ('^' '=') | ('$' '=') | ('*' '=')) ]"`
	Value    string `parser:"[ @(String | String2) ]"`
}

// Pseudo represents the pseudo option for SimpleSelector
type Pseudo struct {
	Pos lexer.Position

	FirstChild  *PseudoFirstChild  `parser:"@@"`
	FirstOfType *PseudoFirstOfType `parser:"| @@"`
	LastChild   *PseudoLastChild   `parser:"| @@"`
	LastOfType  *PseudoLastOfType  `parser:"| @@"`
}

// PseudoFirstChild represents the first-child pseudo
type PseudoFirstChild struct {
	Pos lexer.Position

	Name string `parser:"'first-child'"`
}

// PseudoFirstOfType represents the first-of-type pseudo
type PseudoFirstOfType struct {
	Pos lexer.Position

	Name string `parser:"'first-of-type'"`
}

// PseudoLastChild represents the last-child pseudo
type PseudoLastChild struct {
	Pos lexer.Position

	Name string `parser:"'last-child'"`
}

// PseudoLastOfType represents the last-of-type pseudo
type PseudoLastOfType struct {
	Pos lexer.Position

	Name string `parser:"'last-of-type'"`
}

var (
	queryLexer = lexer.Must(ebnf.New(`
Ident = (alpha | "_") { "_" | "-" |alpha | digit } .
String = "\"" { "\u0000"…"\uffff"-"\""-"\\" | "\\" any } "\"" .
String2 = "'" { "\u0000"…"\uffff"-"'"-"\\" | "\\" any } "'" .
Number = [ "-" | "+" ] digit { digit } .
Punct = "!"…"/" | ":"…"@" | "["…` + "\"`\"" + ` | "{"…"~" .
Whitespace = " " | "\t" | "\n" | "\r" .

alpha = "a"…"z" | "A"…"Z" .
digit = "0"…"9" .
any = "\u0000"…"\uffff" .
`))
	parser = participle.MustBuild(&Query{}, participle.Lexer(queryLexer), participle.Unquote("String", "String2"), participle.Elide("Whitespace"))
)

// Parse parses query and returns query ast
func Parse(q string) (*Query, error) {
	query := &Query{}
	err := parser.ParseString(q, query)
	if err != nil {
		return nil, err
	}
	return query, nil
}

// MustParse parses query and returns query ast
// If failed to parse, fatal occurs
func MustParse(q string) *Query {
	query, err := Parse(q)
	if err != nil {
		log.Fatalf("Cannot parse query. %v", err)
	}
	return query
}
