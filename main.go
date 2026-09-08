package main

import (
	"log"
	"os"
	"os/signal"
	bot "senior_intern_bot/bot"
	config "senior_intern_bot/config"
	filter "senior_intern_bot/filter"
	orchestrator "senior_intern_bot/orchestrator"
	poller "senior_intern_bot/poller"
	sender "senior_intern_bot/sender"
	storer "senior_intern_bot/storer"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable is not set")
	}

	internshipChannelID := os.Getenv("DISCORD_CHANNEL_INTERNSHIPS")
	if internshipChannelID == "" {
		log.Fatal("DISCORD_CHANNEL_INTERNSHIPS environment variable is not set")
	}

	unsortedChannelID := os.Getenv("DISCORD_CHANNEL_UNSORTED")
	if unsortedChannelID == "" {
		log.Fatal("DISCORD_CHANNEL_UNSORTED environment variable is not set")
	}

	cfg := config.Config{
		BotToken:            botToken,
		InternshipChannelID: internshipChannelID,
		UnsortedChannelID:   unsortedChannelID,
	}

	// Start bot session
	botSession, err := bot.New(cfg)
	if err != nil {
		log.Println("error creating new bot session: ", err)
		return
	}

	err = botSession.Open()
	if err != nil {
		log.Println("error opening up bot session: ", err)
		return
	}
	defer botSession.Close()

	greenhouseCompanies := []string{"monzo"}
	/*
		"anthropic", "stripe", "jetbrains", "cloudflare",
			"mongodb", "canonical", "samsara", "celonis", "hellofresh", "doctolib",
			"airbnb", "databricks", "squarespace"
	*/

	// Start up pipeline services.
	pollerServ := poller.New(greenhouseCompanies)
	filterServ := filter.New()
	senderServ := sender.New(botSession)
	storerServ := storer.New()

	// Start ticker (determines how often we poll)
	ticker := time.NewTicker(20 * time.Minute)

	orchestratorServ, err := orchestrator.New(
		pollerServ,
		filterServ,
		senderServ,
		storerServ,
		ticker,
	)
	if err != nil {
		log.Println("error creating orchestrator: ", err)
	}

	if err := orchestratorServ.Start(); err != nil {
		log.Println("error starting orchestrator: ", err)
	}

	// Keep bot running until OS interruption (e.g ctrl + c)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	if err := orchestratorServ.Stop(); err != nil {
		log.Println("orchestrator has not gracefully stopped")
	} else {
		log.Println("orchestrator has gracefully stopped")
	}

}
