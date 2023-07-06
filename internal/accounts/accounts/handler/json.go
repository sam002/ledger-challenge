package handler

import (
	"encoding/json"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"ledger/internal/accounts/accounts"
	"ledger/pkg/money"
	"net/http"
)

var _ accounts.Handlers = &JsonHandler{}

type JsonHandler struct {
	accounts accounts.Accounts
	logger   *zap.Logger
}

func NewJsonHandler(accounts accounts.Accounts, logger *zap.Logger) *JsonHandler {
	return &JsonHandler{accounts: accounts, logger: logger}
}

type ByUserUuidRequest struct {
	Uuid pgtype.UUID `json:"user_id"`
}

type AccountResponse struct {
	ID      pgtype.UUID     `json:"id"`
	UserId  pgtype.UUID     `json:"user_id"`
	Balance decimal.Decimal `json:"balance"`
	Equity  decimal.Decimal `json:"equity"`
}

type IssuerAccountResponse struct {
	ID      pgtype.UUID     `json:"id"`
	UserId  pgtype.UUID     `json:"user_id"`
	Balance decimal.Decimal `json:"balance"`
	Equity  decimal.Decimal `json:"equity"`
	Company string          `json:"company"`
}

type InvestorAccountResponse struct {
	ID      pgtype.UUID     `json:"id"`
	UserId  pgtype.UUID     `json:"user_id"`
	Balance decimal.Decimal `json:"balance"`
	Equity  decimal.Decimal `json:"equity"`
	Vat     string          `json:"vat"`
}

func (j *JsonHandler) CreateIssuerAccount(w http.ResponseWriter, r *http.Request) {
	req := ByUserUuidRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		j.logger.Warn("Wrong query", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	issuer, err := j.accounts.CreateIssuerAccount(req.Uuid, "", "")
	if err != nil {
		j.logger.Error("Cannot create account", zap.Error(err), zap.Any("params", req))
		w.WriteHeader(500)
		return
	}

	issuerAcc := IssuerAccountResponse{
		ID:      issuer.Account.ID,
		UserId:  issuer.Account.UserId,
		Balance: issuer.Account.Balance,
		Equity:  issuer.Account.Equity,
		Company: issuer.Company,
	}
	res, err := json.Marshal(issuerAcc)
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

func (j *JsonHandler) CreateInvestorAccount(w http.ResponseWriter, r *http.Request) {
	req := ByUserUuidRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		j.logger.Warn("Wrong query", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	investor, err := j.accounts.CreateInvestorAccount(req.Uuid, "", "")
	if err != nil {
		j.logger.Error("Cannot create account", zap.Error(err), zap.Any("params", req))
		w.WriteHeader(500)
		return
	}

	investorAcc := InvestorAccountResponse{
		ID:      investor.Account.ID,
		UserId:  investor.Account.UserId,
		Balance: investor.Account.Balance,
		Equity:  investor.Account.Equity,
		Vat:     investor.Vat,
	}
	res, err := json.Marshal(investorAcc)
	if err != nil {
		j.logger.Info("Error marshal investor", zap.Error(err))
		return
	}
	_, err = w.Write(res)
	if err != nil {
		j.logger.Info("Interrupt write response", zap.Error(err))
		return
	}
}

func (j *JsonHandler) GetAccountByUserID(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

type DepositRequest struct {
	Uuid   pgtype.UUID     `json:"account_id"`
	Amount decimal.Decimal `json:"amount"`
}

func (j *JsonHandler) DepositAccount(w http.ResponseWriter, r *http.Request) {
	req := DepositRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		j.logger.Warn("Wrong query", zap.Error(err))
		w.WriteHeader(500)
		return
	}
	acc, err := j.accounts.Deposit(req.Uuid, money.Money{Amount: req.Amount})
	if err != nil {
		w.WriteHeader(500)
		j.logger.Error("Cannot deposit account", zap.Error(err), zap.Any("params", req))
		return
	}

	a := AccountResponse{
		ID:      acc.ID,
		UserId:  acc.UserId,
		Balance: acc.Balance,
		Equity:  acc.Equity,
	}
	res, err := json.Marshal(a)
	if err != nil {
		j.logger.Info("Error marshal investor", zap.Error(err))
		return
	}

	_, err = w.Write(res)
	if err != nil {
		j.logger.Info("Interrupt write response", zap.Error(err))
		return
	}
}
