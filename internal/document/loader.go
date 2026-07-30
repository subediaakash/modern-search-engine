package document

import (
	"os"
	"path/filepath"
)

func LoadDocument(folder string) ([]Document, error) {
	// array that will store the document id , content and  file path so that we can list out what each document has

	var docs []Document

	// we read the folder , so that we can find what files that folder has
	files, err := os.ReadDir(folder)

	if err != nil {
		return nil, err
	}

	// initiaing an id for the first file in the folder
	id := 1
	// looping through the  files in that folder

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// attaching the path inorder to get the value
		path := filepath.Join(folder, file.Name())
		// extract data from that path
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		// now thatdata is extracted we add it in the array of document
		docs = append(docs, Document{
			Id:      id,
			Name:    file.Name(),
			Content: string(data),
		})
		id++
	}
	return docs, nil
}
