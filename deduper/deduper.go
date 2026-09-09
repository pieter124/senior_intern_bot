package deduper

import (
	"senior_intern_bot/domain"
)

type Deduper interface {
	Dedupe([]domain.Posting) ([]domain.Posting, error)
}

type SQLite3Deduper struct {
	// dbConnection

}

func New() *SQLite3Deduper {
	return &SQLite3Deduper{}
}

func (deduper *SQLite3Deduper) Dedupe(postings []domain.Posting) ([]domain.Posting, error) {
	return nil, nil
}
