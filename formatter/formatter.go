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

func Format(file string, w io.Writer) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	comments := parser.ParseComments(buf)

	tree, err := idl.Parse(buf)
	if err != nil {
		return err
	}

	blocks := []Block{}
	for _, c := range comments {
		blocks = append(blocks, Block{
			c.Info().Line, c,
		})
	}
	for _, h := range tree.Headers {
		blocks = append(blocks, Block{
			h.Info().Line, h,
		})
	}
	for _, d := range tree.Definitions {
		blocks = append(blocks, Block{
			d.Info().Line, d,
		})
	}
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Line < blocks[j].Line
	})

	return traverse(blocks, w)
}

type Block struct {
	Line    int
	Content interface{}
}

func traverse(blocks []Block, w io.Writer) error {
	var i printer.Indent

	for _, block := range blocks {
		switch c := block.Content.(type) {
		case *token.Comment:
			printer.PrintComment(w, i, c)
		case *ast.Include:
			printer.PrintInclude(w, i, c)
		case *ast.Namespace:
			printer.PrintNamespace(w, i, c)
		case *ast.Constant:
			printer.PrintConstant(w, i, c)
		case *ast.Enum:
			printer.PrintEnum(w, i, c)
		case *ast.Struct:
			printer.PrintStruct(w, i, c)
		case *ast.Service:
			printer.PrintService(w, i, c)
		}
	}
	return nil
}
