package sender

import (
	domain "senior_intern_bot/domain"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestDiscordBatchify(t *testing.T) {
	// makePostings builds n postings whose formatted lines are all the same length.
	makePostings := func(n int, title string) []domain.Posting {
		postings := make([]domain.Posting, n)
		for i := range postings {
			postings[i] = domain.Posting{Title: title, Company: "acme", URL: "https://example.com/job"}
		}
		return postings
	}

	// A single formatted line for the postings above is well under the limit, but
	// a long title makes a handful of them overflow one batch.
	longTitle := strings.Repeat("a", 150)
	lineLen := utf8.RuneCountInString(formatPosting(makePostings(1, longTitle)[0]))

	var tests = []struct {
		name         string
		postings     []domain.Posting
		wantBatches  int
		wantPerBatch []int // number of postings expected in each batch
	}{
		{
			name:        "no postings",
			postings:    nil,
			wantBatches: 0,
		},
		{
			name:         "single posting",
			postings:     makePostings(1, "Software Intern"),
			wantBatches:  1,
			wantPerBatch: []int{1},
		},
		{
			name:         "several small postings share a batch",
			postings:     makePostings(5, "Software Intern"),
			wantBatches:  1,
			wantPerBatch: []int{5},
		},
		{
			name:         "overflow splits into multiple batches",
			postings:     makePostings(20, longTitle),
			wantBatches:  2,
			wantPerBatch: nil, // checked via the invariants below
		},
		{
			name: "single posting longer than the limit still gets its own batch",
			postings: []domain.Posting{
				{Title: "x", Company: "acme", URL: strings.Repeat("u", discordMsgLimit+10)},
			},
			wantBatches:  1,
			wantPerBatch: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := discordBatchify(tt.postings)

			if len(got) != tt.wantBatches {
				t.Fatalf("got %d batches, want %d", len(got), tt.wantBatches)
			}
			if tt.wantPerBatch != nil {
				for i, b := range got {
					if len(b.postings) != tt.wantPerBatch[i] {
						t.Errorf("batch %d has %d postings, want %d", i, len(b.postings), tt.wantPerBatch[i])
					}
				}
			}

			// Every posting appears exactly once, in the original order.
			var flattened []domain.Posting
			for _, b := range got {
				flattened = append(flattened, b.postings...)
			}
			if len(flattened) != len(tt.postings) {
				t.Fatalf("batches contain %d postings in total, want %d", len(flattened), len(tt.postings))
			}
			for i := range flattened {
				if flattened[i] != tt.postings[i] {
					t.Errorf("posting %d out of order or altered: got %+v, want %+v", i, flattened[i], tt.postings[i])
				}
			}

			for i, b := range got {
				// Content has one line per posting.
				if lines := strings.Count(b.content, "\n") + 1; lines != len(b.postings) {
					t.Errorf("batch %d has %d lines for %d postings", i, lines, len(b.postings))
				}
				// No multi-posting batch may exceed the Discord limit.
				if len(b.postings) > 1 && utf8.RuneCountInString(b.content) > discordMsgLimit {
					t.Errorf("batch %d is %d runes, over the %d limit", i, utf8.RuneCountInString(b.content), discordMsgLimit)
				}
			}
		})
	}

	// Sanity check that the overflow case really is an overflow: 20 lines must not fit in one message.
	if lineLen*20 < discordMsgLimit {
		t.Fatalf("overflow test is not exercising the limit: %d runes total", lineLen*20)
	}
}
