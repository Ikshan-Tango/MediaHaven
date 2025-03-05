package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"mediahaven/pkg/discord/services"

	"github.com/labstack/echo/v4"
)

func Upload(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "upload-*.tmp")
	if err != nil {
		return err
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Copy the uploaded file to the temporary file
	if _, err = io.Copy(tempFile, src); err != nil {
		return err
	}

	// Upload the file to Discord
	err = services.UploadToDiscord(tempFile.Name(), file.Filename)
	if err != nil {
		log.Println("ERROR - while uploading file to Discord:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to upload file to Discord: %s", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "File uploaded to Discord successfully!",
	})
}
