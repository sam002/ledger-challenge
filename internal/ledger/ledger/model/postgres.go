package model

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
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

type InvoiceModel struct {
	gorm.Model
	ID          pgtype.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	IssuerId    pgtype.UUID `gorm:"index;type:uuid"`
	Description string
	Volume      uint
	UnitPrice   decimal.Decimal `gorm:"type:decimal"`
	Status      string
}

func (InvoiceModel) TableName() string {
	return "invoice"
}

func (l *LedgerDB) CreateInvoice(issuerId pgtype.UUID, description string, volume uint, unitPrice money.Money) (*ledger.Invoice, error) {
	inv := InvoiceModel{
		IssuerId:    issuerId,
		Description: description,
		Volume:      volume,
		UnitPrice:   unitPrice.Amount,
		Status:      "", //todo implement statuses
	}
	l.connect.Create(&inv)
	if l.connect.Error != nil {
		return nil, l.connect.Error
	}
	res := ledger.Invoice{
		ID:          inv.ID,
		IssuerId:    inv.IssuerId,
		Description: inv.Description,
		Volume:      inv.Volume,
		UnitPrice:   money.Money{Amount: inv.UnitPrice},
		Status:      inv.Status,
	}
	return &res, nil
}

type DealModel struct {
	gorm.Model
	ID            pgtype.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	InvestorId    pgtype.UUID `gorm:"index;type:uuid"`
	TransactionId pgtype.UUID `gorm:"index;type:uuid"`
	InvoiceId     pgtype.UUID
	Invoice       *InvoiceModel
	Volume        uint
}

func (DealModel) TableName() string {
	return "deal"
}

func (l *LedgerDB) CreateBid(investorId pgtype.UUID, invoiceId pgtype.UUID, volume uint) (*ledger.Deal, error) {
	tx := l.connect.Begin()
	inv, err := l.getInvoiceForNewDeal(invoiceId, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tv, err := l.getTradedVolumeByInvoice(inv.ID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if inv.Volume-tv < volume {
		tx.Rollback()
		err := fmt.Errorf("volume of the bid abouve availible")
		l.logger.Warn(err.Error(), zap.Any("invoice", inv), zap.Uint("traded volume", tv), zap.Uint("bid volume", volume))
		return nil, fmt.Errorf("volume of the bid above availible")
	}
	//todo made transaction

	d := DealModel{
		InvestorId: investorId,
		Invoice:    inv,
		Volume:     volume,
	}
	l.connect.Create(&d)
	if l.connect.Error != nil {
		tx.Rollback()
		return nil, l.connect.Error
	}
	res := ledger.Deal{
		ID:         d.ID,
		InvestorId: d.InvestorId,
		InvoiceId:  d.InvoiceId,
		Volume:     d.Volume,
	}

	tx.Commit()
	return &res, nil
}

func (l *LedgerDB) ApproveInvoice(issueId pgtype.UUID) (*ledger.Invoice, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LedgerDB) DeclineInvoice(issueId pgtype.UUID, amount money.Money) (*ledger.Invoice, error) {
	//TODO implement me
	panic("implement me")
}

// required transaction
func (l *LedgerDB) getInvoiceForNewDeal(id pgtype.UUID, tx *gorm.DB) (*InvoiceModel, error) {
	inv := InvoiceModel{ID: id}
	tx.First(&inv)
	// todo add statuses
	//if inv.Status != "NEW" {
	//	return nil, fmt.Errorf("invoice not accept a new bid")
	//}

	return &inv, nil
}

func (l *LedgerDB) getTradedVolumeByInvoice(invoiceId pgtype.UUID, tx *gorm.DB) (uint, error) {
	res := struct {
		traded uint
	}{}

	tx.Select("SUM(volume) AS traded").
		Where("invoice_id=?", invoiceId).
		Group("invoice_id").
		Scan(&res)

	return res.traded, nil
}
