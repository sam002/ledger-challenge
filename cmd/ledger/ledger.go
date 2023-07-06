package main

import (
	"fmt"
	"go.uber.org/zap"
	"ledger/internal/ledger/config"
	"ledger/internal/ledger/ledger/handler"
	"ledger/internal/ledger/ledger/model"
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

	ledger, err := model.NewLedger(cfg.DSN, logger)
	if err != nil {
		logger.Error("Cannot init ledger service", zap.Error(err))
		return
	}

	handlers := handler.NewJsonHandler(ledger, logger)

	apiServer, err := server.NewServer(cfg.Host, cfg.Port, logger)
	if err != nil {
		logger.Error("Cannot init HTTP API", zap.Error(err))
		return
	}

	apiServer.AddHandler("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ledger API"))
	})

	apiServer.AddHandler("/create-issue", handlers.CreateInvoice)
	apiServer.AddHandler("/create-bid", handlers.CreateBid)
	apiServer.AddHandler("/approve-issue", handlers.ApproveInvoice)
	apiServer.AddHandler("/decline-issue", handlers.DeclineInvoice)

	d := deamonizer.NewDaemonizer(logger)
	go apiServer.Run(&d)

	//todo implement observer for invoices

	d.Start()

	d.GracefulShutdown(deamonizer.DEFAULT_TIMEOUT_SHUTDOWN)
}
