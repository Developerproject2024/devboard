package main

import (
	"net/http"
	"os"

	_ "github.com/Developerproject2024/devboard/docs"
	"github.com/Developerproject2024/devboard/internal/handler"
	"github.com/Developerproject2024/devboard/internal/logger"
	"github.com/Developerproject2024/devboard/internal/middlewares"
	"github.com/Developerproject2024/devboard/internal/server"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Devboard API
// @version 1.0
// @description API REST para gestión de proyectos
// @contact.name Soporte Devboard Ricardo
// @host  localhost:8080
// @BasePath /api/v1

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger := logger.New(logger.ProductionConfig())

	server := server.New(":"+port, logger)

	server.Use(middlewares.Recovery(logger))
	server.Use(middlewares.Logger(logger))

	healthHandler := handler.NewHealthHandler()

	server.RegisterRoutes("GET /docs/", httpSwagger.WrapHandler)
	server.RegisterRoutes("GET /health", healthHandler)
	server.RegisterRoutes("GET /ready", healthHandler)

	// Ruta temporal para probar Recovery
	server.RegisterRoutes("GET /panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("error provocado para probar Recovery")
	}))

	if err := server.Start(); err != nil {
		logger.Error("Error starting server fatal", "error", err)
		os.Exit(1)
	}

}
