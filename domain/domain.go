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
	Title     string   `json:"title"`
	JobID     int      `json:"id"`
	Location  Location `json:"location"`
	PostedAt  string   `json:"first_published"`
	UpdatedAt string   `json:"updated_at"`
	URL       string   `json:"absolute_url"`
	Company   string   `json:"-"` // Set during polling.
	Verdict   Verdict  `json:"-"` // Set during filtering.
}

type AshbyPosting struct {
	Title    string `json:"title"`
	JobID    string `json:"id"`
	Location string `json:"location"`
	PostedAt string `json:"publishedAt"`
}
