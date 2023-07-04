package db

import (
	"ledger/internal/accounts/money"
	"time"
)

type AccountModel struct {
	Id        string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	User      string      `gorm:"uniqueIndex;type:uuid"`
	Balance   money.Money `gorm:"type:money"`
	Equity    money.Money `gorm:"type:money"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
