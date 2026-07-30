package main

import (
	"fmt"

	"github.com/subediaakash/search-engine/internal/document"
	"github.com/subediaakash/search-engine/internal/index"
)

func main() {

	// index logic
	// For every document

	// Create empty seen map

	// Tokenize the document

	// Remove stop words

	// For every word

	//     If already seen

	//         continue

	//     Append document ID

	// Mark as seen
	idx := index.New()

	documents, err := document.LoadDocument("files")
	if err != nil {
		return
	}
	idx.Build(documents)
	fmt.Println(idx.Words)

}
