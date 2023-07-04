package main

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"ledger/internal/accounts/accounts/db"
	"ledger/internal/accounts/config"
	"ledger/internal/accounts/server/server"
	"ledger/pkg/deamonizer"
	"ledger/pkg/uuid"
	"net/http"
)

func main() {
	//logger
	logger, _ := zap.NewProduction()
	defer logger.Info("Exit. Goodbye!")
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Logs not sync: %v\n", err)
		}
	}(logger)

	cfg := config.GetConfig(logger)

	accounts, err := db.NewAccounts(cfg.DSN, logger)
	if err != nil {
		logger.Error("Cannot init Accounts", zap.Error(err))
		return
	}

	apiServer, err := server.NewServer(cfg.Host, cfg.Port, logger)
	if err != nil {
		logger.Error("Cannot init HTTP API", zap.Error(err))
		return
	}

	apiServer.AddHandler("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Articles API"))
	})

	apiServer.AddHandler("/create-issuer-account", func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			UserID   uuid.UUID `json:"UUID"`
			Currency string    `json:"currency"`
			Company  string    `json:"company"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Info("Wrong query", zap.Error(err))
		}
		issuer, err := accounts.CreateIssuerAccount(req.UserID, req.Currency, req.Company)
		if err != nil {
			logger.Error("Cannot create account", zap.Error(err), zap.Any("params", req))
		}
		res, err := json.Marshal(issuer)
		if err != nil {
			logger.Info("Error marshal issuer", zap.Error(err))
		}
		_, err = w.Write(res)
		if err != nil {
			logger.Info("Interrupt write response", zap.Error(err))
		}
	})

	d := deamonizer.NewDaemonizer(logger)
	go apiServer.Run(&d)

	d.Start()

	d.GracefulShutdown(deamonizer.DEFAULT_TIMEOUT_SHUTDOWN)
}
