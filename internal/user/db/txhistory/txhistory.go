package txhistory

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

var _ user.TransactionHistoryStorage = &db{}

type db struct {
	storage *sql.DB
	logger  *log.Logger
}

func NewStorage(storage *sql.DB, logger *log.Logger) user.TransactionHistoryStorage {
	return &db{
		storage: storage,
		logger:  logger,
	}
}

func (s *db) FindById(ctx context.Context, id uint64) (transactionHistory user.TransactionHistory, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = s.storage.QueryRowContext(ctx, "SELECT id, amount, currency, description, done, datetime, wallet_id FROM shortener where id=$1 limit 1;", id).
		Scan(&transactionHistory.Id, &transactionHistory.Amount, &transactionHistory.Currency, &transactionHistory.Done, transactionHistory.Datetime, transactionHistory.WalletId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return transactionHistory, apperror.ErrNotFound
		}
		return transactionHistory, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return transactionHistory, nil
}

func (s *db) FindByWalletId(ctx context.Context, walletId uint64) (transactionHistorys []user.TransactionHistory, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query, err := s.storage.QueryContext(ctx, "SELECT id, amount, currency, description, done, datetime, wallet_id FROM shortener where wallet_id=$1;", walletId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return transactionHistorys, apperror.ErrNotFound
		}
		return transactionHistorys, fmt.Errorf("failed to execute query. error: %w", err)
	}
	defer query.Close()

	for query.Next() {
		var transactionHistory user.TransactionHistory
		err = query.Scan(&transactionHistory.Id, &transactionHistory.Amount, &transactionHistory.Currency, &transactionHistory.Done, transactionHistory.Datetime, transactionHistory.WalletId)
		if err != nil {
			return transactionHistorys, fmt.Errorf("failed to unmarshal error: %w", err)
		}
		transactionHistorys = append(transactionHistorys, transactionHistory)
	}
	if err = query.Err(); err != nil {
		return transactionHistorys, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return transactionHistorys, nil
}

func (s *db) Create(ctx context.Context, amount *big.Float, currency string, description string, done bool, datetime time.Time, walletId uint64) (id uint64, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = s.storage.QueryRowContext(ctx, "INSERT INTO wallets(amount,currency,description,done,datetime,wallet_id) VALUES($1,$2,$3,$4,$5,$6) RETURNING id;", amount, currency, description, done, datetime, walletId).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperror.ErrNotFound
		}
		return 0, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return 0, nil
}
