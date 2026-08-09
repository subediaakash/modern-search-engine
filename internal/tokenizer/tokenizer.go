package tokenizer

import (
	"strings"
	"unicode"
)

func Tokenize(text string) []string {
	// step 1 : lets convert all the texts to lower case
	text = strings.ToLower(text)

	// step 2 : remove the punctuations
	clean := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			return r
		}
		return ' '

	}, text)

	// step 3 : split the words by white spaces

	words := strings.Fields(clean)

	return words
}

func Analyze(text string) []string {
	words := Tokenize(text)
	filteredWords := RemoveStopWords(words)
	return filteredWords
}
