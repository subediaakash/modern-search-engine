package index

import (
	"github.com/subediaakash/search-engine/internal/document"
	"github.com/subediaakash/search-engine/internal/tokenizer"
)

type InvertedIndex struct {
	Words map[string][]int
}

// thsi is to prevent creating map everytime

func New() *InvertedIndex {
	return &InvertedIndex{
		Words: make(map[string][]int),
	}
}

func (idx *InvertedIndex) Build(docs []document.Document) {
	for _, docs := range docs {
		seen := make(map[string]bool)

		words := tokenizer.Tokenize(docs.Content)
		filteredWords := tokenizer.RemoveStopWords(words)
		for _, eachWord := range filteredWords {
			if seen[eachWord] {
				continue
			}
			idx.Words[eachWord] = append(idx.Words[eachWord], docs.Id)
			seen[eachWord] = true
		}
	}
}
