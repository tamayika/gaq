package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/tamayika/gaq/pkg/gaq"
	"github.com/tamayika/gaq/pkg/gaq/query"
)

var version = "dev"

func printText(source []byte, fset *token.FileSet, nodes []ast.Node) {
	for _, node := range nodes {
		pos := fset.Position(node.Pos())
		end := fset.Position(node.End())
		fmt.Println(string(source[pos.Offset:end.Offset]))
	}
}

func printPos(nodes []ast.Node) {
	for _, node := range nodes {
		fmt.Printf("%d,%d\n", node.Pos(), node.End())
	}
}

func main() {
	var format string

	rootCmd := &cobra.Command{
		Use:   "gaq <Query>",
		Short: "gaq is the cli tool to query ast node. STDIN needed as go code.",
		Long: `gaq is the cli tool to query ast node. 
Typical usage is

  cat <go file path> | gaq <Query>

Please see details at https://github.com/tamayika/gaq`,
		Args:    cobra.ExactArgs(1),
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			data, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("Cannot read data from stdin. %v", err)
			}
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "", string(data), parser.ParseComments)
			if err != nil {
				log.Fatalf("Cannot parse source. %v", err)
			}
			node := gaq.MustParseNode(f)

			q := query.MustParse(args[0])
			nodes := node.QuerySelectorAll(q)
			switch format {
			case "text":
				printText(data, fset, nodes)
			case "pos":
				printPos(nodes)
			default:
				log.Fatalf("Format: %s is not supported.", format)
			}
		},
	}
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "Output format, 'text' or 'pos'. Default is 'text'")
	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
