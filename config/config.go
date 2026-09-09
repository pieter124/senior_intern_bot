package config

type Config struct {
	BotToken            string
	InternshipChannelID string // Postings with Verdict=ACCEPT
	UnsortedChannelID   string // Postings with Verdict=REVIEW
	DatabaseDSN         string
}
