package authctx

import (
	"context"
	"errors"
)

type ctxKey int

const (
	userContextKey ctxKey = iota
)

type User struct {
	UserID string
	Email  string
}

func UserFromCtx(ctx context.Context) (User, error) {
	u, ok := ctx.Value(userContextKey).(User)
	if !ok {
		return User{}, errors.New("user not found in context")
	}
	return u, nil
}

func ContextWithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}
