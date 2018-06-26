package printer

import (
	"fmt"
	"io"

	"go.uber.org/thriftrw/ast"
)

func PrintField(w io.Writer, i *Indent, f *ast.Field) {
	fmt.Fprintf(w, "%s%d: %s%s %s\n", i, f.ID, requiredness(f.Requiredness), f.Type, f.Name)
}

func requiredness(r ast.Requiredness) string {
	switch r {
	default:
		return ""
	case ast.Required:
		return "required "
	case ast.Optional:
		return "optional "
	}
}
