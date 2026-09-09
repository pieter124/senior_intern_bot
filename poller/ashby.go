package poller

import (
	domain "senior_intern_bot/domain"
	"strings"
)

type ashbyJob struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Department         string `json:"department"`
	Team               string `json:"team"`
	EmploymentType     string `json:"employmentType"`
	Location           string `json:"location"`
	SecondaryLocations []struct {
		Location string `json:"location"`
	} `json:"secondaryLocations"`
	PublishedAt string `json:"publishedAt"`
	IsListed    bool   `json:"isListed"`
	URL         string `json:"jobUrl"` // VERIFY this key
}

func (j ashbyJob) toPosting(company string) domain.Posting {
	locations := []string{j.Location}
	for _, s := range j.SecondaryLocations {
		locations = append(locations, s.Location)
	}

	return domain.Posting{
		Source:    "ashby",
		Company:   company,
		JobID:     j.ID,
		Title:     j.Title,
		Location:  strings.Join(locations, ", "),
		URL:       j.URL,
		PostedAt:  j.PublishedAt,
		UpdatedAt: j.PublishedAt, // Ashby exposes no updatedAt
	}
}
