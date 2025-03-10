package controller

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"mediahaven/pkg/discord/services"

	"github.com/labstack/echo/v4"
)

func Upload(c echo.Context) error {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to retrieve the file",
		})
	}

	go func() {
		// Open the file
		src, err := file.Open()
		if err != nil {
			log.Println("Failed to open file : ", err)
		}
		defer src.Close()

		// Read file content
		fileContent, err := io.ReadAll(src)
		if err != nil {
			log.Println("Failed to read file content : ", err)
		}

		// Encrypt the file content
		encryptedContent, err := services.EncryptFile(fileContent)
		if err != nil {
			log.Println("Failed to encrypt file content : ", err)
		}

		// Split into 8MB chunks (Discord-friendly size)
		chunkSize := 8 * 1024 * 1024 // 8MB
		chunks := splitIntoChunks(encryptedContent, chunkSize)

		// Upload each chunk to Discord
		for i, chunk := range chunks {
			chunkName := fmt.Sprintf("%s.part%d", file.Filename, i+1)

			// Create in-memory reader for the chunk
			chunkReader := bytes.NewReader(chunk)

			// Upload directly from memory
			if err := services.UploadToDiscord(chunkName, chunkReader); err != nil {
				log.Printf("Failed to upload chunk %d: %v", i+1, err)

			}
		}
	}()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "File upload started to discord",
		"filename": file.Filename,
	})
}

// Helper function to split data into chunks
func splitIntoChunks(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}

func Download(c echo.Context) error {
	// Get the original file name from query parameters
	fileName := c.QueryParam("filename")
	if fileName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Filename is required",
		})
	}

	// Get all file chunks from Discord
	chunkedContent, err := services.DownloadFromDiscord(fileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to fetch file chunks: %s", err),
		})
	}

	// Combine and decrypt the chunks
	decryptedContent, err := services.CombineAndDecryptChunks(chunkedContent)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to process file: %s", err),
		})
	}

	// Determine MIME type from original filename
	contentType := getContentType(fileName)

	return c.Blob(http.StatusOK, contentType, decryptedContent)
}

func getContentType(fileName string) string {
	// Remove .partXX suffix if present
	cleanName := strings.Split(fileName, ".part")[0]

	if ext := filepath.Ext(cleanName); ext != "" {
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg":
			return "image/jpeg"
		case ".png":
			return "image/png"
		case ".gif":
			return "image/gif"
		}
	}
	return "application/octet-stream"
}
