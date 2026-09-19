package filter

import (
	domain "senior_intern_bot/domain"
	"testing"
)

func TestKeywordFilterClassify(t *testing.T) {
	var tests = []struct {
		name    string
		posting domain.Posting
		want    domain.Verdict
	}{
		{
			name:    "internship in EMEA region is accepted",
			posting: domain.Posting{Title: "EMEA Technnology Internship", Location: "EMEA"},
			want:    domain.ACCEPT,
		},
		{
			name:    "summer analyst in a European city is accepted",
			posting: domain.Posting{Title: "Summer Analyst 2027", Location: "Berlin"},
			want:    domain.ACCEPT,
		},
		{
			name:    "non-intern role in Europe goes to review",
			posting: domain.Posting{Title: "Software Engineer", Location: "Europe"},
			want:    domain.REVIEW,
		},
		{
			name:    "out-of-scope discipline is rejected",
			posting: domain.Posting{Title: "Marketing Internship", Location: "London"},
			want:    domain.REJECT,
		},
		{
			name:    "software engineering internship in a UK city is accepted",
			posting: domain.Posting{Title: "Software Engineering Internship", Location: "Manchester"},
			want:    domain.ACCEPT,
		},
		{
			name:    "out-of-scope discipline and region is rejected",
			posting: domain.Posting{Title: "Audit Internship", Location: "San Francisco"},
			want:    domain.REJECT,
		},
	}

	filter := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filter.Classify(tt.posting)
			if got != tt.want {
				t.Errorf("Classify(%+v): got %+v, want %+v", tt.posting, got, tt.want)
			}
		})
	}
}
