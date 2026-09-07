package domain

type Verdict int

const (
	ACCEPT Verdict = iota
	REVIEW
	REJECT
)

type Location struct {
	Name string `json:"name"`
}
type Posting struct {
	Company  string   `json:"company_name"`
	Title    string   `json:"title"`
	JobID    int      `json:"internal_job_id"`
	Location Location `json:"location"`
	PostedAt string   `json:"first_published"`
	URL      string   `json:"absolute_url"`
	RawJSON  string   `json:"-"` // Set during polling.
	Verdict  Verdict  `json:"-"` // Done during filtering.
}
