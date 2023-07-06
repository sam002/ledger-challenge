package accounts

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"ledger/pkg/money"
	"net/http"
)

type Account struct {
	ID      pgtype.UUID
	UserId  pgtype.UUID
	Balance decimal.Decimal
	Equity  decimal.Decimal
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
	CreateInvestorAccount(userId pgtype.UUID, currency string, vat string) (*InvestorAccount, error)
	CreateIssuerAccount(userId pgtype.UUID, currency string, company string) (*IssuerAccount, error)
	GetAccountByUserID(userId pgtype.UUID) (*Account, error)
	Deposit(accountId pgtype.UUID, amount money.Money) (*Account, error)
}

type Handlers interface {
	CreateIssuerAccount(w http.ResponseWriter, r *http.Request)
	CreateInvestorAccount(w http.ResponseWriter, r *http.Request)
	GetAccountByUserID(w http.ResponseWriter, r *http.Request)
	DepositAccount(w http.ResponseWriter, r *http.Request)
}
