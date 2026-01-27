package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/telegram"
)

type Server struct {
	cfg        *config.App
	log        logger.Logger
	httpServer *http.Server
}

func New(cfg *config.App, log logger.Logger, handler http.Handler) *Server {
	return &Server{
		cfg: cfg,
		log: log,
		httpServer: &http.Server{
			Addr:           cfg.GetDSN(),
			Handler:        handler,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			IdleTimeout:    cfg.IdleTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
		},
	}
}

func (s *Server) Run() error {
	go func() {
		s.log.Info("Server is starting on port " + s.httpServer.Addr)
		telegram.Send("✅ <b>Application started successfully</b>")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatal("Server listen error", logger.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	telegram.Send("🚨 <b>Application stopping...</b>")
	s.log.Info("Server is shutting down...")
	
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	
	if err := s.httpServer.Shutdown(ctx); err != nil {
        s.log.Error("Server forced to shutdown", logger.Error(err))
    }

    return nil
}
