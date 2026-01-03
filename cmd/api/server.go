package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/semmidev/ethos-go/config"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// Server wraps the HTTP server with graceful shutdown capabilities
type Server struct {
	httpServer *http.Server
	logger     logger.Logger
}

// NewServer creates a new HTTP server with the given configuration
func NewServer(cfg *config.Config, router chi.Router, logger logger.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.ServerHost + ":" + cfg.ServerPort,
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		logger: logger,
	}
}

// Start begins listening for HTTP requests
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(ctx, "starting HTTP server",
		logger.Field{Key: "addr", Value: s.httpServer.Addr},
	)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info(ctx, "shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

// Addr returns the server address
func (s *Server) Addr() string {
	return s.httpServer.Addr
}
