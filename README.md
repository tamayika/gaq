# GAQ(Go AST Query)

GAQ is the library to query `ast.Node` children like JavaScript `querySelector` or `querySelectorAll`.

## Install

```sh
go get github.com/tamayika/gaq/pkg/gaq
```

## Usage

Please refer [pkg/gaq/example_test.go](pkg/gaq/example_test.go)

## Query Specfication

Heavily inspired by CSS Selector.

```
NodeName [[Combinator] NodeName]
```

Here, NodeName is one of node type of [ast](https://golang.org/pkg/go/ast/).
For example, if you want to find `*ast.StructType`, NodeName is `StructType`.

If you don't know NodeName, VSCode extension [vscode-go-ast-explorer](https://github.com/tamayika/vscode-go-ast-explorer) will help you to find it out.

Please see below supported combinator.

|  Combinator  |           Meaning           |
| ------------ | --------------------------- |
| +            | Adjacent sibling combinator |
| >            | Child combinator            |
| (whitespace) | Descendant combinator       |
