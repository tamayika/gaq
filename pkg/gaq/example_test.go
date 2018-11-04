package gaq_test

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/tamayika/gaq/pkg/gaq"
	"github.com/tamayika/gaq/pkg/gaq/query"
)

func Example() {
	// 1. parse Query
	q, err := query.Parse("File > Ident")
	if err != nil {
		log.Fatalf("Cannot parse query. %v", err)
	}
	// or MustParse fatals if parse failed
	// q := query.MustParse("File > Ident")

	// 2. parse source
	source := `package main`
	node, err := gaq.Parse(source)
	if err != nil {
		log.Fatalf("Cannot parse node. %v", err)
	}

	// 3. run QuerySelector or QuerySelectorAll
	fmt.Printf("%v\n", node.QuerySelector(q))
	fmt.Printf("%v\n", node.QuerySelectorAll(q))
}

// Modify ast.Node to export all field in all struct type.
func ExampleNode_QuerySelectorAll() {
	source := `package main

type User struct {
	name string
	age  int
}`
	q := query.MustParse("StructType > FieldList > Field")

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		log.Fatalf("Cannot parse source. %v", err)
	}
	node := gaq.MustParseNode(f)
	for _, n := range node.QuerySelectorAll(q) {
		field := n.(*ast.Field)
		for i, name := range field.Names {
			field.Names[i].Name = strings.Title(name.Name)
		}
	}
	var buf bytes.Buffer
	err = format.Node(&buf, fset, f)
	if err != nil {
		log.Fatalf("Cannot format node. %v", err)
	}
	fmt.Println(buf.String())
	// Output:
	// package main
	//
	// type User struct {
	//	Name string
	//	Age  int
	// }
}
