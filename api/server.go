package api

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const (
	ServerDefaultTimeOut = 10 * time.Second
	ServerWriteTimeOut   = ServerDefaultTimeOut
	ServerReadTimeOut    = ServerDefaultTimeOut
)

type Server struct {
	logger  *slog.Logger
	handler http.Handler
	address string
}

func NewServer(l *slog.Logger, h http.Handler, a string) *Server {
	return &Server{
		logger:  l,
		handler: h,
		address: a,
	}
}

func (s *Server) Listen(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Handler:      s.handler,
		Addr:         s.address,
		BaseContext:  func(net.Listener) context.Context { return ctx },
		WriteTimeout: ServerWriteTimeOut,
		ReadTimeout:  ServerReadTimeOut,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("listen and serve returned err", slog.Any("err", err))
			return
		}
	}()

	<-ctx.Done()
	s.logger.Info("got interruption signal")
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Error("server shutdown returned an err", slog.Any("err", err))
		return err
	}
	return nil
}
