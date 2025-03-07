package main

import (
	"log"
	"mediahaven/pkg/discord/controller"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// Load env if not set (local dev)
	if _, localEnvSet := os.LookupEnv("PRODUCTION"); !localEnvSet {
		log.Println("loading variables from local env")
		godotenv.Load()
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health", controller.Health)
	e.POST("/upload", controller.Upload)
	e.POST("/download", controller.Download)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
