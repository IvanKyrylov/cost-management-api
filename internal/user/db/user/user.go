package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IvanKyrylov/cost-management-api/internal/apperror"
	"github.com/IvanKyrylov/cost-management-api/internal/user"
)

var _ user.UserStorage = &db{}

type db struct {
	storage *sql.DB
	logger  *log.Logger
}

func NewStorage(storage *sql.DB, logger *log.Logger) user.UserStorage {
	return &db{
		storage: storage,
		logger:  logger,
	}
}

func (s *db) FindById(ctx context.Context, id uint64) (user user.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = s.storage.QueryRowContext(ctx, "SELECT id, name, surname, username FROM users where id=$1 limit 1;", id).
		Scan(&user.Id, &user.Name, &user.Surname, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, apperror.ErrNotFound
		}
		return user, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return user, nil
}

func (s *db) FindByUsername(ctx context.Context, username string) (user user.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = s.storage.QueryRowContext(ctx, "SELECT id, name, surname, username FROM users where username=$1 limit 1;", username).
		Scan(&user.Id, &user.Name, &user.Surname, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, apperror.ErrNotFound
		}
		return user, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return user, nil
}

func (s *db) FindAll(ctx context.Context, limit, page uint64) (users []user.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query, err := s.storage.QueryContext(ctx, "SELECT id, name, surname, username FROM users limit $1 offset $2;", limit, page*limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, apperror.ErrNotFound
		}
		return users, fmt.Errorf("failed to execute query. error: %w", err)
	}
	defer query.Close()

	for query.Next() {
		var user user.User
		err = query.Scan(&user.Id, &user.Name, &user.Surname, &user.Username)
		if err != nil {
			return users, fmt.Errorf("failed to unmarshal error: %w", err)
		}
		users = append(users, user)
	}
	if err = query.Err(); err != nil {
		return users, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return users, nil
}
