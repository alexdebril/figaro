package http

import (
	"log/slog"
	"net/http"
)

type Route struct {
	Pattern string
	Handler http.Handler
}

func NewRouter(logger *slog.Logger, routes ...Route) http.Handler {
	mux := http.NewServeMux()
	for _, route := range routes {
		logger.Info("new route added", "route", route.Pattern)
		mux.Handle(route.Pattern, route.Handler)
	}
	return mux
}