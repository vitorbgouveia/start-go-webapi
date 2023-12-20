package web

import "errors"

type InputJsonRequest interface {
	Validate() error
}

type InputCreateWallet struct {
	AccountId string `json:"account_id"`
}

func (s *InputCreateWallet) Validate() error {
	switch {
	case s.AccountId == "":
		return errors.New("account_id is empty")
	}

	return nil
}

type InputWithdrawWallet struct {
	AccountId string  `json:"account_id"`
	Value     float64 `json:"value"`
}

func (s *InputWithdrawWallet) Validate() error {
	switch {
	case s.AccountId == "":
		return errors.New("account_id is empty")
	case s.Value <= 0:
		return errors.New("invalid value to perform withdraw")
	}

	return nil
}

type InputDepositWallet struct {
	AccountId string  `json:"account_id"`
	Value     float64 `json:"value"`
}

func (s *InputDepositWallet) Validate() error {
	switch {
	case s.AccountId == "":
		return errors.New("account_id is empty")
	case s.Value <= 0:
		return errors.New("invalid value to perform deposit")
	}

	return nil
}

type OutputUpdateWalletBalance struct {
	NewBalance float64 `json:"new_balance"`
	OldBalance float64 `json:"old_balance"`
}

type OutputBalanceWalletUser struct {
	Balance float64 `json:"balance"`
}
