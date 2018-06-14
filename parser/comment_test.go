package parser_test

import (
	"reflect"
	"testing"

	"github.com/minodisk/thriftfmt/parser"
	"github.com/minodisk/thriftfmt/token"
)

func TestComment_ParseComment(t *testing.T) {
	for _, c := range []struct {
		input string
		want  []*token.Comment
	}{
		{
			`/* foo */`,
			[]*token.Comment{
				{
					1,
					`/* foo */`,
				},
			},
		},
		{
			`/* foo *//* bar */`,
			[]*token.Comment{
				{
					1,
					`/* foo */`,
				},
				{
					1,
					`/* bar */`,
				},
			},
		},
		{
			`/* foo */
/* bar */`,
			[]*token.Comment{
				{
					1,
					`/* foo */`,
				},
				{
					2,
					`/* bar */`,
				},
			},
		},
	} {
		got := parser.ParseComments([]byte(c.input))
		if !reflect.DeepEqual(got, c.want) {
			t.Error("want:")
			for _, c := range c.want {
				t.Error(*c)
			}
			t.Error(" got:")
			for _, c := range got {
				t.Error(*c)
			}
		}
	}
}
