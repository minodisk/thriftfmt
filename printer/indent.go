package printer

import "strings"

const indent = "  "

type Indent int

func NewIndent(base int) *Indent {
	i := Indent(base)
	return &i
}

func (i *Indent) Increment() {
	*i++
}

func (i *Indent) Decrement() {
	*i--
}

func (i *Indent) String() string {
	return strings.Repeat(indent, int(*i))
}
