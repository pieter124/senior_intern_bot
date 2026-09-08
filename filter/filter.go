package filter

import (
	domain "senior_intern_bot/domain"
)

type Filterer interface {
	Filter([]domain.Posting) (map[domain.Verdict][]domain.Posting, error)
}

type KeywordFilter struct {
}

func New() *KeywordFilter {
	return &KeywordFilter{}
}

// TODO
func (filter *KeywordFilter) Filter(data []domain.Posting) (map[domain.Verdict][]domain.Posting, error) {
	postingsByVerdictMap := make(map[domain.Verdict][]domain.Posting)
	postingsByVerdictMap[domain.REVIEW] = append(postingsByVerdictMap[domain.ACCEPT], data...)
	return postingsByVerdictMap, nil
}
