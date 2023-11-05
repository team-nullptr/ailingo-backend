package httpserver

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"
)

const (
	defaultShutdownTimeout = 5 * time.Second
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

type Option func(s *Server)

func WithHandler(handler http.Handler) Option {
	return func(s *Server) {
		s.server.Handler = handler
	}
}

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.server.Addr = addr
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

func New(opts ...Option) *Server {
	httpserver := &http.Server{
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	server := &Server{
		server:          httpserver,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) StartTLS(certFile string, keyFile string) {
	go func() {
		s.notify <- s.server.ListenAndServeTLS(certFile, keyFile)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
