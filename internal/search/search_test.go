package search

import (
	"testing"

	"github.com/subediaakash/search-engine/internal/document"
	"github.com/subediaakash/search-engine/internal/index"
)

// buildTestIndex gives every test a small, fixed corpus so results don't
// change when the files/ folder changes.
func buildTestIndex() *index.InvertedIndex {
	idx := index.New()
	idx.Build([]document.Document{
		{Id: 1, Name: "corona.txt", Content: "corona virus is a deadly disease in the world"},
		{Id: 2, Name: "family.txt", Content: "you should spend time with your family"},
		{Id: 3, Name: "someone.txt", Content: "be someone who people admire in this world"},
	})
	return idx
}

func TestSearch(t *testing.T) {
	idx := buildTestIndex()

	tests := []struct {
		name  string
		query string
		want  []int
	}{
		{"single term in one doc", "corona", []int{1}},
		{"term shared by two docs", "world", []int{1, 3}},
		{"term not in corpus", "elephant", nil},
		{"empty query", "", nil},

		// The bug we just fixed: the query used to keep "the", which is
		// stripped at index time, so nothing ever matched.
		{"stop word before real term", "the corona", []int{1}},

		// Analyze lowercases and strips punctuation on both sides.
		{"uppercase query", "CORONA", []int{1}},
		{"query with punctuation", "corona!!!", []int{1}},

		// A query made only of stop words analyzes to zero tokens.
		{"only stop words", "the a of", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Search(idx, tt.query)
			if !equal(got, tt.want) {
				t.Errorf("Search(%q) = %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}

// TestSearchMultiTermAND documents Phase 1's remaining work: a two-term query
// should return only documents containing BOTH terms. Today Search looks at
// words[0] only, so it returns every doc matching just the first term.
func TestSearchMultiTermAND(t *testing.T) {
	t.Skip("not implemented yet: Search only uses the first query term")

	idx := buildTestIndex()

	// "world" is in docs 1 and 3; "corona" is only in doc 1.
	// AND semantics means the answer is doc 1 alone.
	got := Search(idx, "corona world")
	if !equal(got, []int{1}) {
		t.Errorf("Search(%q) = %v, want [1]", "corona world", got)
	}
}

// equal compares two doc-ID slices, treating nil and empty as the same thing.
func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
