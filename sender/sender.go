package sender

import (
	"fmt"
	"log"
	bot "senior_intern_bot/bot"
	domain "senior_intern_bot/domain"
)

type Sender interface {
	Send([]domain.Posting) error
}

type DiscordSender struct {
	sessionHandler *bot.SessionHandler
}

func New(sessionHandler *bot.SessionHandler) *DiscordSender {
	return &DiscordSender{
		sessionHandler: sessionHandler,
	}
}

func (sender *DiscordSender) formatPostingsToMessages(postings []domain.Posting) ([]string, error) {
	var formattedMessages []string
	for _, p := range postings {
		message := fmt.Sprintf("%s | %s : %s", p.Title, p.Company, p.URL)
		formattedMessages = append(formattedMessages, message)
	}
	return formattedMessages, nil
}

func (sender *DiscordSender) Send(postings []domain.Posting) error {
	formattedMessages, err := sender.formatPostingsToMessages(postings)
	if err != nil {
		return fmt.Errorf("error formatting postings to messages: %s", err)
	}

	var currentMessage string
	var messagesToBatch []string
	for _, msg := range formattedMessages {
		if len(currentMessage)+len(msg) <= 1000 {
			currentMessage += fmt.Sprintf("\n%s", msg) // Add to current message to be batched
		} else {
			messagesToBatch = append(messagesToBatch, currentMessage)
			currentMessage = msg
		}
	}
	messagesToBatch = append(messagesToBatch, currentMessage)

	for _, batch := range messagesToBatch {
		if err := sender.sessionHandler.SendToInternshipChannel(batch); err != nil {
			log.Println("Error sending message: ", err)
		}
	}
	return nil
}
