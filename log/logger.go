package log

import (
	"io"
	"log/slog"
)

type Format int

const (
	JSON Format = iota
	TEXT
)

func Build(out io.Writer, f Format, debug bool) *slog.Logger {
	var opts *slog.HandlerOptions
	if debug {
		opts = &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}
	}
	var handler slog.Handler
	switch f {
	case JSON:
		handler = slog.NewJSONHandler(out, opts)
	default:
		handler = slog.NewTextHandler(out, opts)
	}

	return slog.New(handler)
}
