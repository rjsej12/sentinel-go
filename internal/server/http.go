package server

import (
	"context"
	"log"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(addr string, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *HTTPServer) Start() {
	go func() {
		log.Printf("HTTP server started on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	log.Println("shutting down http server...")
	return s.server.Shutdown(ctx)
}
