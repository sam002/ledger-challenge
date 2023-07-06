package main

import (
	"fmt"
	"go.uber.org/zap"
	"ledger/internal/accounts/accounts/handler"
	"ledger/internal/accounts/accounts/model"
	"ledger/internal/accounts/config"
	"ledger/pkg/deamonizer"
	"ledger/pkg/server/server"
	"net/http"
)

func main() {
	//logger
	logger, _ := zap.NewProduction()
	defer logger.Info("Exit. Goodbye!")
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Logs not sync: %v\n", err)
		}
	}(logger)

	cfg := config.GetConfig(logger)

	accounts, err := model.NewAccounts(cfg.DSN, logger)
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
		w.Write([]byte("Accounts API"))
	})

	handlers := handler.NewJsonHandler(accounts, logger)
	apiServer.AddHandler("/create-investor-account", handlers.CreateInvestorAccount)
	apiServer.AddHandler("/create-issuer-account", handlers.CreateIssuerAccount)
	apiServer.AddHandler("/get-account", handlers.GetAccountByUserID)
	apiServer.AddHandler("/deposit-account", handlers.DepositAccount)

	d := deamonizer.NewDaemonizer(logger)
	go apiServer.Run(&d)

	//TODO implement transactions

	d.Start()

	d.GracefulShutdown(deamonizer.DEFAULT_TIMEOUT_SHUTDOWN)
}
