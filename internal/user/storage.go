package user

import (
	"context"
	"math/big"
	"time"
)

type UserStorage interface {
	FindById(ctx context.Context, id uint64) (User, error)
	FindByUsername(ctx context.Context, username string) (User, error)
	FindAll(ctx context.Context, limit, page uint64) ([]User, error)
}

type WalletStorage interface {
	FindById(ctx context.Context, id uint64) (Wallet, error)
	FindByUserId(ctx context.Context, userId uint64) ([]Wallet, error)

	Create(ctx context.Context, amount *big.Float, currency string, userId uint64) (uint64, error)
	Update(ctx context.Context, walletId uint64, amount *big.Float) (uint64, error)
}

type TransactionHistoryStorage interface {
	FindById(ctx context.Context, id uint64) (TransactionHistory, error)
	FindByWalletId(ctx context.Context, walletId uint64) ([]TransactionHistory, error)
	Create(ctx context.Context, amount *big.Float, currency string, description string, done bool, datetime time.Time, walletId uint64) (uint64, error)
}
