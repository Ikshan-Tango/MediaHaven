package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mediahaven/pkg/config"

	"github.com/bwmarrin/discordgo"
)

func UploadToDiscord(filename string, fileContent io.Reader) error {
	dg, err := discordgo.New("Bot " + config.Get().BotToken)
	discordChannelID := config.Get().ChannelId
	if err != nil {
		return err
	}
	_, err = dg.ChannelFileSend(discordChannelID, filename, fileContent)
	return err
}

// Encrypt data using AES
func EncryptFile(data []byte) ([]byte, error) {
	// Decode the hex key into a byte slice
	key, err := hex.DecodeString(config.Get().SecretKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %v", err)
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a GCM mode for encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt the data
	encrypted := gcm.Seal(nonce, nonce, data, nil)
	return encrypted, nil
}
