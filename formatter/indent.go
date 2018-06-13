package formatter

import "strings"

const indent = "  "

type Indent int

func (i Indent) String() string {
	return strings.Repeat(indent, int(i))
}
