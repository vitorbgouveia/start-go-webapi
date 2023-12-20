package repositories

import (
	"context"
	"database/sql"
	"time"
)

const (
	ErrMakePrepareStmt = "could not make stmt in database"
)

type WalletUpdateBalanceParams struct {
	AccountID  string
	NewBalance float64
}

type walletRepository struct {
	dbConn *sql.DB
}

type WalletRepository interface {
	HaveWallet(ctx context.Context, accID string) (bool, error)
	CreateWallet(ctx context.Context, accID string) error
	Balance(ctx context.Context, accID string) (float64, error)
	UpdateBalance(_ context.Context, _ WalletUpdateBalanceParams) error
}

func NewWalletRepository(dbConn *sql.DB) WalletRepository {
	return &walletRepository{
		dbConn: dbConn,
	}
}

func (s *walletRepository) HaveWallet(ctx context.Context, accID string) (bool, error) {
	var qtd int
	err := s.dbConn.QueryRowContext(ctx, `
		SELECT count(1) FROM banking.wallet w WHERE w.account_id = $1`, accID).Scan(&qtd)
	if err != nil {
		return false, err
	}

	return qtd > 0, err
}

func (s *walletRepository) Balance(ctx context.Context, accID string) (float64, error) {
	var balance *float64
	err := s.dbConn.QueryRowContext(ctx, `
		SELECT w.balance FROM banking.wallet w WHERE w.account_id = $1`, accID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	if balance == nil {
		return 0, err
	}

	return *balance, nil
}

func (s *walletRepository) CreateWallet(ctx context.Context, accID string) error {
	if _, err := s.dbConn.ExecContext(ctx,
		"INSERT INTO banking.wallet (account_id, created) VALUES ($1, $2)",
		accID, time.Now()); err != nil {
		return err
	}

	return nil
}

func (s *walletRepository) UpdateBalance(ctx context.Context, params WalletUpdateBalanceParams) error {
	if _, err := s.dbConn.ExecContext(ctx,
		"UPDATE banking.wallet w SET balance=$1 WHERE w.account_id=$2",
		params.NewBalance, params.AccountID); err != nil {
		return err
	}

	return nil
}
