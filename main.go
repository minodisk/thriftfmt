package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/idl"
)

type options struct {
	display bool
	report  bool
	list    bool
	write   bool
}

func main() {
	if err := _main(); err != nil {
		_, e := os.Stderr.WriteString(fmt.Sprintf("%s", err))
		if e != nil {
			// do nothing
		}
		os.Exit(2)
	}
}

func _main() error {
	_, files := parseFlag()
	for _, file := range files {
		if err := format(file, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

func format(file string, w io.Writer) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	tree, err := idl.Parse(buf)
	if err != nil {
		return err
	}
	return traverse(tree, w)
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

func traverse(tree *ast.Program, w io.Writer) error {

	for _, header := range tree.Headers {
		fmt.Println(header.Info().Line)
	}
	for _, definition := range tree.Definitions {
		fmt.Println(definition.Info().Line)
	}

	var indent Indent

	var isInclude = false
	for _, header := range tree.Headers {
		switch h := header.(type) {
		case *ast.Include:
			if !isInclude {
				fmt.Fprintf(w, "\n")
				isInclude = true
			}
			fmt.Fprintf(w, "%sinclude ", indent)
			if h.Name != "" {
				fmt.Fprintf(w, "%s ", h.Name)
			}
			fmt.Fprintf(w, "\"%s\"\n", h.Path)
		case *ast.Namespace:
			if isInclude {
				fmt.Fprintf(w, "\n")
				isInclude = false
			}
			fmt.Fprintf(w, "%snamespace %s %s\n", indent, h.Scope, h.Name)
		}
	}

	fmt.Fprintf(w, "\n")

	for _, definition := range tree.Definitions {
		//fmt.Println(definition.Info())

		switch d := definition.(type) {
		case *ast.Constant:
			printDoc(w, indent, d.Doc)
			fmt.Fprintf(w, "%sconst %s %s = %s\n", indent, d.Type, d.Name, constantValue(d.Value))

		case *ast.Enum:
			printDoc(w, indent, d.Doc)
			fmt.Fprintf(w, "%senum %s {", indent, d.Name)
			if len(d.Items) == 0 {
				fmt.Fprintf(w, "}")
			} else {
				fmt.Fprintf(w, "\n")
				indent++
				for _, item := range d.Items {
					fmt.Fprintf(w, "%s%s = %d\n", indent, item.Name, *item.Value)
				}
				indent--
				fmt.Fprintf(w, "%s}\n\n", indent)
			}

		case *ast.Struct:
			printDoc(w, indent, d.Doc)
			fmt.Fprintf(w, "%sstruct %s {", indent, d.Name)
			if len(d.Fields) > 0 {
				fmt.Fprintf(w, "\n")
				indent++
				for _, field := range d.Fields {
					printField(w, indent, field)
				}
				indent--
			}
			fmt.Fprintf(w, "%s}\n\n", indent)

		case *ast.Service:
			printDoc(w, indent, d.Doc)
			fmt.Fprintf(w, "service %s {", d.Name)
			if len(d.Functions) > 0 {
				fmt.Fprint(w, "\n\n")
				indent++
				for _, function := range d.Functions {
					printDoc(w, indent, function.Doc)
					fmt.Fprintf(w, "%s%s%s%s(", indent, oneWay(function.OneWay), returnType(function.ReturnType), function.Name)
					if len(function.Parameters) > 0 {
						fmt.Fprint(w, "\n")
						indent++
						for _, parameter := range function.Parameters {
							printField(w, indent, parameter)
						}
						indent--
						fmt.Fprintf(w, "%s)", indent)
					} else {
						fmt.Fprintf(w, ")")
					}

					if len(function.Exceptions) > 0 {
						fmt.Fprintf(w, " throws (\n")
						indent++
						for _, exception := range function.Exceptions {
							printField(w, indent, exception)
						}
						indent--
						fmt.Fprintf(w, "%s)", indent)
					}

					// function.Annotations
					io.WriteString(w, "\n\n")
				}
				indent--
			}
			fmt.Fprintf(w, "}\n\n")
		}
	}
	return nil
}

func printDoc(w io.Writer, indent Indent, doc string) {
	if len(doc) == 0 {
		return
	}
	fmt.Fprintf(w, "%s/**\n", indent)
	for _, line := range strings.Split(doc, "\n") {
		fmt.Fprintf(w, "%s * %s\n", indent, line)
	}
	fmt.Fprintf(w, "%s */\n", indent)
}

func printField(w io.Writer, indent Indent, field *ast.Field) {
	fmt.Fprintf(w, "%s%d: %s%s %s\n", indent, field.ID, requiredness(field.Requiredness), field.Type, field.Name)
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

func requiredness(r ast.Requiredness) string {
	switch r {
	default:
		return ""
	case ast.Required:
		return "required "
	case ast.Optional:
		return "optional "
	}
}

type Indent int

func (indent Indent) String() string {
	return strings.Repeat("  ", int(indent))
}

func parseFlag() (options, []string) {
	display := flag.Bool("d", false, "display diffs instead of rewriting files")
	report := flag.Bool("e", false, "report all errors (not just the first 10 on different lines)")
	list := flag.Bool("l", false, "list files whose formatting differs from gofmt's")
	write := flag.Bool("w", false, "write result to (source) file instead of stdout")
	flag.Parse()
	files := flag.Args()
	return options{
		*display,
		*report,
		*list,
		*write,
	}, files
}
