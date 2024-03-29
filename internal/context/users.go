package context

import (
	"context"

	"github.com/ethanjmarchand/learnhtmx/internal/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, contact *models.Contact) context.Context {
	return context.WithValue(ctx, userKey, contact)
}

func User(ctx context.Context) *models.Contact {
	val := ctx.Value(userKey)
	user, ok := val.(*models.Contact)
	if !ok {
		// The most likely case is that nothing was ever stored in the context, so it doesn't have a type of *models.User. It is also possible that other code in this package wrote an invalid value using the user key.
		return nil
	}
	return user
}
