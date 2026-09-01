package sender

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Message struct {
	Title string
	URL   string
}

type Dispatcher struct {
	Messages       chan Message
	DiscordSession *discordgo.Session
	ChannelID      string
}

func New(discordSession *discordgo.Session, channelId string) *Dispatcher {
	return &Dispatcher{
		Messages:       make(chan Message, 1),
		DiscordSession: discordSession,
		ChannelID:      channelId,
	}
}

func (d *Dispatcher) Run() {
	go func() {
		for message := range d.Messages {
			if _, err := d.DiscordSession.ChannelMessageSend(d.ChannelID, message.Title+": "+message.URL); err != nil {
				log.Println("sender: error sending message: ", err)
			}
		}
	}()
}
