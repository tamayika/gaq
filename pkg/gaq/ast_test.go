package gaq

import (
	"go/ast"
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
			assert.NotNil(t, got)
			assert.IsType(t, tt.want, got)
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
			"File",
			MustParse(`package main`),
			args{
				query.MustParse("File"),
			},
			[]ast.Node{&ast.File{}},
		},
		{
			"File Ident",
			MustParse(`package main
			`),
			args{
				query.MustParse("File Ident"),
			},
			[]ast.Node{&ast.Ident{Name: "main"}},
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
				&ast.Ident{Name: "main"},
				&ast.Ident{Name: "f"},
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
				&ast.Ident{Name: "main"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.QuerySelectorAll(tt.args.q)
			assert.NotNil(t, got)
			assert.Len(t, got, len(tt.want))
			for i := 0; i < len(tt.want); i++ {
				assert.IsType(t, tt.want[i], got[i])
			}
		})
	}
}
