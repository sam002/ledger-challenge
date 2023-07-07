package main

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	modelAccounts "ledger/internal/accounts/accounts/model"
	modelLedger "ledger/internal/ledger/ledger/model"
	"ledger/internal/migration/config"
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
	//zap.IncreaseLevel(zap.DebugLevel)

	cfg := config.GetConfig(logger)

	undo := zap.ReplaceGlobals(logger)
	defer undo()
	//todo make wrapper for gorm logger
	//dbLogger := zapgorm2.New(zap.L())

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy:         nil,
		FullSaveAssociations:   false,
		//Logger:                                   dbLogger,
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		IgnoreRelationshipsWhenMigrating:         false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		TranslateError:                           false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	})
	if err != nil {
		logger.Error("Cannot start migration", zap.Error(err))
		return
	}

	// todo write to file, then use other migration tool
	db.DryRun = true
	if err := db.AutoMigrate(
		&modelAccounts.Account{},
		&modelAccounts.IssuerAccount{},
		&modelAccounts.InvestorAccount{},
		&modelLedger.InvoiceModel{},
		&modelLedger.DealModel{},
	); err != nil {
		logger.Error("Error in check migrate models", zap.Error(err))
		return
	}

	db.DryRun = false
	if err := db.AutoMigrate(
		&modelAccounts.Account{},
		&modelAccounts.IssuerAccount{},
		&modelAccounts.InvestorAccount{},
		&modelLedger.InvoiceModel{},
		&modelLedger.DealModel{},
	); err != nil {
		logger.Error("Error in migrate models", zap.Error(err))
		return
	}
}
