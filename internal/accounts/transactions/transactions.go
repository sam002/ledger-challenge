package transactions

import (
	"database/sql/driver"
	"github.com/jackc/pgx/v5/pgtype"
	"ledger/internal/accounts/accounts"
)

type Transaction struct {
	ID          pgtype.UUID
	FromAccount pgtype.UUID
	ToAccount   pgtype.UUID
	Status      TransactionStatus
}

type Transactions interface {
	CreateTransaction(from accounts.Account, to accounts.Account) (Transaction, error)
	ApproveTransaction(transaction Transaction) error
	DeclineTransaction(transaction Transaction) error
}

// todo separate
type TransactionStatus string

const (
	NEW               TransactionStatus = "NEW"
	APPROVING_PROCESS TransactionStatus = "APPROVING_PROCESS"
	DECLINING_PROCESS TransactionStatus = "DECLINING_PROCESS"
	APPROVED          TransactionStatus = "APPROVED"
	DECLINED          TransactionStatus = "DECLINED"
)

func (ts *TransactionStatus) Scan(value interface{}) error {
	*ts = TransactionStatus(value.([]byte))
	return nil
}

func (ts *TransactionStatus) Value() (driver.Value, error) {
	return string(*ts), nil
}
