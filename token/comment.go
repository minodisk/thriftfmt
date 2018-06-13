package token

type Comment struct {
	Line int
	Body string
}

type CommentInfo struct {
	Line int
}

func (c *Comment) Info() CommentInfo {
	return CommentInfo{Line: c.Line}
}
