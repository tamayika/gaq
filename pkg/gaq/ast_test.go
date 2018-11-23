package gaq

import (
	"go/ast"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tamayika/gaq/pkg/gaq/query"
)

func TestNode_QuerySelector(t *testing.T) {
	type args struct {
		q *query.Query
	}
	tests := []struct {
		name string
		n    *Node
		args args
		want ast.Node
	}{
		{
			"Not match",
			MustParse(`package main
			func f() {

			}`),
			args{
				query.MustParse("GenDecl"),
			},
			nil,
		},
		{
			"File",
			MustParse(`package main`),
			args{
				query.MustParse("File"),
			},
			&ast.File{},
		},
		{
			"File Ident",
			MustParse(`package main
			`),
			args{
				query.MustParse("File Ident"),
			},
			&ast.Ident{Name: "main"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.QuerySelector(tt.args.q)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.IsType(t, tt.want, got)
			}
		})
	}
}

func TestNode_QuerySelectorAll(t *testing.T) {
	type args struct {
		q *query.Query
	}
	tests := []struct {
		name string
		n    *Node
		args args
		want []ast.Node
	}{
		{

			"Not match",
			MustParse(`package main
			func f() {

			}`),
			args{
				query.MustParse("GenDecl"),
			},
			[]ast.Node{},
		},
		{
			"*",
			MustParse(`package main`),
			args{
				query.MustParse("*"),
			},
			[]ast.Node{
				&ast.File{
					Name: &ast.Ident{
						Name: "main",
					},
				},
				&ast.Ident{
					Name: "main",
				},
			},
		},
		{
			"File",
			MustParse(`package main`),
			args{
				query.MustParse("File"),
			},
			[]ast.Node{
				&ast.File{
					Name: &ast.Ident{
						Name: "main",
					},
				},
			},
		},
		{
			"File Ident",
			MustParse(`package main
			`),
			args{
				query.MustParse("File Ident"),
			},
			[]ast.Node{
				&ast.Ident{
					Name: "main",
				},
			},
		},
		{
			"File, Ident",
			MustParse(`package main
			`),
			args{
				query.MustParse("File, Ident"),
			},
			[]ast.Node{
				&ast.File{
					Name: &ast.Ident{
						Name: "main",
					},
				},
				&ast.Ident{
					Name: "main",
				},
			},
		},
		{
			"File, File",
			MustParse(`package main
			`),
			args{
				query.MustParse("File, File"),
			},
			[]ast.Node{
				&ast.File{
					Name: &ast.Ident{
						Name: "main",
					},
				},
			},
		},
		{
			"File Ident",
			MustParse(`package main
			import (
				"os"
			)
			func f() {

			}
			`),
			args{
				query.MustParse("File Ident"),
			},
			[]ast.Node{
				&ast.Ident{
					Name: "main",
				},
				&ast.Ident{
					Name: "f",
				},
			},
		},
		{
			"File > Ident",
			MustParse(`package main
			import (
				"os"
			)
			func f() {

			}
			`),
			args{
				query.MustParse("File > Ident"),
			},
			[]ast.Node{
				&ast.Ident{
					Name: "main",
				},
			},
		},
		{
			"GenDecl + FuncDecl",
			MustParse(`package main
			import (
				"os"
			)
			func f() {

			}
			`),
			args{
				query.MustParse("GenDecl + FuncDecl"),
			},
			[]ast.Node{
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "f",
					},
				},
			},
		},
		{
			"Not matched FuncDecl + FuncDecl",
			MustParse(`package main
			import (
				"os"
			)
			func f() {

			}
			`),
			args{
				query.MustParse("FuncDecl + FuncDecl"),
			},
			[]ast.Node{},
		},
		{
			"GenDecl ~ FuncDecl",
			MustParse(`package main
			import (
				"os"
			)
			func f() {

			}
			func f2() {

			}
			`),
			args{
				query.MustParse("GenDecl ~ FuncDecl"),
			},
			[]ast.Node{
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "f",
					},
				},
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "f2",
					},
				},
			},
		},
		{
			"Not matched GenDecl ~ FuncDecl",
			MustParse(`package main
			import (
				"os"
			)
			`),
			args{
				query.MustParse("GenDecl ~ FuncDecl"),
			},
			[]ast.Node{},
		},
		{
			"File Ident[Name='a']",
			MustParse(`package main
			var a string
			var ab string
			var b string
			var ba string
			`),
			args{
				query.MustParse("File Ident[Name='a']"),
			},
			[]ast.Node{
				&ast.Ident{
					Name: "a",
				},
			},
		},
		{
			"File Ident[Name^='a']",
			MustParse(`package main
			var a string
			var ab string
			var b string
			var ba string
			`),
			args{
				query.MustParse("File Ident[Name^='a']"),
			},
			[]ast.Node{
				&ast.Ident{
					Name: "a",
				},
				&ast.Ident{
					Name: "ab",
				},
			},
		},
		{
			"File Ident[Name$='a']",
			MustParse(`package main
			var a string
			var ab string
			var b string
			var ba string
			`),
			args{
				query.MustParse("File Ident[Name$='a']"),
			},
			[]ast.Node{
				&ast.Ident{
					Name: "a",
				},
				&ast.Ident{
					Name: "ba",
				},
			},
		},
		{
			"File Ident[Name*='a']",
			MustParse(`package foo
			var a string
			var ab string
			var b string
			var ba string
			`),
			args{
				query.MustParse("File Ident[Name*='a']"),
			},
			[]ast.Node{
				&ast.Ident{
					Name: "a",
				},
				&ast.Ident{
					Name: "ab",
				},
				&ast.Ident{
					Name: "ba",
				},
			},
		},
		{
			"StructType FieldList:empty",
			MustParse(`package foo

			type a struct {
			}
			
			type b struct {
				f string
			}
			`),
			args{
				query.MustParse("StructType FieldList:empty"),
			},
			[]ast.Node{
				&ast.FieldList{},
			},
		},
		{
			"File StructType Field:first-child",
			MustParse(`package foo

			type s struct {
				hoge string
				huga string
			}
			`),
			args{
				query.MustParse("File StructType Field:first-child"),
			},
			[]ast.Node{
				&ast.Field{
					Names: []*ast.Ident{
						&ast.Ident{Name: "hoge"},
					},
					Type: &ast.Ident{Name: "string"},
				},
			},
		},
		{
			"Not matched File > FuncDecl:first-child",
			MustParse(`package foo
			import (
				"os"
			)
			func f() {

			}
			`),
			args{
				query.MustParse("File > FuncDecl:first-child"),
			},
			[]ast.Node{},
		},
		{
			"File > FuncDecl:first-of-type",
			MustParse(`package foo
			import (
				"os"
			)
			func f() {

			}
			`),
			args{
				query.MustParse("File > FuncDecl:first-of-type"),
			},
			[]ast.Node{
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "f",
					},
				},
			},
		},
		{
			"File TypeSpec:has(Field)",
			MustParse(`package foo

			type s struct {
				hoge string
				huga string
			}

			type ss struct {

			}
			`),
			args{
				query.MustParse("File TypeSpec:has(Field)"),
			},
			[]ast.Node{
				&ast.TypeSpec{
					Name: &ast.Ident{Name: "s"},
				},
			},
		},
		{
			"Not matched File TypeSpec:has(>Field)",
			MustParse(`package foo

			type s struct {
				hoge string
				huga string
			}

			type ss struct {

			}
			`),
			args{
				query.MustParse("File TypeSpec:has(>Field)"),
			},
			[]ast.Node{},
		},
		{
			"*:is(InterfaceType, StructType)",
			MustParse(`package foo
			type I interface {}
			type S struct {}
			`),
			args{
				query.MustParse("*:is(InterfaceType, StructType)"),
			},
			[]ast.Node{
				&ast.InterfaceType{},
				&ast.StructType{},
			},
		},
		{
			"File StructType Field:last-child",
			MustParse(`package foo

			type s struct {
				hoge string
				huga string
			}
			`),
			args{
				query.MustParse("File StructType Field:last-child"),
			},
			[]ast.Node{
				&ast.Field{
					Names: []*ast.Ident{
						&ast.Ident{Name: "huga"},
					},
					Type: &ast.Ident{Name: "string"},
				},
			},
		},
		{
			"Not matched File > GenDecl:last-child",
			MustParse(`package foo
			import (
				"os"
			)
			func f() {

			}
			`),
			args{
				query.MustParse("File > GenDecl:last-child"),
			},
			[]ast.Node{},
		},
		{
			"File > GenDecl:last-of-type",
			MustParse(`package foo
			import (
				"os"
			)
			func f() {

			}
			`),
			args{
				query.MustParse("File > GenDecl:last-of-type"),
			},
			[]ast.Node{
				&ast.GenDecl{},
			},
		},
		{
			"TypeSpec>*:not(InterfaceType):not(Ident)",
			MustParse(`package foo
			type I interface {}
			type S struct {}
			`),
			args{
				query.MustParse("TypeSpec>*:not(InterfaceType):not(Ident)"),
			},
			[]ast.Node{
				&ast.StructType{},
			},
		},
		{
			"TypeSpec>*:not(InterfaceType, Ident)",
			MustParse(`package foo
			type I interface {}
			type S struct {}
			`),
			args{
				query.MustParse("TypeSpec>*:not(InterfaceType, Ident)"),
			},
			[]ast.Node{
				&ast.StructType{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.QuerySelectorAll(tt.args.q)
			assert.NotNil(t, got)
			l := len(tt.want)
			assert.Len(t, got, l)
			if l > len(got) {
				l = len(got)
			}
			for i := 0; i < l; i++ {
				equalIdent(t, tt.want[i], got[i])
			}
		})
	}
}

func equalIdent(t *testing.T, n1 ast.Node, n2 ast.Node) bool {
	sameType := assert.IsType(t, n1, n2)
	if !sameType {
		return false
	}
	n1Ident, n1Identok := n1.(*ast.Ident)
	n2Ident, n2Identok := n2.(*ast.Ident)
	if n1Identok && n2Identok {
		return assert.Equal(t, n1Ident.Name, n2Ident.Name)
	}
	n1Value := reflect.ValueOf(n1).Elem()
	n2Value := reflect.ValueOf(n2).Elem()
	for i := 0; i < n1Value.NumField(); i++ {
		n1Field := n1Value.Field(i)
		n2Field := n2Value.Field(i)

		switch n1Field.Type().Kind() {
		case reflect.Ptr, reflect.Interface:
			if n1Field.IsNil() && n2Field.IsNil() {
				continue
			}
			n1FieldNode, n1Fieldok := n1Field.Interface().(ast.Node)
			n2FieldNode, n2Fieldok := n2Field.Interface().(ast.Node)
			if n1Field.IsNil() || n2Field.IsNil() {
				isIdent := false
				if n1FieldNode != nil {
					_, ok := n1FieldNode.(*ast.Ident)
					isIdent = isIdent || ok
				}
				if n2FieldNode != nil {
					_, ok := n2FieldNode.(*ast.Ident)
					isIdent = isIdent || ok
				}
				if isIdent {
					assert.Equal(t, n1FieldNode, n2FieldNode)
				}
				continue
			}

			if n1Fieldok && n2Fieldok {
				equalIdent(t, n1FieldNode, n2FieldNode)
			} else {
				assert.Equal(t, n1Field.Interface(), n2Field.Interface())
			}
		case reflect.Slice, reflect.Array:
			if n1Field.IsNil() || n2Field.IsNil() {
				continue
			}
			l := n1Field.Len()
			if l > n2Field.Len() {
				l = n2Field.Len()
			}
			for i := 0; i < l; i++ {
				n1FieldNode, n1Fieldok := n1Field.Index(i).Interface().(ast.Node)
				n2FieldNode, n2Fieldok := n2Field.Index(i).Interface().(ast.Node)
				if n1Fieldok && n2Fieldok {
					equalIdent(t, n1FieldNode, n2FieldNode)
				} else {
					assert.Equal(t, n1Field.Index(i).Interface(), n2Field.Index(i).Interface())
				}
			}
		}
	}
	return true
}
