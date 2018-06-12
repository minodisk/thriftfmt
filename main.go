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
	"bytes"
	"sort"
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
		//buf := bytes.Buffer{}
		if err := format(file, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

type Comment struct {
	Line int
	Body string
}

type CommentInfo struct {
	Line int
}

func (c Comment) Info() CommentInfo  {
	return CommentInfo{Line: c.Line}
}

func format(file string, w io.Writer) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	comments := []Comment{}
	line := 1
	length := len(buf)
	for i := 0; i < length; i++ {
		b := buf[i]
		switch b {
		case '\n':
			line++
		case '/':
			if buf[i+1] == '*' && buf[i+2] != '*' {
				l := line
				body := &bytes.Buffer{}
				commentLoop:
				for ; i < length; i++ {
					c := buf[i]
					body.WriteByte(c)
					switch c {
					case '\n':
						line++
					case '*':
						if buf[i+1] == '/' {
							body.WriteByte(buf[i+1])
							comments = append(comments, Comment{
								Line: l,
								Body: body.String(),
							})
							break commentLoop
						}
					}
				}
			}
		}
	}

	tree, err := idl.Parse(buf)
	if err != nil {
		return err
	}

	blocks := []Block{}
	for _, c := range comments {
		blocks = append(blocks,Block{
			c.Info().Line, c,
		})
	}
	for _,h := range tree.Headers {
		blocks = append(blocks,Block{
			h.Info().Line, h,
		})
	}
	for _,d := range tree.Definitions {
		blocks = append(blocks,Block{
			d.Info().Line, d,
		})
	}
	sort.Slice(blocks, func(i,j int)bool {
		return blocks[i].Line < blocks[j].Line
	})

	return traverse(blocks, w)
}

type Block struct {
	Line int
	Content interface{}
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

func traverse(blocks []Block, w io.Writer) error {

	var indent Indent

	for _, block := range blocks{
		switch c := block.Content.(type) {
		case Comment:
			fmt.Fprintf(w, "%s\n\n", c.Body)

		case *ast.Include:
			fmt.Fprintf(w, "%sinclude ", indent)
			if c.Name != "" {
				fmt.Fprintf(w, "%s ", c.Name)
			}
			fmt.Fprintf(w, "\"%s\"\n", c.Path)
		case *ast.Namespace:
			fmt.Fprintf(w, "%snamespace %s %s\n", indent, c.Scope, c.Name)
		case *ast.Constant:
			printDoc(w, indent, c.Doc)
			fmt.Fprintf(w, "%sconst %s %s = %s\n", indent, c.Type, c.Name, constantValue(c.Value))

		case *ast.Enum:
			printDoc(w, indent, c.Doc)
			fmt.Fprintf(w, "%senum %s {", indent, c.Name)
			if len(c.Items) == 0 {
				fmt.Fprintf(w, "}")
			} else {
				fmt.Fprintf(w, "\n")
				indent++
				for _, item := range c.Items {
					//fmt.Println(item.Line)
					fmt.Fprintf(w, "%s%s = %d\n", indent, item.Name, *item.Value)
				}
				indent--
				fmt.Fprintf(w, "%s}\n\n", indent)
			}

		case *ast.Struct:
			printDoc(w, indent, c.Doc)
			fmt.Fprintf(w, "%sstruct %s {", indent, c.Name)
			if len(c.Fields) > 0 {
				fmt.Fprintf(w, "\n")
				indent++
				for _, field := range c.Fields {
					printField(w, indent, field)
				}
				indent--
			}
			fmt.Fprintf(w, "%s}\n\n", indent)

		case *ast.Service:
			printDoc(w, indent, c.Doc)
			fmt.Fprintf(w, "service %s {", c.Name)
			if len(c.Functions) > 0 {
				fmt.Fprint(w, "\n\n")
				indent++
				for _, function := range c.Functions {
					//fmt.Println(function.Line)
					printDoc(w, indent, function.Doc)
					fmt.Fprintf(w, "%s%s%s%s(", indent, oneWay(function.OneWay), returnType(function.ReturnType), function.Name)
					if len(function.Parameters) > 0 {
						fmt.Fprint(w, "\n")
						indent++
						for _, parameter := range function.Parameters {
							//fmt.Println(parameter.Line)
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
