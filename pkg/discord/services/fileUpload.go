package services

import (
	"log"
	"mediahaven/pkg/config"
	"os"

	"github.com/bwmarrin/discordgo"
)

func UploadToDiscord(filePath, fileName string) error {
	log.Println("filepath: ", filePath)
	discordToken := config.Get().BotToken
	log.Println("discordToken: ", discordToken)
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
	discordChannelID := "971806504975994941"
	_, err = dg.ChannelFileSend(discordChannelID, fileName, file)
	if err != nil {
		return err
	}

	return nil
}
