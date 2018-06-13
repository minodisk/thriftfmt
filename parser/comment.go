package parser

import (
	"bytes"
)

type Comment struct {
	Line int
	Body string
}

type CommentInfo struct {
	Line int
}

func (c Comment) Info() CommentInfo {
	return CommentInfo{Line: c.Line}
}

func ParseComments(buf []byte) []Comment {
	comments := []Comment{}
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
				for ; i < length; i++ {
					c := buf[i]
					body.WriteByte(c)
					switch c {
					case '\n':
						line++
					case '*':
						if buf[i+1] == '/' {
							body.WriteByte(buf[i+1])
							comments = append(comments, Comment{
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
