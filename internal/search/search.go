package search

// implementation of inverted index.
import (
	"github.com/subediaakash/search-engine/internal/index"
	"github.com/subediaakash/search-engine/internal/tokenizer"
)

func Search(idx *index.InvertedIndex, query string) []int {
	words := tokenizer.Tokenize(query)
	// check for some cheeky users entering "" as query
	if len(words) == 0 {
		return []int{}
	}

	word := words[0]

	docs, exists := idx.Words[word]
	if !exists {
		return []int{}
	}

	return docs
}
