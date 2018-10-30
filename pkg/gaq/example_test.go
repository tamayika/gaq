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

var source = `package main

type User struct {
	name string
	age  int
}
`

var q = query.MustParse("StructType > FieldList > Field")

func ExampleNode_QuerySelectorAll_export_all_fields() {
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
