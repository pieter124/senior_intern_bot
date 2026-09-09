package domain

import "fmt"

type Verdict int

const (
	REVIEW Verdict = iota
	ACCEPT
	REJECT
)

type Source string

const (
	Greenhouse Source = "greenhouse"
	Ashby      Source = "ashby"
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

type Posting struct {
	Source    string
	Title     string
	JobID     string
	Location  string
	PostedAt  string
	UpdatedAt string
	URL       string
	Company   string
	Verdict   Verdict
}
