package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/alexdebril/figaro/api"
	"github.com/alexdebril/figaro/cmd"
	"github.com/alexdebril/figaro/log"
	figaro_otel "github.com/alexdebril/figaro/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const appName = "example-with-otel"

func main() {
	// Parse flags
	f := getFlags()
	// start a slog.Logger
	l := log.Build(os.Stdout, f.GetLogFormat(), f.Debug)

	l.Info("starting application")
	l.Debug("startup params", slog.Any("flags", f), slog.String("address", f.address))

	ctx := context.Background()
	shutdowns, err := figaro_otel.Setup(ctx, figaro_otel.WithMetrics(), figaro_otel.WithTracing())
	if err != nil {
		l.Error("could not start Open Telemetry", "error", err.Error())
	}
	defer func(ctx context.Context) {
		err := shutdowns(ctx)
		if err != nil {
			l.Error("Error shutting down Open Telemetry", "error", err.Error())
		}
	}(ctx)

	// Create a new server
	routes := initRouting(l)
	opts := []api.MuxOption{
		api.WithAllowOrigin("https://example.com"),
		api.WithAllowHeaders([]string{"X-Header", "Access-Control-Allow-Origin", "Content-Type"}),
	}
	builder := api.NewMuxBuilder(l, opts...)
	mux := builder.Build(routes)

	handler := otelhttp.NewHandler(mux, "/")
	s := api.NewServer(l, handler, f.address)
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
	d := NewDefault(l, otel.Tracer(appName), otel.Meter(appName))
	return []*api.Route{
		api.NewRouteFunc(api.Public, "GET /test", d.Get),
		api.NewRouteFunc(api.Public, "POST /test", d.Post),
	}
}

type Default struct {
	logger *slog.Logger
	tracer trace.Tracer
	meter  metric.Meter
}

func NewDefault(logger *slog.Logger, tracer trace.Tracer, meter metric.Meter) *Default {
	return &Default{
		logger: logger,
		tracer: tracer,
		meter:  meter,
	}
}

func (d *Default) Get(w http.ResponseWriter, r *http.Request) {
	ctx, span := d.tracer.Start(r.Context(), "get_test")
	defer span.End()
	d.logger.Info("GET request received", "path", r.URL.Path)
	w.WriteHeader(http.StatusOK)

	cnt, err := d.meter.Int64Counter("get")
	if err != nil {
		d.logger.Error("Error in counting GET request", "error", err)
		return
	}
	cnt.Add(ctx, 1)
}

func (d *Default) Post(w http.ResponseWriter, r *http.Request) {
	ctx, span := d.tracer.Start(r.Context(), "post_test")
	defer span.End()
	d.logger.Info("POST request received", "path", r.URL.Path)
	w.WriteHeader(http.StatusOK)

	cnt, err := d.meter.Int64Counter("post")
	if err != nil {
		d.logger.Error("Error in counting POST request", "error", err)
		return
	}
	cnt.Add(ctx, 1)
}
