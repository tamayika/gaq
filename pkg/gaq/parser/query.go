package parser

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
)

// Query represents parsed query ast
type Query struct {
	Pos lexer.Position

	Entries []*Entry `parser:"{ @@ }"`
}

// Entry represents ast node name
type Entry struct {
	Pos lexer.Position

	Combinator string `parser:"[ @('>' | '+') ]"`
	Name       string `parser:"@Ident"`
}

var (
	queryLexer = lexer.Must(ebnf.New(`
Ident = (alpha | "_") { "_" | alpha | digit } .
String = "\"" { "\u0000"…"\uffff"-"\""-"\\" | "\\" any } "\"" .
Punct = "!"…"/" | ":"…"@" | "["…`+"\"`\""+` | "{"…"~" .
Whitespace = " " | "\t" | "\n" | "\r" .

alpha = "a"…"z" | "A"…"Z" .
digit = "0"…"9" .
any = "\u0000"…"\uffff" .
`, ebnf.Elide("Whitespace")))
	parser = participle.MustBuild(&Query{}, participle.Lexer(queryLexer), participle.Unquote("String"))
)

// ParseQuery parses query and returns query ast
func ParseQuery(q string) (*Query, error) {
	query := &Query{}
	err := parser.ParseString(q, query)
	if err != nil {
		return nil, err
	}
	return query, nil
}
