package services

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io"
	"mediahaven/pkg/config"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func DownloadFromDiscord(baseName string) ([][]byte, error) {
	botToken := config.Get().BotToken
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, fmt.Errorf("discord session failed: %w", err)
	}

	discordChannelID := config.Get().ChannelId
	messages, err := dg.ChannelMessages(discordChannelID, 100, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("message fetch failed: %w", err)
	}

	var chunks [][]byte
	var chunkNumbers []int

	// Collect all matching chunks
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if strings.HasPrefix(attachment.Filename, baseName+".part") {
				// Extract chunk number
				parts := strings.Split(attachment.Filename, ".part")
				if len(parts) != 2 {
					continue
				}

				chunkNum, err := strconv.Atoi(parts[1])
				if err != nil {
					continue
				}

				// Download chunk content
				resp, err := http.Get(attachment.URL)
				if err != nil {
					return nil, fmt.Errorf("chunk download failed: %w", err)
				}
				defer resp.Body.Close()

				content, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("chunk read failed: %w", err)
				}

				chunks = append(chunks, content)
				chunkNumbers = append(chunkNumbers, chunkNum)
			}
		}
	}

	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunks found for %s", baseName)
	}

	// Sort chunks by part number
	sort.Slice(chunks, func(i, j int) bool {
		return chunkNumbers[i] < chunkNumbers[j]
	})

	return chunks, nil
}

func CombineAndDecryptChunks(chunks [][]byte) ([]byte, error) {
	// Combine chunks in order
	var combined []byte
	for _, chunk := range chunks {
		combined = append(combined, chunk...)
	}

	// Decrypt the combined content
	decrypted, err := DecryptFile(combined)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return decrypted, nil
}

// Decrypt data using AES
func DecryptFile(data []byte) ([]byte, error) {
	// Decode the hex key into a byte slice
	key, err := hex.DecodeString(config.Get().SecretKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a GCM mode for decryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract the nonce from the encrypted data
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt the data
	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
