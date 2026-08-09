package index

import (
	"github.com/subediaakash/search-engine/internal/document"
	"github.com/subediaakash/search-engine/internal/tokenizer"
)

type InvertedIndex struct {
	Words     map[string][]int
	Documents map[int]document.Document
}

// thsi is to prevent creating map everytime

func New() *InvertedIndex {
	return &InvertedIndex{
		Words:     make(map[string][]int),
		Documents: make(map[int]document.Document),
	}
}

func (idx *InvertedIndex) Build(docs []document.Document) {
	for _, doc := range docs {
		seen := make(map[string]bool)

		idx.Documents[doc.Id] = doc

		words := tokenizer.Analyze(doc.Content)
		for _, eachWord := range words {
			if seen[eachWord] {
				continue
			}
			idx.Words[eachWord] = append(idx.Words[eachWord], doc.Id)
			seen[eachWord] = true
		}
	}
}
