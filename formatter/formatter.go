package formatter

import (
	"io"
	"io/ioutil"
	"sort"

	"github.com/minodisk/thriftfmt/formatter/printer"
	"github.com/minodisk/thriftfmt/parser"
	"github.com/minodisk/thriftfmt/token"
	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/idl"
)

type Line struct {
	Number int
	Nodes  []interface{}
}

func (l *Line) Append(n interface{}) {
	l.Nodes = append(l.Nodes, n)
}

type Lines []*Line

func (ls *Lines) At(n int) *Line {
	for _, l := range *ls {
		if l.Number == n {
			return l
		}
	}
	l := &Line{n, []interface{}{}}
	*ls = append(*ls, l)
	return l
}

func Format(r io.Reader, w io.Writer) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	comments := parser.ParseComments(buf)

	tree, err := idl.Parse(buf)
	if err != nil {
		return err
	}

	ls := Lines{}
	for _, c := range comments {
		l := ls.At(c.Info().Line)
		l.Append(c)
	}
	for _, h := range tree.Headers {
		l := ls.At(h.Info().Line)
		l.Append(h)
	}
	for _, d := range tree.Definitions {
		l := ls.At(d.Info().Line)
		l.Append(d)
	}
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].Number < ls[j].Number
	})

	return traverse(ls, w)
}

func traverse(ls Lines, w io.Writer) error {
	i := printer.NewIndent(0)
	for _, l := range ls {
		for _, node := range l.Nodes {
			switch n := node.(type) {
			case *token.Comment:
				printer.PrintComment(w, i, n)
			case *ast.Include:
				printer.PrintInclude(w, i, n)
			case *ast.Namespace:
				printer.PrintNamespace(w, i, n)
			case *ast.Constant:
				printer.PrintConstant(w, i, n)
			case *ast.Enum:
				printer.PrintEnum(w, i, n)
			case *ast.Struct:
				printer.PrintStruct(w, i, n)
			case *ast.Service:
				printer.PrintService(w, i, n)
			}
		}
	}
	return nil
}
