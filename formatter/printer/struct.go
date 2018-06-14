package printer

import (
	"fmt"
	"io"

	"go.uber.org/thriftrw/ast"
)

func PrintStruct(w io.Writer, i *Indent, s *ast.Struct) {
	PrintDoc(w, i, s.Doc)
	fmt.Fprintf(w, "%sstruct %s {", i, s.Name)
	if len(s.Fields) > 0 {
		fmt.Fprintf(w, "\n")
		*i++
		for _, field := range s.Fields {
			PrintField(w, i, field)
		}
		*i--
	}
	fmt.Fprintf(w, "%s}\n\n", i)
}
