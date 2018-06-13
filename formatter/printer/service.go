package printer

import (
	"fmt"
	"io"

	"go.uber.org/thriftrw/ast"
)

func PrintService(w io.Writer, i Indent, s *ast.Service) {
	PrintDoc(w, i, s.Doc)
	fmt.Fprintf(w, "service %s {", s.Name)
	if len(s.Functions) > 0 {
		fmt.Fprint(w, "\n\n")
		i++
		for _, function := range s.Functions {
			//fmt.Println(function.Line)
			PrintDoc(w, i, function.Doc)
			fmt.Fprintf(w, "%s%s%s%s(", i, oneWay(function.OneWay), returnType(function.ReturnType), function.Name)
			if len(function.Parameters) > 0 {
				fmt.Fprint(w, "\n")
				i++
				for _, parameter := range function.Parameters {
					//fmt.Println(parameter.Line)
					PrintField(w, i, parameter)
				}
				i--
				fmt.Fprintf(w, "%s)", i)
			} else {
				fmt.Fprintf(w, ")")
			}

			if len(function.Exceptions) > 0 {
				fmt.Fprintf(w, " throws (\n")
				i++
				for _, exception := range function.Exceptions {
					PrintField(w, i, exception)
				}
				i--
				fmt.Fprintf(w, "%s)", i)
			}

			// function.Annotations
			io.WriteString(w, "\n\n")
		}
		i--
	}
	fmt.Fprintf(w, "}\n\n")
}

func oneWay(oneWay bool) string {
	if oneWay {
		return "oneway "
	}
	return ""
}

func returnType(returnType ast.Type) string {
	if returnType == nil {
		return "void "
	}
	return fmt.Sprintf("%s ", returnType.String())
}
