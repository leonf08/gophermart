package http

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	err    chan error
}

// NewServer creates a new HTTP server.
func NewServer(handler http.Handler, address string) *Server {
	s := &Server{
		server: &http.Server{
			Addr:    address,
			Handler: handler,
		},
		err: make(chan error, 1),
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.err <- s.server.ListenAndServe()
		close(s.err)
	}()
}

// Err returns a channel with an error.
func (s *Server) Err() <-chan error {
	return s.err
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
