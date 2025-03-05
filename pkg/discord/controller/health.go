package controller

import (
	"log"

	"github.com/labstack/echo/v4"
)

func Health(c echo.Context) error {
	log.Println("Health check")
	return c.String(200, "OK")
}
