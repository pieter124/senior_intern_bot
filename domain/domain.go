package domain

import "fmt"

type Verdict int

const (
	REVIEW Verdict = iota
	ACCEPT
	REJECT
)

var verdictNames = [...]string{
	REVIEW: "REVIEW",
	ACCEPT: "ACCEPT",
	REJECT: "REJECT",
}

func (v Verdict) String() string {
	if v < 0 || int(v) >= len(verdictNames) {
		return fmt.Sprintf("Verdict(%d)", int(v))
	}
	return verdictNames[v]
}

type Location struct {
	Name string `json:"name"`
}
type Posting struct {
	Title       string   `json:"title"`
	JobID       int      `json:"id"`
	Location    Location `json:"location"`
	PostedAt    string   `json:"first_published"`
	URL         string   `json:"absolute_url"`
	Content     string   `json:"content"`
	Company     string   `json:"-"` // Set during polling.
	RawJSON     string   `json:"-"` // Set during polling.
	Verdict     Verdict  `json:"-"` // Set during filtering.
	ContentHash string   `json:"-"` // Set during storing.
}
