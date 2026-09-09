package deduper

import (
	"database/sql"
	"errors"
	"log"
	"senior_intern_bot/domain"
)

type Deduper interface {
	Dedupe([]domain.Posting) ([]domain.Posting, error)
}

type SQLiteDeduper struct {
	dbConnection *sql.DB
}

func New(dbConnection *sql.DB) *SQLiteDeduper {
	return &SQLiteDeduper{
		dbConnection: dbConnection,
	}
}

func (deduper *SQLiteDeduper) Dedupe(postings []domain.Posting) ([]domain.Posting, error) {
	dedupedPostings := make([]domain.Posting, 0, len(postings))

	var query string = `SELECT 1 FROM postings
						WHERE source = ?
						AND company = ?
						AND job_id = ?
						LIMIT 1;`
	for _, p := range postings {
		var exists int
		err := deduper.dbConnection.QueryRow(query, p.Source, p.Company, p.JobID).Scan(&exists)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			dedupedPostings = append(dedupedPostings, p)
		case err != nil:
			log.Printf("error checking posting for duplicate: %s|%s|%s\n", p.Source, p.Company, p.JobID)
			return nil, err
		}
	}

	return dedupedPostings, nil
}
