package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server holds dependencies and implements the GracefulShutdownServer interface.
type Server struct {
	httpServer    *http.Server
	logger        *log.Logger
	cancelTimeout time.Duration
	killChannel   chan os.Signal
}

// GracefulShutdownServer is the interface that defines a server capable of graceful shutdown.
type GracefulShutdownServer interface {
	Run()
}

// New returns a Server configured with the http server passed as argument.
func New(s *http.Server, l *log.Logger, timeout time.Duration) *Server {
	return &Server{
		httpServer:    s,
		logger:        l,
		cancelTimeout: timeout,
		killChannel:   make(chan os.Signal, 1),
	}
}

// Run starts the server by telling it to listen on the configured port and to serve HTTP requests.
// Then it waits for a SIGINT/SIGTERM signal to initiate the
func (s *Server) Run() {
	signal.Notify(s.killChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Listen: %s\n", err)
		}
	}()
	s.logger.Printf("Server started on %s", s.httpServer.Addr)

	<-s.killChannel
	s.logger.Printf("Server stopping")

	//ctx, cancel := context.WithTimeout(context.Background(), s.cancelTimeout)
	//defer cancel()

	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		s.logger.Fatalf("Server shutdown failed: %+v", err)
	}
	s.logger.Printf("Server terminated gracefully")
}

func (s *Server) Stop() {
	s.killChannel <- syscall.SIGINT
}
