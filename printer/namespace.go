package printer

import (
	"fmt"
	"io"

	"go.uber.org/thriftrw/ast"
)

func PrintNamespace(w io.Writer, i *Indent, n *ast.Namespace, eol func(w io.Writer)) {
	fmt.Fprintf(w, "%snamespace %s %s", i, n.Scope, n.Name)
	eol(w)
}
