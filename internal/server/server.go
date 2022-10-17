package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(path string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + path,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
