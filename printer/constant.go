package printer

import (
	"fmt"
	"io"

	"go.uber.org/thriftrw/ast"
)

func PrintConstant(w io.Writer, i *Indent, c *ast.Constant) {
	PrintDoc(w, i, c.Doc)
	fmt.Fprintf(w, "%sconst %s %s = %s\n", i, c.Type, c.Name, constantValue(c.Value))
}

func constantValue(value ast.ConstantValue) string {
	switch v := value.(type) {
	default:
		return fmt.Sprint(v)
		// case ast.ConstantBoolean,
		// 	ast.ConstantInteger,
		// 	ast.ConstantString,
		// 	ast.ConstantDouble:
		// 	return fmt.Sprint(v)
	case ast.ConstantMap:
		return "{}"
	case ast.ConstantList:
		return "[]"
	case ast.ConstantReference:
		return v.Name
	}
}
