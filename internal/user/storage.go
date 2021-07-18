package user

import (
	"context"
)

type Storage interface {
	FindById(ctx context.Context, id string) (User, error)
	FindByUsername(ctx context.Context, username string) (User, error)
	FindAll(ctx context.Context, limit, page int64) ([]User, error)
}
