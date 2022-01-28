package home

import (
	"context"

	"github.com/r523/dafdardaar/internal/model"
)

// office stores the office model into the database.
type Office interface {
	Set(ctx context.Context, office *model.Office) error
	Get(ctx context.Context, id string) (model.Office, error)
}
