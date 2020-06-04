package session

import (
	"context"
)

type Session struct {
	ID       string
	UserId   string
	UserName string
}

func SessionFromContext(ctx context.Context) (*Session, bool) {
	session, exist := ctx.Value("sessionKey").(*Session)
	return session, exist
}
