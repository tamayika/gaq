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
Query:
    Selector [, Selector]

Selector:
    SimpleSelector [Combinator Selector]

SimpleSelector:
    NodeName[Attribute][Pseudo]

Attribute:
    '[' Field [ AttributeOperator Value ] ']'

Pseudo:
    ':' Name [ '(' Expression ')' ]
```

Full used example
```
    File > Ident[Name*='test']:first-child, File > Ident[Name*='test']:last-child
```

Here, NodeName is one of node type of [ast](https://golang.org/pkg/go/ast/).
For example, if you want to find `*ast.StructType`, NodeName is `StructType`.

If you don't know NodeName, VSCode extension [vscode-go-ast-explorer](https://github.com/tamayika/vscode-go-ast-explorer) will help you to find it out.

### Supported combinators.

|  Combinator  |            Name             |                                                 Meaning                                                 |
| ------------ | --------------------------- | ------------------------------------------------------------------------------------------------------- |
| +            | Adjacent sibling combinator | The second node directly follows the first, and both share the same parent.                             |
| ~            | General sibling combinator  | The second node follows the first (though not necessarily immediately), and both share the same parent. |
| >            | Child combinator            | Selects nodes that are direct children of the first node.                                               |
| (whitespace) | Descendant combinator       | Selects nodes that are descendants of the first node.                                                   |

### Supported Attribute Syntax

|    Syntax    |                                                                 Meaning                                                                 |
| ------------ | --------------------------------------------------------------------------------------------------------------------------------------- |
| `[f]`        | Represents Node with an field name of f.                                                                                                |
| `[f=value]`  | Represents Node with an field name of f whose value is exactly value.                                                                   |
| `[f~=value]` | Represents Node with an field name of f whose value is a whitespace-separated list of words, one of which is exactly value.             |
| `[f|=value]` | Represents Node with an field name of f whose value can be exactly value or can begin with value immediately followed by a hyphen, `-`. |
| `[f^=value]` | Represents Node with an field name of f whose value is prefixed (preceded) by value.                                                    |
| `[f$=value]` | Represents Node with an field name of f whose value is suffixed (followed) by value.                                                    |
| `[f*=value]` | Represents Node with an field name of f whose value contains at least one occurrence of value within the string.                        |

### Supported Pseudo Class

|     Syntax     |                          Meaning                          |
| -------------- | --------------------------------------------------------- |
| `:first-child` | Represents the first node among a group of sibling nodes. |
| `:last-child`  | Represents the last node among a group of sibling nodes.  |