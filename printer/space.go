package printer

import (
	"fmt"
	"io"
)

func PrintSpace(w io.Writer) {
	fmt.Fprint(w, " ")
}
