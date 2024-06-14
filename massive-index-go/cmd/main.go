package main

import (
	"NicholasSynovic/utils"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func main() {
	config := utils.ParseCommandLine()

	lineChan := make(chan string)

	files, fileCount := utils.ListFilesInDirectory(config.OAWorkJSONDirectoryPath)

	bar := progressbar.Default(int64(fileCount), "Reading files from "+filepath.Base(config.OAWorkJSONDirectoryPath))
	for idx := range files {
		if !utils.CheckExtension(files[idx], ".json") {
			bar.Add(1)
			continue
		}

		fp := utils.OpenFile(files[idx])
		go utils.ReadLines(fp, lineChan)

		bar.Add(1)
	}
	bar.Exit()

	spinner := progressbar.Default(-1, "Iterating through lines...")
	for line := range lineChan {
		_ = line
		spinner.Add(1)
	}

}
