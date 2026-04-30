package otel

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

const (
	CtxAppName    = "app_name"
	CtxVersion    = "version"
	CtxInstanceID = "instance_id"
)

// taken from: https://opentelemetry.io/docs/languages/go/getting-started/

type Opt func(ctx context.Context, res *resource.Resource, fns []func(context.Context) error) error

func WithTracing() Opt {
	return func(ctx context.Context, res *resource.Resource, fns []func(context.Context) error) error {
		tracerProvider, err := newTracerProvider(ctx, res)
		if err != nil {
			return err
		}
		fns = append(fns, tracerProvider.Shutdown)
		otel.SetTracerProvider(tracerProvider)
		return nil
	}
}

func WithMetrics() Opt {
	return func(ctx context.Context, res *resource.Resource, fns []func(context.Context) error) error {
		meterProvider, err := newMeterProvider(ctx, res)
		if err != nil {
			return err
		}
		fns = append(fns, meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)
		return nil
	}
}

func NewHTTPHandler(ctx context.Context, l *slog.Logger, h http.Handler, opts ...Opt) (http.Handler, func(context.Context)) {
	shutdowns, err := Setup(ctx, opts...)
	if err != nil {
		l.Error("could not start Open Telemetry", "error", err.Error())
		return h, func(_ context.Context) {
			l.Warn("no Open Telemetry stack to shutdown")
		}
	}
	shutdown := func(ctx context.Context) {
		err := shutdowns(ctx)
		if err != nil {
			l.Error("Error shutting down Open Telemetry", "error", err.Error())
		}
	}
	handler := otelhttp.NewHandler(h, "/")
	return handler, shutdown
}

// Setup bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func Setup(ctx context.Context, opts ...Opt) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	var err error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(fromContext(ctx, CtxAppName)),
		semconv.ServiceVersionKey.String(fromContext(ctx, CtxVersion)),
		semconv.ServiceInstanceIDKey.String(fromContext(ctx, CtxInstanceID)),
	)
	for _, opt := range opts {
		optErr := opt(ctx, res, shutdownFuncs)
		if optErr != nil {
			handleErr(optErr)
			return shutdown, err
		}
	}

	return shutdown, err
}

func fromContext(ctx context.Context, key string) string {
	val, ok := ctx.Value(key).(string)
	if ok {
		return val
	}
	return ""
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider(ctx context.Context, res *resource.Resource) (*trace.TracerProvider, error) {
	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return nil, err
	}
	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(spanExporter),
	)
	return tracerProvider, nil
}

func newMeterProvider(ctx context.Context, res *resource.Resource) (*metric.MeterProvider, error) {
	// Create metric reader using environment variables
	metricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return nil, err
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metricReader),
		metric.WithResource(res),
	)
	return meterProvider, nil
}
