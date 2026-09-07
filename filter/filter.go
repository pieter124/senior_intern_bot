package filter

import (
	domain "senior_intern_bot/domain"
)

type Filter struct {
}

func New() (Filter, error) {
	return Filter{}, nil
}

// TODO
func (filter *Filter) Filter(data []domain.Posting) ([]domain.Posting, error) {
	return data, nil
}
