package search

// implementation of inverted index.
import (
	"github.com/subediaakash/search-engine/internal/index"
	"github.com/subediaakash/search-engine/internal/tokenizer"
)

// intersect returns the sorted doc IDs present in both a and b.
// TODO: two-pointer walk. a and b are already sorted (Build appends in
// ascending doc.Id order) so no sorting needed here.
func intersect(a, b []int) []int {
	i := 0
	j := 0
	result := []int{}

	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			result = append(result, a[i])
			i++
			j++
		} else if a[i] < b[j] {
			i++
		} else {
			j++
		}
	}

	return result
}

func Search(idx *index.InvertedIndex, query string) []int {
	words := tokenizer.Analyze(query)
	// check for some cheeky users entering "" as query
	if len(words) == 0 {
		return []int{}
	}

	// --- Layer 1: gather ---
	// TODO: replace the single `word := words[0]` lookup below with a loop
	// over ALL of `words`, collecting each term's postings list into a
	// [][]int. Decide what happens when a term isn't in idx.Words at all
	// (strict AND -> the whole query can't match -> return early).

	word := words[0]

	docs, exists := idx.Words[word]
	if !exists {
		return []int{}
	}

	// --- Layer 2: fold ---
	// TODO: once you have `lists [][]int` from Layer 1, seed an accumulator
	// with lists[0] and intersect() it against each remaining list in turn.
	// Return the final accumulator instead of `docs` below.

	return docs
}

func IntersectionSort(idx *index.InvertedIndex, query string) []int {
	// Same tokenizing step as Search: turn "Corona World!" into ["corona", "world"].
	words := tokenizer.Analyze(query)
	if len(words) == 0 {
		return []int{}
	}

	// --- Layer 1: gather ---
	// lists will hold one postings list per query word:
	//   lists[0] = idx.Words["corona"], lists[1] = idx.Words["world"], ...
	// We don't know len(words) in advance in general, but we do know it
	// here, so we can preallocate with make([]T, 0, cap) to avoid repeated
	// reallocation as we append. This is an optimization, not a requirement
	// — `var lists [][]int` followed by append would also work.
	lists := make([][]int, 0, len(words))

	for _, word := range words {
		docs, exists := idx.Words[word]
		if !exists {
			// Strict AND: if even one term matches nothing, the whole
			// query can't match anything either. No point gathering the
			// rest — return immediately.
			return []int{}
		}
		lists = append(lists, docs)
	}

	// --- Layer 2: fold ---
	// Seed the accumulator with the first postings list, then repeatedly
	// intersect it with each remaining list. acc only ever shrinks (or
	// stays the same) as more lists are folded in — it can never grow.
	acc := lists[0]
	for i := 1; i < len(lists); i++ {
		// Early exit: once the accumulator is empty, intersecting it with
		// anything else is still empty, so there's nothing left to gain
		// by looking at the remaining lists.
		if len(acc) == 0 {
			break
		}
		acc = intersect(acc, lists[i])
	}

	return acc
}
