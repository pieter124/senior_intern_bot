package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	poller "senior_intern_bot/poller"
	sender "senior_intern_bot/sender"

	"github.com/bwmarrin/discordgo"
)

var BotToken string
var InternshipsChannelID string

func Run() {
	// Create a session
	discordSession, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal("Error message: ", err)
	}
	// Open session
	discordSession.Open()
	defer discordSession.Close() // Close session, after function termination.

	fmt.Println("Bot running...")
	if _, err := discordSession.ChannelMessageSend(InternshipsChannelID, "Bot running..."); err != nil {
		log.Println("Error sending startup message: ", err)
	}

	// Start the db connection

	// Start the sender
	dispatcher := sender.New(discordSession, InternshipsChannelID)
	dispatcher.Run()

	// Start the poller
	poller.Run(dispatcher.Messages)

	// Keep bot running until OS interruption (e.g ctrl + c)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
