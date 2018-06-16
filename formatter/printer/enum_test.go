package printer_test

import (
	"testing"

	"bytes"

	"github.com/minodisk/thriftfmt/formatter/printer"
	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/idl"
)

func TestEnum_PrintEnum(t *testing.T) {
	for _, c := range []struct {
		name   string
		input  string
		indent *printer.Indent
		want   string
	}{
		{
			`empty`,
			`enum TweetType {}`,
			printer.NewIndent(0),
			`enum TweetType {}
`,
		},
		{
			`values`,
			`enum TweetType { TWEET, RETWEET = 2, DM = 5, REPLY }`,
			printer.NewIndent(0),
			`enum TweetType {
  TWEET,
  RETWEET = 2,
  DM = 5,
  REPLY
}
`,
		},
		{
			`annotations`,
			`enum TweetType { TWEET } (foo="foo",bar="bar")`,
			printer.NewIndent(0),
			`enum TweetType {
  TWEET
} (foo = "foo", bar = "bar")
`,
		},
		{
			`item annotations`,
			`enum TweetType { TWEET (foo="foo"), RETWEET=2 (baz="baz") }`,
			printer.NewIndent(0),
			`enum TweetType {
  TWEET (foo = "foo"),
  RETWEET = 2 (baz = "baz")
}
`,
		},
		{
			`values`,
			`enum TweetType{TWEET(foo="foo"),RETWEET=2,DM=5(baz="baz"),REPLY}(quz="qux")`,
			printer.NewIndent(0),
			`enum TweetType {
  TWEET (foo = "foo"),
  RETWEET = 2,
  DM = 5 (baz = "baz"),
  REPLY
} (quz = "qux")
`,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			p, err := idl.Parse([]byte(c.input))
			if err != nil {
				t.Fatal(err)
			}
			buf := bytes.NewBuffer([]byte{})
			printer.PrintEnum(buf, c.indent, p.Definitions[0].(*ast.Enum))
			got := buf.String()
			if got != c.want {
				t.Errorf("\nWANT:\n%s\nGOT:\n%s\n", c.want, got)
			}
		})
	}
}
