package services

import (
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
	log.Println("Bot Token:", botToken)
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Println("ERROR - while creating Discord session:", err)
		return nil, err
	}
	discordChannelID := "971806504975994941"
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
