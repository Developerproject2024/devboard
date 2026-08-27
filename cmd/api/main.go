package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Developerproject2024/devboard/internal/handler"
	"github.com/Developerproject2024/devboard/internal/logger"
	"github.com/Developerproject2024/devboard/internal/middlewares"
	"github.com/Developerproject2024/devboard/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger := logger.New(logger.DefaultConfig())

	server := server.New(":" + port)

	server.Use(middlewares.Recovery(logger))
	server.Use(middlewares.Logger(logger))

	healthHandler := handler.NewHealthHandler()

	server.RegisterRoutes("GET /health", healthHandler)
	server.RegisterRoutes("GET /ready", healthHandler)

	// Ruta temporal para probar Recovery
	server.RegisterRoutes("GET /panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("error provocado para probar Recovery")
	}))

	logger.Info("servidor iniciando", slog.String("addr", ":8080"))

	if err := server.Start(); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

}
