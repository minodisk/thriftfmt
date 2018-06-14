package printer

import (
	"fmt"
	"io"

	"go.uber.org/thriftrw/ast"
)

func PrintInclude(w io.Writer, i *Indent, include *ast.Include) {
	fmt.Fprintf(w, "%sinclude ", i)
	if include.Name != "" {
		fmt.Fprintf(w, "%s ", include.Name)
	}
	fmt.Fprintf(w, "\"%s\"\n", include.Path)
}
