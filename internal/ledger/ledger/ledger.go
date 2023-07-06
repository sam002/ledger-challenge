package ledger

import (
	"github.com/jackc/pgx/v5/pgtype"
	"ledger/pkg/money"
	"net/http"
)

type Invoice struct {
	ID          pgtype.UUID
	IssuerId    pgtype.UUID
	Description string
	Volume      uint
	UnitPrice   money.Money
	Status      string //todo implement ENUM
}

type Deal struct {
	ID            pgtype.UUID
	InvestorId    pgtype.UUID
	InvoiceId     pgtype.UUID
	Quantity      uint
	TransactionId pgtype.UUID
}

type Ledger interface {
	CreateInvoice(issuerId pgtype.UUID, description string, volume uint, unitPrice money.Money) (*Invoice, error)
	CreateBid(investorId pgtype.UUID, invoiceId pgtype.UUID, volume uint) (*Deal, error)
	ApproveInvoice(issueId pgtype.UUID) (*Invoice, error)
	DeclineInvoice(issueId pgtype.UUID, amount money.Money) (*Invoice, error)
}

type Handlers interface {
	CreateInvoice(w http.ResponseWriter, r *http.Request)
	CreateBid(w http.ResponseWriter, r *http.Request)
	ApproveInvoice(w http.ResponseWriter, r *http.Request)
	DeclineInvoice(w http.ResponseWriter, r *http.Request)
}
