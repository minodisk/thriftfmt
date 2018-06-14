package printer

import (
	"fmt"
	"io"

	"go.uber.org/thriftrw/ast"
)

func PrintEnum(w io.Writer, i *Indent, c *ast.Enum) {
	PrintDoc(w, i, c.Doc)
	fmt.Fprintf(w, "%senum %s {", i, c.Name)
	if len(c.Items) == 0 {
		fmt.Fprintf(w, "}")
	} else {
		fmt.Fprintf(w, "\n")
		*i++
		for _, item := range c.Items {
			//fmt.Println(item.Line)
			fmt.Fprintf(w, "%s%s = %d\n", i, item.Name, *item.Value)
		}
		*i--
		fmt.Fprintf(w, "%s}\n\n", i)
	}
}
