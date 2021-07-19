package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/IvanKyrylov/cost-management-api/internal/apperror"
)

var _ Service = &service{}

type Service interface {
	GetById(ctx context.Context, id uint64) (UserDTO, error)
	GetByUsername(ctx context.Context, userName string) (UserDTO, error)
	GetAll(ctx context.Context, limit, page uint64) ([]UserDTO, error)
}

type service struct {
	userStorage               UserStorage
	walletStorage             WalletStorage
	transactionHistoryStorage TransactionHistoryStorage
	logger                    *log.Logger
}

func NewService(userStorage UserStorage, walletStorage WalletStorage, transactionHistoryStorage TransactionHistoryStorage, logger *log.Logger) (Service, error) {
	return &service{
		userStorage:               userStorage,
		walletStorage:             walletStorage,
		transactionHistoryStorage: transactionHistoryStorage,
		logger:                    logger,
	}, nil
}

func (s service) GetById(ctx context.Context, id uint64) (userDto UserDTO, err error) {
	user, err := s.userStorage.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound.Err) {
			return userDto, err
		}
		return userDto, fmt.Errorf("failed to find user by id. error: %w", err)
	}

	wallets, err := s.walletStorage.FindByUserId(ctx, user.Id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound.Err) {
			return userDto, err
		}
		return userDto, fmt.Errorf("failed to find wallets by user_id. error: %w", err)
	}

	walletsDto := make([]WalletDTO, 0, len(wallets))
	for i := 0; i < len(wallets); i++ {
		walletsDto = append(walletsDto, walletMapping(wallets[i]))
	}
	userDto = userMapping(user)
	userDto.Wallets = walletsDto
	return userDto, nil
}

func (s service) GetByUsername(ctx context.Context, userName string) (userDto UserDTO, err error) {
	user, err := s.userStorage.FindByUsername(ctx, userName)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound.Err) {
			return userDto, err
		}
		return userDto, fmt.Errorf("failed to find user by username. error: %w", err)
	}

	wallets, err := s.walletStorage.FindByUserId(ctx, user.Id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound.Err) {
			return userDto, err
		}
		return userDto, fmt.Errorf("failed to find wallets by user_id. error: %w", err)
	}

	walletsDto := make([]WalletDTO, 0, len(wallets))
	for i := 0; i < len(wallets); i++ {
		walletsDto = append(walletsDto, walletMapping(wallets[i]))
	}
	userDto = userMapping(user)
	userDto.Wallets = walletsDto
	return userDto, nil
}

func (s service) GetAll(ctx context.Context, limit, page uint64) (usersDto []UserDTO, err error) {
	users, err := s.userStorage.FindAll(ctx, limit, page)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound.Err) {
			return usersDto, err
		}
		return usersDto, fmt.Errorf("failed to find user by username. error: %w", err)
	}
	for i := 0; i < len(users); i++ {
		wallets, err := s.walletStorage.FindByUserId(ctx, users[i].Id)
		if err != nil {
			if errors.Is(err, apperror.ErrNotFound.Err) {
				return usersDto, err
			}
			return usersDto, fmt.Errorf("failed to find wallets by user_id. error: %w", err)
		}
		walletsDto := make([]WalletDTO, 0, len(wallets))
		for i := 0; i < len(wallets); i++ {
			walletsDto = append(walletsDto, walletMapping(wallets[i]))
		}
		userDto := userMapping(users[i])
		userDto.Wallets = walletsDto
		usersDto = append(usersDto, userDto)
	}

	return usersDto, nil
}

func walletMapping(wallet Wallet) (walletDto WalletDTO) {
	walletDto.Id = wallet.Id
	walletDto.Amount = wallet.Amount
	walletDto.Currency = wallet.Currency
	return
}

func userMapping(user User) (userDto UserDTO) {
	userDto.Id = user.Id
	userDto.Name = user.Name
	userDto.Surname = user.Surname
	userDto.Username = user.Username
	return
}
