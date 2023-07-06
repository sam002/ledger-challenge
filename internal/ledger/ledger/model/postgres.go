package model

import (
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ledger/internal/ledger/ledger"
	"ledger/pkg/money"
)

var _ ledger.Ledger = &LedgerDB{}

type LedgerDB struct {
	connect *gorm.DB
	logger  *zap.Logger
}

func NewLedger(dsn string, logger *zap.Logger) (ledger.Ledger, error) {
	c, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Cannot init Users")
		return nil, err
	}

	return &LedgerDB{connect: c, logger: logger}, nil
}

func (l *LedgerDB) CreateInvoice(issuerId pgtype.UUID, description string, volume uint, unitPrice money.Money) (*ledger.Invoice, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LedgerDB) CreateBid(investorId pgtype.UUID, invoiceId pgtype.UUID, volume uint) (*ledger.Deal, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LedgerDB) ApproveInvoice(issueId pgtype.UUID) (*ledger.Invoice, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LedgerDB) DeclineInvoice(issueId pgtype.UUID, amount money.Money) (*ledger.Invoice, error) {
	//TODO implement me
	panic("implement me")
}
