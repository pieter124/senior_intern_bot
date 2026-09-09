package poller

import (
	domain "senior_intern_bot/domain"
	"strconv"
)

type greenhouseJob struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"absolute_url"`
	UpdatedAt string `json:"updated_at"`
	PostedAt  string `json:"first_published"`
	Location  struct {
		Name string `json:"name"`
	} `json:"location"`
}

func (j greenhouseJob) toPosting(company string) domain.Posting {
	return domain.Posting{
		Source:    "greenhouse",
		Company:   company,
		JobID:     strconv.Itoa(j.ID),
		Title:     j.Title,
		Location:  j.Location.Name,
		URL:       j.URL,
		PostedAt:  j.PostedAt,
		UpdatedAt: j.UpdatedAt,
	}
}
