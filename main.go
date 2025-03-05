package main

import (
	"mediahaven/pkg/discord/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health", controller.Health)
	e.POST("/upload", controller.Upload)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
