package user

import (
	"context"

	"github.com/r523/dafdardaar/internal/model"
)

// User stores and retrieves users.
type User interface {
	Set(ctx context.Context, user model.User) error
	Get(ctx context.Context, id string) (model.User, error)
}
