package bot

import (
	"log"
	config "senior_intern_bot/config"

	"github.com/bwmarrin/discordgo"
)

type SessionHandler struct {
	session *discordgo.Session
	config.Config
}

// SendToInternshipChannel
func (sh *SessionHandler) SendToInternshipChannel(content string) error {
	_, err := sh.session.ChannelMessageSend(sh.InternshipChannelID, content)
	return err
}

// SendToUnsortedChannel
func (sh *SessionHandler) SendToUnsortedChannel(content string) error {
	_, err := sh.session.ChannelMessageSend(sh.UnsortedChannelID, content)
	return err
}

func New(cfg config.Config) (*SessionHandler, error) {
	// Create a session
	discordSession, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		log.Println("Error message: ", err)
		return nil, err
	}

	return &SessionHandler{
		session: discordSession,
		Config:  cfg,
	}, nil
}

func (sh *SessionHandler) Open() error {
	if err := sh.session.Open(); err != nil {
		log.Println("Error starting session: ", err)
		return err
	}

	// TO-BE-REMOVED
	if _, err := sh.session.ChannelMessageSend(sh.InternshipChannelID, "Bot running..."); err != nil {
		log.Println("Error sending startup message: ", err)
		return err
	}
	return nil
}

func (sh *SessionHandler) Close() error {
	if err := sh.session.Close(); err != nil {
		log.Println("Error closing session: ", err)
		return err
	}
	return nil
}
