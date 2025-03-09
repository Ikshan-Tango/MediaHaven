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

	fileContent, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	encryptedContent, err := services.EncryptFile(fileContent)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to encrypt file: %s", err),
		})
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "upload-*.tmp")
	if err != nil {
		return err
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Write the encrypted content to the temporary file
	if _, err = tempFile.Write(encryptedContent); err != nil {
		return err
	}

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

func Download(c echo.Context) error {
	// Get the file name from query parameters
	fileName := c.QueryParam("filename")
	if fileName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Filename is required",
		})
	}

	// Fetch the file from Discord
	encryptedFileContent, err := services.DownloadFromDiscord(fileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to fetch file from Discord: %s", err),
		})
	}
	// Decrypt the file content
	decryptedContent, err := services.DecryptFile(encryptedFileContent)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to decrypt file: %s", err),
		})
	}
	// Determine the MIME type based on the file extension
	contentType := "application/octet-stream" // Default MIME type
	if len(fileName) > 4 {
		extension := fileName[len(fileName)-4:]
		switch extension {
		case ".jpg", "jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		}
	}

	// Serve the file to the client
	return c.Blob(http.StatusOK, contentType, decryptedContent)
}
