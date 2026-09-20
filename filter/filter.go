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

type axis struct {
	inScope    []string
	outOfScope []string
}
type KeywordFilter struct {
	role       axis
	region     axis
	discipline axis
}

func New() *KeywordFilter {
	return &KeywordFilter{
		role:       RoleAxis,
		region:     RegionAxis,
		discipline: DisciplineAxis,
	}
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

func (a *axis) evaluate(text string) signal {
	hitInScope := util.ContainsAny(text, a.inScope)
	hitOutOfScope := util.ContainsAny(text, a.outOfScope)

	switch {
	case hitInScope && !hitOutOfScope:
		return yes
	case !hitInScope && hitOutOfScope:
		return no
	default:
		return unknown
	}
}
