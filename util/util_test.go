package util

import (
	"testing"
)

func TestContainsAny(t *testing.T) {
	var tests = []struct {
		name    string
		s       string
		phrases []string
		want    bool
	}{
		{
			name:    "whole word is matched",
			s:       "Software Intern",
			phrases: []string{"marketing", "intern"},
			want:    true,
		},
		{
			name:    "substring of a longer word is not matched",
			s:       "International Sales",
			phrases: []string{"intern"},
			want:    false,
		},
		{
			name:    "phrase longer than the word is not matched",
			s:       "intern",
			phrases: []string{"international"},
			want:    false,
		},
		{
			name:    "matching is case-insensitive",
			s:       "Engineer",
			phrases: []string{"intern", "software", "engineer"},
			want:    true,
		},
		{
			name:    "phrase case is ignored",
			s:       "engineer",
			phrases: []string{"Engineer"},
			want:    true,
		},
		{
			name:    "punctuation in the text is ignored",
			s:       "Intern (Summer, 2027)",
			phrases: []string{"summer"},
			want:    true,
		},
		{
			name:    "punctuation in a phrase is ignored",
			s:       "Full Time Co Op",
			phrases: []string{"co-op"},
			want:    true,
		},
		{
			name:    "multi-word phrase is matched",
			s:       "Sophia Antipolis, France",
			phrases: []string{"sophia antipolis"},
			want:    true,
		},
		{
			name:    "words in the wrong order are not matched",
			s:       "engineer software",
			phrases: []string{"software engineer"},
			want:    false,
		},
		{
			name:    "no phrase is contained",
			s:       "summer",
			phrases: []string{"emea", "software", "privacy"},
			want:    false,
		},
		{
			name:    "empty phrase list never matches",
			s:       "summer",
			phrases: nil,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsAny(tt.s, tt.phrases)
			if got != tt.want {
				t.Errorf("ContainsAny(%q, %q) = %v, want %v", tt.s, tt.phrases, got, tt.want)
			}
		})
	}
}
