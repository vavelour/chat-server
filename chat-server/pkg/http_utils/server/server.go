package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg HttpServerConfig, h http.Handler) *Server {
	return &Server{httpServer: &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Addr),
		Handler:        h,
		MaxHeaderBytes: 1 << cfg.MaxHeaderBytes,
		ReadTimeout:    cfg.ReadTimeout * time.Second,
		WriteTimeout:   cfg.WriteTimeout * time.Second,
	}}
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listening and serving: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
