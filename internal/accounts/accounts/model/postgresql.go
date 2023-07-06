package model

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ledger/internal/accounts/accounts"
	"ledger/pkg/money"
)

var _ accounts.Accounts = &AccountsDB{}

type AccountsDB struct {
	connect *gorm.DB
	logger  *zap.Logger
}

type Account struct {
	gorm.Model
	ID      pgtype.UUID     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserId  pgtype.UUID     `gorm:"index;type:uuid"`
	Balance decimal.Decimal `gorm:"type:decimal"`
	Equity  decimal.Decimal `gorm:"type:decimal"`
}

type IssuerAccount struct {
	ID        pgtype.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Company   string
	AccountID pgtype.UUID
	Account   Account
}

type InvestorAccount struct {
	ID        pgtype.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	VAT       string      `gorm:"type:varchar(20)"`
	AccountID pgtype.UUID
	Account   Account
}

func NewAccounts(dsn string, logger *zap.Logger) (accounts.Accounts, error) {
	//todo add schemas
	c, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Cannot init Users")
		return nil, err
	}

	acc := AccountsDB{
		connect: c,
		logger:  logger,
	}
	return &acc, nil
}

func (a *AccountsDB) CreateInvestorAccount(userId pgtype.UUID, currency string, vat string) (*accounts.InvestorAccount, error) {
	acc := InvestorAccount{
		VAT: vat,
		Account: Account{
			UserId: userId,
		},
	}
	a.connect.Create(&acc)
	if a.connect.Error != nil {
		return nil, a.connect.Error
	}
	res := accounts.InvestorAccount{
		Account: &accounts.Account{
			ID:      acc.AccountID,
			UserId:  acc.Account.UserId,
			Balance: acc.Account.Balance,
			Equity:  acc.Account.Equity,
		},
		Vat: vat,
	}
	return &res, nil
}

func (a *AccountsDB) CreateIssuerAccount(userId pgtype.UUID, currency string, company string) (*accounts.IssuerAccount, error) {
	acc := IssuerAccount{
		Company: company,
		Account: Account{
			UserId: userId,
		},
	}
	a.connect.Create(&acc)
	if a.connect.Error != nil {
		return nil, a.connect.Error
	}
	res := accounts.IssuerAccount{
		Account: &accounts.Account{
			ID:      acc.AccountID,
			UserId:  acc.Account.UserId,
			Balance: acc.Account.Balance,
			Equity:  acc.Account.Equity,
		},
		Company: company,
	}
	return &res, nil
}

func (a *AccountsDB) GetAccountByUserID(userId pgtype.UUID) (*accounts.Account, error) {
	return &accounts.Account{}, nil
}

func (a *AccountsDB) Deposit(accountId pgtype.UUID, amount money.Money) (*accounts.Account, error) {
	if amount.Amount.IsNegative() {
		e := fmt.Errorf("deposit requred only positive amount")
		a.logger.Warn(e.Error())
		return nil, e
	}
	tx := a.connect.Begin()

	acc := Account{ID: accountId}
	a.connect.First(&acc)
	if a.connect.Error != nil {
		tx.Rollback()
		return nil, a.connect.Error
	}
	//todo check money compatibility
	acc.Balance = amount.Amount.Add(acc.Balance)
	acc.Equity = amount.Amount.Add(acc.Equity)
	a.connect.Save(&acc)

	tx.Commit()

	res := &accounts.Account{
		ID:      acc.ID,
		UserId:  acc.UserId,
		Balance: acc.Balance,
		Equity:  acc.Equity,
	}
	return res, nil
}
