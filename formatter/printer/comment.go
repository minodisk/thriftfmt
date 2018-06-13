package printer

import (
	"fmt"
	"io"

	"github.com/minodisk/thriftfmt/token"
)

func PrintComment(w io.Writer, i Indent, c *token.Comment) {
	fmt.Fprintf(w, "%s\n\n", c.Body)
}
