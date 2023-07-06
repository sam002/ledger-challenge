package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"ledger/pkg/deamonizer"
	iserver "ledger/pkg/server"
	"moul.io/chizap"
	"net/http"
)

type Server struct {
	logger  *zap.Logger
	port    int
	host    string
	handler *chi.Mux
}

func NewServer(host string, port int, logger *zap.Logger) (iserver.Server, error) {
	r := chi.NewRouter()
	r.Use(chizap.New(logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))
	r.Use(middleware.Heartbeat("/ping"))
	s := Server{
		handler: r,
		logger:  logger,
		host:    host,
		port:    port,
	}

	return s, nil
}

func (c Server) Run(d *deamonizer.Daemonizer) {
	url := fmt.Sprintf("%s:%d", c.host, c.port)
	server := &http.Server{Addr: url, Handler: c.handler}

	c.logger.Info("Start http server", zap.String("host", url))
	d.AddDaemon(server.ListenAndServe, func() error {
		return server.Shutdown(context.Background())
	},
	)
}

func (c Server) AddHandler(route string, handler http.HandlerFunc) {
	c.handler.Post(route, handler)
}
