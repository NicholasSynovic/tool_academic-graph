package utils

import (
	"os"
	"path/filepath"
)

func ListFilesInDirectory(directory string) ([]string, int) {
	files := []string{}

	directoryReader, err := os.ReadDir(directory)
	if err != nil {
		panic(err)
	}

	for _, file := range directoryReader {
		filename := file.Name()
		fp, _ := filepath.Abs(filename)
		files = append(files, fp)
	}

	return files, len(files)
}
