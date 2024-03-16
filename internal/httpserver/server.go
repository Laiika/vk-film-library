package httpserver

import (
	"net/http"
)

type Server struct {
	server *http.Server
}

func New(handler http.Handler, address string) *Server {
	httpServer := &http.Server{
		Handler: handler,
		Addr:    address,
	}

	s := &Server{
		server: httpServer,
	}

	return s
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.server.Close()
}
