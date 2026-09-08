package filter

import (
	domain "senior_intern_bot/domain"
	"strings"
	"unicode"
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
	title := normalise(posting.Title)
	location := normalise(posting.Location.Name)

	role := filter.role.evaluate(title)
	region := filter.region.evaluate(location)
	discipline := filter.discipline.evaluate(title)

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

func (a *axis) evaluate(normalised string) signal {
	hitInScope := containsAny(normalised, a.inScope)
	hitOutOfScope := containsAny(normalised, a.outOfScope)

	switch {
	case hitInScope && !hitOutOfScope:
		return yes
	case !hitInScope && hitOutOfScope:
		return no
	default:
		return unknown
	}
}

func normalise(s string) string {
	words := strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	return " " + strings.Join(words, " ") + " "
}

func containsAny(s string, phrases []string) bool {
	for _, p := range phrases {
		if strings.Contains(s, " "+p+" ") {
			return true
		}
	}
	return false
}
