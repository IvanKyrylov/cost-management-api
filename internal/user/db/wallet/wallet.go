package wallet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/IvanKyrylov/cost-management-api/internal/apperror"
	"github.com/IvanKyrylov/cost-management-api/internal/user"
)

var _ user.WalletStorage = &db{}

type db struct {
	storage *sql.DB
	logger  *log.Logger
}

func NewStorage(storage *sql.DB, logger *log.Logger) user.WalletStorage {
	return &db{
		storage: storage,
		logger:  logger,
	}
}

func (s *db) FindById(ctx context.Context, id uint64) (wallet user.Wallet, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = s.storage.QueryRowContext(ctx, "SELECT id, amount, currency, user_id FROM wallets where id=$1 limit 1;", id).
		Scan(&wallet.Id, &wallet.Amount, &wallet.Currency, &wallet.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wallet, apperror.ErrNotFound
		}
		return wallet, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return wallet, nil
}

func (s *db) FindByUserId(ctx context.Context, userId uint64) (wallets []user.Wallet, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query, err := s.storage.QueryContext(ctx, "SELECT id, amount, currency, user_id FROM wallets where user_id=$1;", userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wallets, apperror.ErrNotFound
		}
		return wallets, fmt.Errorf("failed to execute query. error: %w", err)
	}
	defer query.Close()

	for query.Next() {
		var wallet user.Wallet
		err = query.Scan(&wallet.Id, &wallet.Amount, &wallet.Currency, &wallet.UserId)
		if err != nil {
			return wallets, fmt.Errorf("failed to unmarshal error: %w", err)
		}
		wallets = append(wallets, wallet)
	}
	if err = query.Err(); err != nil {
		return wallets, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return wallets, nil
}

func (s *db) Create(ctx context.Context, amount *big.Float, currency string, userId uint64) (id uint64, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = s.storage.QueryRowContext(ctx, "INSERT INTO wallets(amount,currency,user_id) VALUES($1,$2,$3) RETURNING id;", amount, currency, userId).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperror.ErrNotFound
		}
		return 0, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return 0, nil
}

func (s *db) Update(ctx context.Context, walletId uint64, amount *big.Float) (id uint64, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = s.storage.QueryRowContext(ctx, "UPDATE wallets SET amount=$1 WHERE id=$2 RETURNING id;", amount, walletId).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperror.ErrNotFound
		}
		return 0, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return 0, nil
}
