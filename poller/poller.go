package poller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	domain "senior_intern_bot/domain"
	"time"
)

type Poller interface {
	Poll() ([]domain.Posting, error)
}

type GreenhousePoller struct {
	greenhouseCompanies []string
	client              *http.Client
}

func getGreenhouseURL(company string) string {
	return fmt.Sprintf("https://boards-api.greenhouse.io/v1/boards/%s/jobs?content=true", company)
}

func New(greenhouseCompanies []string) *GreenhousePoller {
	return &GreenhousePoller{
		greenhouseCompanies: greenhouseCompanies,
		client:              &http.Client{Timeout: 10 * time.Second},
	}
}

func (poller *GreenhousePoller) Poll() ([]domain.Posting, error) {
	var (
		polled []domain.Posting
		errs   []error
	)

	// Poll greenhouse companies
	for _, company := range poller.greenhouseCompanies {
		postings, err := poller.pollGreenhouse(company)
		if err != nil {
			errs = append(errs, fmt.Errorf("greenhouse %s: %w", company, err))
			continue
		}
		polled = append(polled, postings...)
	}
	return polled, errors.Join(errs...)
}

func (poller *GreenhousePoller) pollGreenhouse(company string) ([]domain.Posting, error) {
	resp, err := poller.client.Get(getGreenhouseURL(company))
	if err != nil {
		return nil, fmt.Errorf("error making get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var body struct {
		Jobs []json.RawMessage `json:"jobs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	postings := make([]domain.Posting, 0, len(body.Jobs))
	for _, raw := range body.Jobs {
		var p domain.Posting
		if err := json.Unmarshal(raw, &p); err != nil {
			log.Println("error unmarshalling json into posting ", err)
			continue
		}
		p.Company = company
		p.RawJSON = string(raw)
		postings = append(postings, p)
	}
	return postings, nil
}
