package poller

import (
	"fmt"
	"log"
	"net/http"
	domain "senior_intern_bot/domain"
)

type Poller struct {
	greenhouseCompanies []string
}

func GetGreenhouseURL(company string) string {
	return fmt.Sprintf("https://boards-api.greenhouse.io/v1/boards/%s/jobs", company)
}

func New(greenhouseCompanies []string) (Poller, error) {
	return Poller{
		greenhouseCompanies: greenhouseCompanies,
	}, nil
}

func (poller *Poller) Poll() ([]domain.Posting, error) {
	// Poll greenhouse companies
	for _, company := range poller.greenhouseCompanies {
		resp, err := http.Get(GetGreenhouseURL(company))
		if err != nil {
			log.Println("error making get request: ", err)
		}

	}

}
