package otel

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

// taken from: https://opentelemetry.io/docs/languages/go/getting-started/

type Opt func(ctx context.Context, fns []func(context.Context) error) error

func WithTracing() Opt {
	return func(ctx context.Context, fns []func(context.Context) error) error {
		tracerProvider, err := newTracerProvider(ctx)
		if err != nil {
			return err
		}
		fns = append(fns, tracerProvider.Shutdown)
		otel.SetTracerProvider(tracerProvider)
		return nil
	}
}

func WithMetrics() Opt {
	return func(ctx context.Context, fns []func(context.Context) error) error {
		meterProvider, err := newMeterProvider(ctx)
		if err != nil {
			return err
		}
		fns = append(fns, meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)
		return nil
	}
}

func WithLogger() Opt {
	return func(ctx context.Context, fns []func(context.Context) error) error {
		loggerProvider, err := newLoggerProvider()
		if err != nil {
			return err
		}
		fns = append(fns, loggerProvider.Shutdown)
		global.SetLoggerProvider(loggerProvider)
		return nil
	}
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

	for _, opt := range opts {
		optErr := opt(ctx, shutdownFuncs)
		if optErr != nil {
			handleErr(optErr)
			return shutdown, err
		}
	}

	return shutdown, err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider(ctx context.Context) (*trace.TracerProvider, error) {
	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return nil, err
	}
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(spanExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
	)
	return tracerProvider, nil
}

func newMeterProvider(ctx context.Context) (*metric.MeterProvider, error) {
	// Create metric reader using environment variables
	metricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return nil, err
	}
	meterProvider := metric.NewMeterProvider(metric.WithReader(metricReader))
	return meterProvider, nil
}

func newLoggerProvider() (*log.LoggerProvider, error) {
	logExporter, err := stdoutlog.New(stdoutlog.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}
