package handler

import (
	"go.uber.org/zap"
	"ledger/internal/ledger/ledger"
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

func (j *JsonHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (j *JsonHandler) CreateBid(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (j *JsonHandler) ApproveInvoice(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (j *JsonHandler) DeclineInvoice(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
