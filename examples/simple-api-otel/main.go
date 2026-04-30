package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/alexdebril/figaro/api"
	"github.com/alexdebril/figaro/cmd"
	"github.com/alexdebril/figaro/log"
	"github.com/alexdebril/figaro/otel"
)

const appName = "example-with-otel"

func main() {
	// Parse flags
	f := getFlags()
	// start a slog.Logger
	l := log.Build(os.Stdout, f.GetLogFormat(), f.Debug)

	l.Info("starting application")
	l.Debug("startup params", slog.Any("flags", f), slog.String("address", f.address))

	// build a context with custom attributes to define the resource's name and version
	ctx := newContext()
	// build the API's Mux in a factory function ...
	mux := newAPIMux(ctx, l)
	// wrap it with OTEL's tracing Mux ...
	handler, shutdown := otel.NewHTTPHandler(ctx, l, mux, otel.WithTracing(), otel.WithMetrics())
	defer shutdown(ctx)
	// and finally run the server.
	s := api.NewServer(l, handler, f.address)
	if err := s.Listen(ctx); err != nil {
		l.Error("server returned an error", slog.Any("err", err))
		os.Exit(1)
	}
}

func newContext() context.Context {
	ctx := context.Background()
	return context.WithValue(
		context.WithValue(ctx, otel.CtxAppName, "api-otel"),
		otel.CtxVersion, "1.0.0")
}

type flags struct {
	*cmd.Flags
	address string
}

func getFlags() flags {
	coreFlags, fs := cmd.GetCoreFlags()
	flg := flags{
		address: ":8080",
		Flags:   coreFlags,
	}
	fs.StringVar(&flg.address, "http", flg.address, "HTTP address to listen to")
	_ = fs.Parse(os.Args[1:])
	return flg
}
