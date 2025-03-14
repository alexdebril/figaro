package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/alexdebril/figaro/api"
	"github.com/alexdebril/figaro/cmd"
	"github.com/alexdebril/figaro/log"
)

func main() {
	// Parse flags
	f := getFlags()
	// start a slog.Logger
	l := log.Build(os.Stdout, f.GetLogFormat(), f.Debug)

	l.Info("starting application")
	l.Debug("startup params", slog.Any("flags", f), slog.String("address", f.address))

	ctx := context.Background()

	// Create a new server
	routes := initRouting(l)
	opts := []api.MuxOption{
		api.WithAllowOrigin("https://example.com"),
		api.WithAllowHeaders([]string{"X-Header", "Access-Control-Allow-Origin", "Content-Type"}),
	}
	builder := api.NewMuxBuilder(l, opts...)
	mux := builder.Build(routes)

	s := api.NewServer(l, mux, f.address)
	if err := s.Listen(ctx); err != nil {
		l.Error("server returned an error", slog.Any("err", err))
		os.Exit(1)
	}
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

func initRouting(l *slog.Logger) []*api.Route {
	d := NewDefault(l)
	return []*api.Route{
		api.NewRouteFunc(api.Public, "GET /test", d.Get),
		api.NewRouteFunc(api.Public, "POST /test", d.Post),
	}
}

type Default struct {
	logger *slog.Logger
}

func NewDefault(logger *slog.Logger) *Default {
	return &Default{
		logger: logger,
	}
}

func (d *Default) Get(w http.ResponseWriter, r *http.Request) {
	d.logger.Info("GET request received", "path", r.URL.Path)
	w.WriteHeader(http.StatusOK)
}

func (d *Default) Post(w http.ResponseWriter, r *http.Request) {
	d.logger.Info("POST request received", "path", r.URL.Path)
	w.WriteHeader(http.StatusOK)
}
