package printer

import (
	"fmt"
	"io"
	"strings"
)

func PrintDoc(w io.Writer, i *Indent, d string) {
	if len(d) == 0 {
		return
	}
	fmt.Fprintf(w, "%s/**\n", i)
	for _, line := range strings.Split(d, "\n") {
		fmt.Fprintf(w, "%s * %s\n", i, line)
	}
	fmt.Fprintf(w, "%s */\n", i)
}
