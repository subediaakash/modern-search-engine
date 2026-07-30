// this file is to remove insignificant words like is, a ,of , in etc.

package tokenizer

var stopWords = map[string]struct{}{
	"a":    {},
	"an":   {},
	"the":  {},
	"is":   {},
	"are":  {},
	"of":   {},
	"to":   {},
	"in":   {},
	"on":   {},
	"for":  {},
	"and":  {},
	"or":   {},
	"with": {},
}

func RemoveStopWords(words []string) []string {
	// variable that will store the final results
	var result []string
	// using nested loop statergy currently  , in future will search a better approach ,i can think of this only at the moment
	// [ TODO: optimise this]
	for _, word := range words {
		if _, exsists := stopWords[word]; exsists {
			continue
		}
		result = append(result, word)
	}

	return result
}
