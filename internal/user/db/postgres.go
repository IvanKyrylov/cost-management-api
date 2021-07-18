package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/IvanKyrylov/cost-management-api/internal/user"
)

var _ user.Storage = &db{}

type db struct {
	storage    *sql.DB
	collection string
	logger     *log.Logger
}

func NewStorage(storage *sql.DB, collection string, logger *log.Logger) user.Storage {
	return &db{
		storage:    storage,
		collection: collection,
		logger:     logger,
	}
}

func (s *db) FindById(ctx context.Context, uuid string) (user user.User, err error) {
}

func (s *db) FindByUsername(ctx context.Context, username string) (user user.User, err error) {

}

func (s *db) FindAll(ctx context.Context, limit, page int64) (users []user.User, err error) {

}
