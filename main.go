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

	cfg := config.Config{
		BotToken:            botToken,
		InternshipChannelID: internshipChannelID,
	}

	// Start bot session
	botSession := bot.New(cfg)
	botSession.Open()
	defer botSession.Close()

	greenhouseCompanies := []string{"monzo", "anthropic", "stripe", "jetbrains", "cloudflare",
		"mongodb", "canonical", "samsara", "celonis", "hellofresh", "doctolib",
		"airbnb", "databricks", "squarespace",
	}

	pollerServ, err := poller.New(greenhouseCompanies)
	if err != nil {
		log.Println("error creating poller: ", err)
	}

	filterServ, err := filter.New()
	if err != nil {
		log.Println("error creating filter: ", err)
	}

	senderServ, err := sender.New(&botSession)
	if err != nil {
		log.Println("error creating sender: ", err)
	}

	storerServ, err := storer.New()
	if err != nil {
		log.Println("error creating storer: ", err)
	}

	// Start ticker (determines how often we poll)
	ticker := time.NewTicker(1 * time.Minute)

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
