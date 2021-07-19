package user

import (
	"math/big"
	"time"
)

type User struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
}

type Wallet struct {
	Id       uint64     `json:"id"`
	Amount   *big.Float `json:"amount"`
	Currency string     `json:"currency"`
	UserId   uint64     `json:"user_id"`
}

type TransactionHistory struct {
	Id          uint64     `json:"id"`
	Amount      *big.Float `json:"amount"`
	Currency    string     `json:"currency"`
	Description string     `json:"description"`
	Done        bool       `json:"done"`
	Datetime    time.Time  `json:"datetime"`
	WalletId    uint64     `json:"wallet_id"`
}

type UserDTO struct {
	Id       uint64
	Name     string
	Surname  string
	Username string
	Wallets  []WalletDTO
}

type WalletDTO struct {
	Id       uint64
	Amount   *big.Float
	Currency string
}
