package filter

import (
	domain "senior_intern_bot/domain"
	util "senior_intern_bot/util"
)

type Filterer interface {
	Filter([]domain.Posting) (map[domain.Verdict][]domain.Posting, error)
}

type signal int

const (
	unknown signal = iota
	yes
	no
)

// Axis is one dimension a posting is judged on, described by keyword phrases
// that put it in scope or out of scope.
type Axis struct {
	InScope    []string
	OutOfScope []string
}

type KeywordFilter struct {
	role       Axis
	region     Axis
	discipline Axis
}

// New builds a KeywordFilter from the given axes.
func New(role, region, discipline Axis) *KeywordFilter {
	return &KeywordFilter{
		role:       role,
		region:     region,
		discipline: discipline,
	}
}

// NewDefault builds a KeywordFilter from the keyword lists in keywords.go.
func NewDefault() *KeywordFilter {
	return New(RoleAxis, RegionAxis, DisciplineAxis)
}

func (filter *KeywordFilter) Classify(posting domain.Posting) domain.Verdict {
	role := filter.role.evaluate(posting.Title)
	region := filter.region.evaluate(posting.Location)
	discipline := filter.discipline.evaluate(posting.Title)

	switch {
	case role == no || region == no || discipline == no:
		return domain.REJECT
	case role == yes && region == yes:
		return domain.ACCEPT
	default:
		return domain.REVIEW
	}
}

func (filter *KeywordFilter) Filter(postings []domain.Posting) (map[domain.Verdict][]domain.Posting, error) {
	postingsByVerdict := make(map[domain.Verdict][]domain.Posting)
	for _, p := range postings {
		verdict := filter.Classify(p)
		postingsByVerdict[verdict] = append(postingsByVerdict[verdict], p)
	}
	return postingsByVerdict, nil
}

func (a *Axis) evaluate(text string) signal {
	hitInScope := util.ContainsAny(text, a.InScope)
	hitOutOfScope := util.ContainsAny(text, a.OutOfScope)

	switch {
	case hitInScope && !hitOutOfScope:
		return yes
	case !hitInScope && hitOutOfScope:
		return no
	default:
		return unknown
	}
}
