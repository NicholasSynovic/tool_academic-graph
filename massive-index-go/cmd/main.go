package main

import (
	"NicholasSynovic/utils"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func main() {
	config := utils.ParseCommandLine()

	files, fileCount := utils.ListFilesInDirectory(config.OAWorkJSONDirectoryPath)

	bar := progressbar.Default(int64(fileCount), "Reading files from "+filepath.Base(config.OAWorkJSONDirectoryPath))

	for idx := range files {
		idx += 1
		bar.Add(1)
	}

}
