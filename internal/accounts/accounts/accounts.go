package accounts

import (
	"ledger/internal/accounts/money"
	"ledger/pkg/uuid"
)

type Account struct {
	Id      string
	UserId  string
	Balance money.Money
	Equity  money.Money
}

type IssuerAccount struct {
	Account *Account
	Company string
}

type InvestorAccount struct {
	Account *Account
	Vat     string
}

type Accounts interface {
	CreateInvestorAccount(userId uuid.UUID, currency string, vat string) (InvestorAccount, error)
	CreateIssuerAccount(userId uuid.UUID, currency string, company string) (IssuerAccount, error)
	GetAccountByUserID(userId uuid.UUID) (Account, error)
}
