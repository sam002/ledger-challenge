package server

import (
	"ledger/pkg/deamonizer"
	"net/http"
)

// todo move to /pkg
type Server interface {
	Run(d *deamonizer.Daemonizer)
	AddHandler(route string, handler http.HandlerFunc)
}
