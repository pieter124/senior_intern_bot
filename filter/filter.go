package filter

import (
	domain "senior_intern_bot/domain"
)

type Filterer interface {
	Filter([]domain.Posting) ([]domain.Posting, error)
}

type KeywordFilter struct {
}

func New() *KeywordFilter {
	return &KeywordFilter{}
}

// TODO
func (filter *KeywordFilter) Filter(data []domain.Posting) ([]domain.Posting, error) {

	return data, nil
}
