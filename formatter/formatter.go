package formatter

import (
	"io"
	"io/ioutil"

	"github.com/minodisk/thriftfmt/parser"
	"github.com/minodisk/thriftfmt/printer"
	"github.com/minodisk/thriftfmt/token"
	"go.uber.org/thriftrw/ast"
)

func Format(r io.Reader, w io.Writer) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	lines, err := parser.Parse(buf)
	if err != nil {
		return err
	}

	return Print(lines, w)
}

func Print(ls parser.Lines, w io.Writer) error {
	eol := func(w io.Writer) {
		printer.PrintEOL(w)
	}

	i := printer.NewIndent(0)
	for _, l := range ls {
		for in, node := range l.Nodes {
			if in > 0 {
				printer.PrintSpace(w)
			}
			switch n := node.(type) {
			case *ast.Include:
				printer.PrintInclude(w, i, n)
			case *ast.Namespace:
				printer.PrintNamespace(w, i, n, eol)
			case *ast.Constant:
				printer.PrintConstant(w, i, n)
			case *ast.Enum:
				printer.PrintEnum(w, i, n)
			case *ast.Struct:
				printer.PrintStruct(w, i, n)
			case *ast.Service:
				printer.PrintService(w, i, n)
			case *token.Comment:
				printer.PrintComment(w, i, n)
			}
		}
		printer.PrintEOL(w)
	}
	return nil
}
