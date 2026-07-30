package main

import (
	"fmt"

	"github.com/subediaakash/search-engine/internal/document"
)

func main() {
	docs, err := document.LoadDocument("files")
	if err != nil {
		panic(err)
	}
	for _, doc := range docs {

		fmt.Println("ID:", doc.Id)
		fmt.Println("Name:", doc.Name)
		fmt.Println(doc.Content)
		fmt.Println("----------------")
	}

}
