package handler

import (
	"encoding/json"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"ledger/internal/ledger/ledger"
	"ledger/pkg/money"
	"net/http"
)

var _ ledger.Handlers = &JsonHandler{}

type JsonHandler struct {
	ledger ledger.Ledger
	logger *zap.Logger
}

func NewJsonHandler(ledger ledger.Ledger, logger *zap.Logger) ledger.Handlers {
	return &JsonHandler{ledger: ledger, logger: logger}
}

type CreateInvoiceRequest struct {
	IssuerId    pgtype.UUID     `json:"issuer_id"`
	Description string          `json:"description"`
	Volume      uint            `json:"volume"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
}

type InvoiceResponse struct {
	ID          pgtype.UUID     `json:"id"`
	IssuerId    pgtype.UUID     `json:"issuer_id"`
	Description string          `json:"description"`
	Volume      uint            `json:"volume"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	Status      string          `json:"status"`
}

func (j *JsonHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	req := CreateInvoiceRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		j.logger.Warn("Wrong query", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	invoice, err := j.ledger.CreateInvoice(req.IssuerId, req.Description, req.Volume, money.Money{Amount: req.UnitPrice})
	if err != nil {
		j.logger.Error("Cannot create invoice", zap.Error(err), zap.Any("params", req))
		w.WriteHeader(500)
		return
	}

	resp := InvoiceResponse{
		ID:          invoice.ID,
		IssuerId:    invoice.IssuerId,
		Description: invoice.Description,
		Volume:      invoice.Volume,
		UnitPrice:   invoice.UnitPrice.Amount,
		Status:      invoice.Status,
	}
	res, err := json.Marshal(resp)
	if err != nil {
		j.logger.Info("Error marshal issuer", zap.Error(err))
		return
	}
	_, err = w.Write(res)
	if err != nil {
		j.logger.Info("Interrupt write response", zap.Error(err))
		return
	}
}

type CreateBidRequest struct {
	InvestorId pgtype.UUID `json:"investor_id"`
	InvoiceId  pgtype.UUID `json:"invoice_id"`
	Volume     uint        `json:"volume"`
}
type CreateBidResponse struct {
	Id         pgtype.UUID `json:"id"`
	InvestorId pgtype.UUID `json:"investor_id"`
	InvoiceId  pgtype.UUID `json:"invoice_id"`
	Volume     uint        `json:"volume"`
}

func (j *JsonHandler) CreateBid(w http.ResponseWriter, r *http.Request) {
	req := CreateBidRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		j.logger.Warn("Wrong query", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	bid, err := j.ledger.CreateBid(req.InvestorId, req.InvoiceId, req.Volume)
	if err != nil {
		j.logger.Error("Cannot create bid", zap.Error(err), zap.Any("params", req))
		w.WriteHeader(500)
		return
	}

	resp := CreateBidResponse{
		Id:         bid.ID,
		InvestorId: bid.InvestorId,
		InvoiceId:  bid.InvoiceId,
		Volume:     bid.Volume,
	}
	res, err := json.Marshal(resp)
	if err != nil {
		j.logger.Info("Error marshal issuer", zap.Error(err))
		return
	}
	_, err = w.Write(res)
	if err != nil {
		j.logger.Info("Interrupt write response", zap.Error(err))
		return
	}
}

func (j *JsonHandler) ApproveInvoice(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (j *JsonHandler) DeclineInvoice(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
