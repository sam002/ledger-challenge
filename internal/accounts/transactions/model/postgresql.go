package model

import (
	"gorm.io/gorm"
	"ledger/internal/accounts/transactions"
)

type TransactionsDB struct {
	gorm.Model
	transactions.Transaction
}
