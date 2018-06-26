package printer

import (
	"bytes"
	"testing"

	"github.com/minodisk/thriftfmt/printer"
)

func TestDoc_PrintDoc(t *testing.T) {
	for _, c := range []struct {
		doc    string
		indent *printer.Indent
		want   string
	}{
		{
			doc:    "foo",
			indent: printer.NewIndent(0),
			want: `/**
 * foo
 */
`,
		},
		{
			doc: `foo
bar
baz`,
			indent: printer.NewIndent(0),
			want: `/**
 * foo
 * bar
 * baz
 */
`,
		},
		{
			doc:    "foo",
			indent: printer.NewIndent(3),
			want: `      /**
       * foo
       */
`,
		},
		{
			doc: `foo
bar
baz`,
			indent: printer.NewIndent(3),
			want: `      /**
       * foo
       * bar
       * baz
       */
`,
		},
	} {
		buf := bytes.NewBuffer([]byte{})
		printer.PrintDoc(buf, c.indent, c.doc)
		got := buf.String()
		if got != c.want {
			t.Errorf("\nwant: %s\n got: %s", c.want, got)

		}
	}
}
