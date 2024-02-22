package http

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

func NewServer(logger *slog.Logger, handler http.Handler, address string) *Server {
	return &Server{
		logger:  logger,
		handler: handler,
		address: address,
	}
}

func (s *Server) Listen(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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
			s.logger.Error("listen and serve returned err", err)
			return
		}
	}()

	<-ctx.Done()
	s.logger.Info("got interruption signal")
	if err := srv.Shutdown(context.TODO()); err != nil {
		s.logger.Error("server shutdown returned an err", err)
		return err
	}
	return nil
}
