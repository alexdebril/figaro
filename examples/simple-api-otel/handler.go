package main

import (
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type defHandler struct {
	logger *slog.Logger
	tracer trace.Tracer
	meter  metric.Meter
}

func newDefHandler(logger *slog.Logger, tracer trace.Tracer, meter metric.Meter) *defHandler {
	return &defHandler{
		logger: logger,
		tracer: tracer,
		meter:  meter,
	}
}

func (d *defHandler) Get(w http.ResponseWriter, r *http.Request) {
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

func (d *defHandler) Post(w http.ResponseWriter, r *http.Request) {
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
