package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mediahaven/pkg/config"
	"os"

	"github.com/bwmarrin/discordgo"
)

func UploadToDiscord(filePath, fileName string) error {
	discordToken := config.Get().BotToken
	// Create a new Discord session
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		return err
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Upload the file to Discord
	discordChannelID := config.Get().ChannelId
	_, err = dg.ChannelFileSend(discordChannelID, fileName, file)
	if err != nil {
		return err
	}

	return nil
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
