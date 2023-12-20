package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/vitorbgouveia/start-project-go/internal/repositories"
	pkg "github.com/vitorbgouveia/start-project-go/package"
	"github.com/vitorbgouveia/start-project-go/package/dto"
	"go.uber.org/zap"
)

type walletService struct {
	logger *zap.SugaredLogger
	repo   repositories.WalletRepository
	mu     sync.Mutex
}

type WalletService interface {
	CreateWallet(_ context.Context, accID string) error
	Balance(_ context.Context, accID string) (float64, error)
	WalletWithdraw(_ context.Context, _ dto.WalletWithdraw) (dto.WalletUpdateBalance, error)
	WalletDeposit(_ context.Context, _ dto.WalletDeposit) (dto.WalletUpdateBalance, error)
}

type walletUpdateBalanceParams struct {
	accountId  string
	value      float64
	isWithdraw bool
}

const (
	ErrBalanceNotFound     = "fail to get balance, user not found"
	ErrWalletAlreadyExists = "fail to create wallet, already exists"
)

func NewWalletServie(logger *zap.SugaredLogger, repo repositories.WalletRepository) WalletService {
	return &walletService{
		logger: logger,
		repo:   repo,
	}
}

func (s *walletService) Balance(ctx context.Context, accID string) (float64, error) {
	ctxTimeOut, cancel := context.WithTimeout(ctx, MaxDurationFastQueryDB)
	defer cancel()

	balance, err := s.repo.Balance(ctxTimeOut, accID)
	if errors.Is(err, sql.ErrNoRows) {
		s.logger.Warnw(ErrBalanceNotFound, zap.String(pkg.AccountIDKey, accID), zap.Error(err))
		return 0, nil
	}
	return balance, err
}

func (s *walletService) CreateWallet(ctx context.Context, accID string) error {
	ctxTimeOut, cancel := context.WithTimeout(ctx, MaxDurationFastQueryDB)
	defer cancel()

	exists, err := s.repo.HaveWallet(ctxTimeOut, accID)
	if err != nil {
		return err
	}

	if exists {
		s.logger.Warnw(ErrWalletAlreadyExists, zap.String(pkg.AccountIDKey, accID), zap.Error(err))
		return errors.New(ErrWalletAlreadyExists)
	}

	return s.repo.CreateWallet(ctxTimeOut, accID)
}

func (s *walletService) WalletWithdraw(ctx context.Context, params dto.WalletWithdraw) (dto.WalletUpdateBalance, error) {
	ctxTimeOut, cancel := context.WithTimeout(ctx, MaxDurationSlowQueryDB)
	defer cancel()

	return s.walletUpdateBalance(ctxTimeOut, walletUpdateBalanceParams{
		accountId:  params.AccountId,
		value:      params.Value * -1,
		isWithdraw: true,
	})
}

func (s *walletService) WalletDeposit(ctx context.Context, params dto.WalletDeposit) (dto.WalletUpdateBalance, error) {
	ctxTimeOut, cancel := context.WithTimeout(ctx, MaxDurationSlowQueryDB)
	defer cancel()

	return s.walletUpdateBalance(ctxTimeOut, walletUpdateBalanceParams{
		accountId: params.AccountId,
		value:     params.Value,
	})
}

func (s *walletService) walletUpdateBalance(ctx context.Context, params walletUpdateBalanceParams) (dto.WalletUpdateBalance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentBalance, err := s.repo.Balance(ctx, params.accountId)
	if err != nil {
		return dto.WalletUpdateBalance{}, nil
	}

	newBalance := currentBalance + params.value
	if newBalance < 0 && params.isWithdraw {
		return dto.WalletUpdateBalance{},
			fmt.Errorf("fail to update balance, final balance is not valid. finalBalance: %f",
				newBalance)
	}

	if err := s.repo.UpdateBalance(ctx, repositories.WalletUpdateBalanceParams{
		AccountID:  params.accountId,
		NewBalance: newBalance,
	}); err != nil {
		return dto.WalletUpdateBalance{}, err
	}

	return dto.WalletUpdateBalance{
		NewBalance: newBalance,
		OldBalance: currentBalance,
	}, nil
}
