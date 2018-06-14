package formatter_test

import (
	"testing"

	"bytes"

	"github.com/minodisk/thriftfmt/formatter"
)

func TestFormatter_CommentPosition(t *testing.T) {
	for _, c := range []struct {
		input string
		want  string
	}{
		{
			`/* foo */
const i64 foo = 1
/* bar */
const i64 bar = 2`,
			`/* foo */
const i64 foo = 1
/* bar */
const i64 bar = 2`,
		},
	} {
		i := bytes.NewBufferString(c.input)
		o := bytes.NewBuffer([]byte{})
		if err := formatter.Format(i, o); err != nil {
			t.Fatal(err)
		}
		got := o.String()
		if got != c.want {
			t.Errorf("\nwant:\n%s\n got:\n%s", c.want, got)
		}
	}
}
