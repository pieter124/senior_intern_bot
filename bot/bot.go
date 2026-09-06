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

func New(cfg config.Config) SessionHandler {
	// Create a session
	discordSession, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		log.Fatal("Error message: ", err)
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

func (sh *SessionHandler) Open() {
	if err := sh.session.Open(); err != nil {
		log.Println("Error starting session: ", err)
		return
	}

	if _, err := sh.session.ChannelMessageSend(sh.InternshipChannelID, "Bot running..."); err != nil {
		log.Println("Error sending startup message: ", err)
	}
}

func (sh *SessionHandler) Close() {
	if err := sh.session.Close(); err != nil {
		log.Println("Error closing session: ", err)
	}
}
