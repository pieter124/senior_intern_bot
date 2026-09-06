package domain

type Verdict int

const (
	ACCEPT Verdict = iota
	REVIEW
	REJECT
)

type Posting struct {
	Company  string
	JobID    string
	PostedAt string
	RawJSON  string
	Verdict  Verdict
}
