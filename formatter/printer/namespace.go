package printer

import (
	"fmt"
	"io"

	"go.uber.org/thriftrw/ast"
)

func PrintNamespace(w io.Writer, i Indent, n *ast.Namespace) {
	fmt.Fprintf(w, "%snamespace %s %s\n", i, n.Scope, n.Name)
}
