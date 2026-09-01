package main

import (
	"log"
	"os"
	bot "senior_intern_bot/bot"

	"github.com/joho/godotenv"
)

func main() {
	// Load env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get env variable for bot token
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable is not set")
	}

	// Get env variable for the #internships channel ID
	internshipsChannel := os.Getenv("DISCORD_CHANNEL_INTERNSHIPS")
	if internshipsChannel == "" {
		log.Fatal("DISCORD_CHANNEL_INTERNSHIPS environment variable is not set")
	}

	// Start bot
	bot.BotToken = botToken
	bot.InternshipsChannelID = internshipsChannel
	bot.Run()
}
