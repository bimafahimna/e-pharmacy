package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/logger"
)

type Server struct {
	http.Server
	GracePeriod time.Duration
}

func NewServer(config *config.Config, handler http.Handler) Server {
	return Server{
		Server: http.Server{
			Addr:    config.App.ServerAddress,
			Handler: handler,
		},
		GracePeriod: config.App.ServerGracePeriod,
	}
}

func (s *Server) Run() {
	go func() {
		logger.Log.Infof("listening on port %s...", s.Addr)
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatalf("failed to start run server: %v", err)
		}
		logger.Log.Info("server is not receiving new requests...")
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	timeout := time.Duration(s.GracePeriod) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	<-stop

	logger.Log.Info("shutting down server...")
	if err := s.Shutdown(ctx); err != nil {
		logger.Log.Fatalf("failed to shutdown server: %v", err)
	}

	logger.Log.Info("server shutdown gracefully")
}
