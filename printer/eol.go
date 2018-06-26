package printer

import (
	"fmt"
	"io"
)

func PrintEOL(w io.Writer) {
	fmt.Fprint(w, "\n")
}
