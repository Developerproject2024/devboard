package server

import (
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
}

func New(addr string) *Server {
	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,

		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return &Server{
		httpServer: httpServer,
		mux:        mux,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (server *Server) RegisterRoutes(pattern string, handler http.Handler) {
	server.mux.Handle(pattern, handler)
}

func (server *Server) ServerHttp(write http.ResponseWriter, request *http.Request) {
	server.mux.ServeHTTP(write, request)
}

func (server *Server) Use(middleware func(http.Handler) http.Handler) {
	server.httpServer.Handler = middleware(server.httpServer.Handler)
}
