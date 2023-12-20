package dto

type WalletWithdraw struct {
	AccountId string
	Value     float64
}

type WalletDeposit struct {
	AccountId string
	Value     float64
}

type WalletUpdateBalance struct {
	NewBalance float64
	OldBalance float64
}
