package parser

import (
	"bytes"

	"github.com/minodisk/thriftfmt/token"
)

func ParseComments(buf []byte) []*token.Comment {
	var comments []*token.Comment
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
				for ; i < length-1; i++ {
					c := buf[i]
					body.WriteByte(c)
					switch c {
					case '\n':
						line++
					case '*':
						if buf[i+1] == '/' {
							i++
							body.WriteByte(buf[i])
							comments = append(comments, &token.Comment{
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
	return comments
}
