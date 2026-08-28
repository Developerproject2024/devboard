package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	logger     *slog.Logger
}

func New(addr string, logger *slog.Logger) *Server {
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
		logger:     logger,
	}
}

func (server *Server) Start() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	serverErr := make(chan error, 1)
	go func() {
		server.logger.Info("Starting server", slog.String("addr", server.httpServer.Addr))
		if error := server.httpServer.ListenAndServe(); error != nil && error != http.ErrServerClosed {
			serverErr <- error
		}
	}()

	select {
	case sig := <-quit:
		server.logger.Info("Shutting down server...", slog.String("signal", sig.String()))
		/*ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.httpServer.Shutdown(ctx); err != nil {
			server.logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
			return err
		}
		server.logger.Info("Server exiting")*/
	case err := <-serverErr:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.httpServer.Shutdown(ctx); err != nil {
		server.logger.Error("Server forced to shutdown", slog.Any("error", err))
		return err
	}
	server.logger.Info("Server exiting")

	return nil
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
