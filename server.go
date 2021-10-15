package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type APIServer struct {
	httpServer *http.Server
}

func NewAPIServer(addr string, port int) *APIServer {
	return &APIServer{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", addr, port),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}

func (s *APIServer) SetRoutes(routes http.Handler) {
	s.httpServer.Handler = routes
}

func (s *APIServer) Start() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()
	log.Infof("Server started on %s", s.httpServer.Addr)

	<-done
	log.Info("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %+v", err)
	}
	log.Info("Server exited gracefully")
}
