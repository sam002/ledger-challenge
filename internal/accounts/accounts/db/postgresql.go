package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ledger/internal/accounts/accounts"
	"ledger/internal/accounts/money"
	"ledger/pkg/uuid"
	"time"
)

type AccountModel struct {
	accounts.Account
	Id        string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	User      string      `gorm:"uniqueIndex;type:uuid"`
	Balance   money.Money `gorm:"type:money"`
	Equity    money.Money `gorm:"type:money"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type Accounts struct {
	connect *gorm.DB
	logger  *zap.Logger
}

func NewAccounts(dsn string, logger *zap.Logger) (accounts.Accounts, error) {
	c, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Cannot init Users")
		return nil, err
	}

	acc := Accounts{
		connect: c,
		logger:  logger,
	}
	return acc, nil
}

func (a Accounts) CreateInvestorAccount(userId uuid.UUID, currency string, vat string) (accounts.InvestorAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (a Accounts) CreateIssuerAccount(userId uuid.UUID, currency string, company string) (accounts.IssuerAccount, error) {
	acc := AccountModel{
		User: "",
		Balance: money.Money{
			Amount:   0,
			Currency: currency,
			Digits:   2,
		},
		CreatedAt: nil,
		UpdatedAt: nil,
	}

	a.connect.Create(acc)
	res := accounts.IssuerAccount{
		Account: &accounts.Account{
			Id:      acc.Id,
			UserId:  acc.UserId,
			Balance: acc.Balance,
			Equity:  acc.Equity,
		},
		Company: company,
	}
	return res, nil
}

func (a Accounts) GetAccountByUserID(userId uuid.UUID) (accounts.Account, error) {
	//TODO implement me
	panic("implement me")
}
