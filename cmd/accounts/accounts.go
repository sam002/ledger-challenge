package main

import (
	"fmt"
	"ledger/cmd/accounts/config"
	"ledger/internal/accounts/users/db"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	//logger
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Logs not sync: %v", err)
		}
	}(logger)

	cfg := config.GetConfig(logger)

	users, err := db.NewUsers(cfg.DSN, logger)

	if err != nil {
		logger.Error("Cannot init users", zap.Error(err))
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		u := users.GetByEmail("test@sam002.net")

		w.Write([]byte(fmt.Sprintf("%v", u)))
	})
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)

	//run http server
}
