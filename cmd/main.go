package main

import (
	"fmt"

	"github.com/subediaakash/search-engine/internal/document"
	"github.com/subediaakash/search-engine/internal/tokenizer"
)

func main() {
	docs, err := document.LoadDocument("files")
	if err != nil {
		panic(err)
	}
	for _, doc := range docs {

		// print tokens from each docs
		words := tokenizer.Tokenize(doc.Content)
		fmt.Println("Tokens : ", words)
		fmt.Println("------------------")

	}

}
