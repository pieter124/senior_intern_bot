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

func (sh *SessionHandler) SendToInternshipChannel(content string) error {
	_, err := sh.session.ChannelMessageSend(sh.InternshipChannelID, content)
	return err
}

func New(cfg config.Config) SessionHandler {
	// Create a session
	discordSession, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		log.Println("Error message: ", err)
	}

	//	fmt.Println("Bot running...")
	//if _, err := discordSession.ChannelMessageSend(cfg.InternshipChannelID, "Bot running..."); err != nil {
	//log.Println("Error sending startup message: ", err)
	//}
	return SessionHandler{
		session: discordSession,
		Config:  cfg,
	}
}

func (sh *SessionHandler) Open() error {
	if err := sh.session.Open(); err != nil {
		log.Println("Error starting session: ", err)
		return err
	}

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
