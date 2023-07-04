package transfers

import (
	"ledger/internal/accounts/accounts"
	"ledger/pkg/uuid"
)

type Transfer struct {
	Id          uuid.UUID
	FromAccount *accounts.Account
	ToAccount   *accounts.Account
	Status      string //todo ENUM[new, processing approve, processing decline, approved, declined]
}

type Transfers interface {
	CreateTransfer(from *accounts.Account, to *accounts.Account) (Transfer, error)
	ApproveTransfer(transfer Transfer) error
	DeclineTransfer(transfer Transfer) error
}
