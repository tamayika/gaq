package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

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

func replaceByCommand(source []byte, fset *token.FileSet, nodes []ast.Node, commands []string) []byte {
	ret := []byte{}
	var lastNode ast.Node
	for _, node := range nodes {
		pos := fset.Position(node.Pos())
		end := fset.Position(node.End())
		nodeText := source[pos.Offset:end.Offset]
		var stderr bytes.Buffer
		cmd := exec.Command(commands[0], commands[1:]...)
		cmd.Stderr = &stderr
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatalf("Cannot get stdin pipe. %v", err)
		}
		io.WriteString(stdin, string(nodeText))
		stdin.Close()
		replacedText, err := cmd.Output()
		if err != nil {
			log.Fatalf("Command failed.\nerr: %v\nstderr: %s\nnodeText: %s", err, strings.TrimSuffix(string(stderr.String()), "\n"), string(nodeText))
		}
		if lastNode == nil {
			ret = append(ret, source[:pos.Offset]...)
		} else {
			ret = append(ret, source[fset.Position(lastNode.End()).Offset:pos.Offset]...)
		}
		ret = append(ret, replacedText...)
		lastNode = node
	}
	if lastNode == nil {
		ret = source
	} else {
		ret = append(ret, source[fset.Position(lastNode.End()).Offset:]...)
	}
	return ret
}

func main() {
	var format string
	var mode string

	rootCmd := &cobra.Command{
		Use:   "gaq <Query>",
		Short: "gaq is the cli tool to query ast node. STDIN needed as go code.",
		Long: `gaq is the cli tool to query ast node. 
Typical usage is

  cat <go file path> | gaq <Query>
  cat <go file path> | gaq -m replace <Query> <Replace command>

Please see details at https://github.com/tamayika/gaq`,
		Args:    cobra.MinimumNArgs(1),
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
			switch mode {
			case "filter":
				switch format {
				case "text":
					printText(data, fset, nodes)
				case "pos":
					printPos(nodes)
				default:
					log.Fatalf("Format: %s is not supported.", format)
				}
			case "replace":
				if len(args) < 2 {
					log.Fatalf("One or more command and args are expected in replace mode.")
				}
				fmt.Println(string(replaceByCommand(data, fset, nodes, args[1:])))
			default:
				log.Fatalf("Mode: %s is not supported.", mode)
			}
		},
	}
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "Output format, 'text' or 'pos'. Default is 'text'")
	rootCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "filter", "Execution mode, 'filter' or 'replace'. Default is 'filter'")
	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
