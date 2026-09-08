package storer

import "senior_intern_bot/domain"

type Storer interface {
	Store([]domain.Posting) error
}

type SQLiteStorer struct {
}

func New() *SQLiteStorer {
	// TODO
	return &SQLiteStorer{}
}

func (storer *SQLiteStorer) Store(sentPostings []domain.Posting) error {
	// TODO
	return nil
}
