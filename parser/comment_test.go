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
	} {
		got := parser.ParseComments([]byte(c.input))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("\nwant: %v\n got: %v", c.want, got)
		}
	}
}
