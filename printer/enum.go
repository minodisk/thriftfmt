package printer

import (
	"fmt"
	"io"

	"go.uber.org/thriftrw/ast"
)

func PrintEnum(w io.Writer, indent *Indent, e *ast.Enum) {
	PrintDoc(w, indent, e.Doc)
	fmt.Fprintf(w, "%senum %s {", indent, e.Name)
	if len(e.Items) == 0 {
		fmt.Fprintf(w, "}\n")
		return
	}
	fmt.Fprintf(w, "\n")
	*indent++
	last := len(e.Items) - 1
	for i, item := range e.Items {
		fmt.Fprintf(w, "%s%s", indent, item.Name)
		if item.Value != nil {
			fmt.Fprintf(w, " = %d", *item.Value)
		}
		if len(item.Annotations) > 0 {
			fmt.Fprintf(w, " %s", ast.FormatAnnotations(item.Annotations))
		}
		if i != last {
			fmt.Fprintf(w, ",")
		}
		fmt.Fprintf(w, "\n")
	}
	*indent--
	fmt.Fprintf(w, "%s}", indent)
	if len(e.Annotations) > 0 {
		fmt.Fprintf(w, " %s", ast.FormatAnnotations(e.Annotations))
	}
	fmt.Fprint(w, "\n")
}
