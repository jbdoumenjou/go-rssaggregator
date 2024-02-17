package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

// Server represents an HTTP server.
type Server struct {
	*http.Server
}

// NewServer creates a new server.
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}

}

// Start runs the HTTP server.
func (s *Server) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("starting server on %s\n", s.Addr)
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("got interruption signal")
	if err := s.Shutdown(context.TODO()); err != nil {
		return fmt.Errorf("server shutdown returned an err: %w\n", err)
	}

	log.Println("server closed")
	return nil
}
