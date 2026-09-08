package sender

import (
	"errors"
	"fmt"
	bot "senior_intern_bot/bot"
	domain "senior_intern_bot/domain"
	"unicode/utf8"
)

const (
	discordMsgLimit = 2000 // Discord character limit per message.
	maxTitleRunes   = 200  // Constant used to truncate title of posting to 200 chars.
)

type Sender interface {
	Send(map[domain.Verdict][]domain.Posting) ([]domain.Posting, error)
}

type DiscordSender struct {
	sessionHandler *bot.SessionHandler
}

func New(sessionHandler *bot.SessionHandler) *DiscordSender {
	return &DiscordSender{
		sessionHandler: sessionHandler,
	}
}

type batch struct {
	content  string
	postings []domain.Posting
}

func (b *batch) add(line string, p domain.Posting) {
	if b.content != "" {
		b.content += "\n"
	}
	b.content += line
	b.postings = append(b.postings, p)
}

func formatPosting(posting domain.Posting) string {
	title := posting.Title
	if utf8.RuneCountInString(title) > maxTitleRunes {
		title = string([]rune(title)[:maxTitleRunes]) + "..."
	}
	return fmt.Sprintf("%s | %s : %s", title, posting.Company, posting.URL)
}

// discordBatchify packs postings into batches without exceeding the Discord character limit.
// Every posting appears in one batch only.
func discordBatchify(postings []domain.Posting) []batch {
	var (
		batches []batch
		current batch
	)

	// Makes the batches, ensuring each batch does not exceed the Discord character limit.
	for _, p := range postings {
		line := formatPosting(p)
		if len(current.postings) > 0 && utf8.RuneCountInString(current.content)+utf8.RuneCountInString(line) >= discordMsgLimit {
			batches = append(batches, current)
			current = batch{}
		}
		current.add(line, p)
	}
	if len(current.postings) > 0 {
		batches = append(batches, current)
	}
	return batches
}

func sendBatches(postings []domain.Posting, send func(string) error) ([]domain.Posting, []error) {
	var (
		sent []domain.Posting
		errs []error
	)

	for _, batch := range discordBatchify(postings) {
		if err := send(batch.content); err != nil {
			errs = append(errs, err)
			continue
		}
		sent = append(sent, batch.postings...)
	}
	return sent, errs
}

// Send handles sending the postings (categorised by verdict), which happens after the filtering stage.
// REJECT postings are deliberately not sent anywhere.
func (sender *DiscordSender) Send(postingsByVerdict map[domain.Verdict][]domain.Posting) ([]domain.Posting, error) {
	sentAcceptPostings, errsFromAcceptedMessages := sendBatches(postingsByVerdict[domain.ACCEPT], sender.sessionHandler.SendToInternshipChannel)
	sentReviewPostings, errsFromReviewMessages := sendBatches(postingsByVerdict[domain.REVIEW], sender.sessionHandler.SendToUnsortedChannel)

	sent := make([]domain.Posting, 0, len(sentAcceptPostings)+len(sentReviewPostings))
	sent = append(sent, sentAcceptPostings...)
	sent = append(sent, sentReviewPostings...)

	errs := make([]error, 0, len(errsFromAcceptedMessages)+len(errsFromReviewMessages))
	errs = append(errs, errsFromAcceptedMessages...)
	errs = append(errs, errsFromReviewMessages...)

	return sent, errors.Join(errs...)
}
