package services

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mediahaven/pkg/config"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func DownloadFromDiscord(fileName string) ([]byte, error) {
	// Create a new Discord session
	botToken := config.Get().BotToken
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Println("ERROR - while creating Discord session:", err)
		return nil, err
	}
	discordChannelID := config.Get().ChannelId
	// Fetch messages from the channel
	messages, err := dg.ChannelMessages(discordChannelID, 100, "", "", "")
	if err != nil {
		log.Println("ERROR - while fetching messages from Discord:", err)
		return nil, err
	}

	// Iterate through messages to find the file
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if attachment.Filename == fileName {
				// Download the file
				resp, err := http.Get(attachment.URL)
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()

				// Read the file content
				fileContent, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, err
				}

				return fileContent, nil
			}
		}
	}

	return nil, fmt.Errorf("file not found")

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
