package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/alexdebril/figaro/api"
	"go.opentelemetry.io/otel"
)

func newAPIMux(_ context.Context, l *slog.Logger) *http.ServeMux {
	// Create a new server
	routes := initRouting(l)
	opts := []api.MuxOption{
		api.WithAllowOrigin("https://example.com"),
		api.WithAllowHeaders([]string{"X-Header", "Access-Control-Allow-Origin", "Content-Type"}),
	}
	builder := api.NewMuxBuilder(l, opts...)
	return builder.Build(routes)
}

func initRouting(l *slog.Logger) []*api.Route {
	d := newDefHandler(l, otel.Tracer(appName), otel.Meter(appName))
	return []*api.Route{
		api.NewRouteFunc(api.Public, "GET /test", d.Get),
		api.NewRouteFunc(api.Public, "POST /test", d.Post),
	}
}
