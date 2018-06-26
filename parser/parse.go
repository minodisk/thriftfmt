package parser

import (
	"fmt"
	"sort"

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

//func (ls Lines) At(n int) *Line {
//}

func (ls *Lines) Get(n int) *Line {
	for _, l := range *ls {
		if l.Number == n {
			return l
		}
	}
	l := &Line{n, []interface{}{}}
	*ls = append(*ls, l)
	return l
}

func Parse(buf []byte) (Lines, error) {
	lines := Lines{}
	tree, err := idl.Parse(buf)
	if err != nil {
		return nil, err
	}
	for _, node := range tree.Headers {
		fmt.Printf("%+v\n", node)
		lines.Get(node.Info().Line).Append(node)
	}
	for _, node := range tree.Definitions {
		lines.Get(node.Info().Line).Append(node)
	}
	sort.Slice(lines, func(i, j int) bool {
		return lines[i].Number < lines[j].Number
	})

	comments := ParseComments(buf)
	for _, node := range comments {
		lines.Get(node.Info().Line).Append(node)
	}

	return lines, nil
}
