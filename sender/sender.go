package sender

import (
	bot "senior_intern_bot/bot"
)

type Sender struct {
	sessionHandler bot.SessionHandler
}

func New(sessionHandler bot.SessionHandler) (Sender, error) {
	return Sender{
		sessionHandler: sessionHandler,
	}, nil
}
