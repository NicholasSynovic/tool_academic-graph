package main

import (
	"NicholasSynovic/utils"
	"fmt"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func readLines(directory string, outChannel chan string) {
	files, fileCount := utils.ListFilesInDirectory(directory)

	barMessage := "Reading files from " + filepath.Base(directory)
	bar := progressbar.Default(int64(fileCount), barMessage)
	defer bar.Exit()

	for idx := range files {
		if !utils.CheckExtension(files[idx], ".json") {
			bar.Add(1)
			continue
		}

		fp := utils.OpenFile(files[idx])

		// Handles closing outChannel
		go utils.ReadLines(fp, outChannel)

		bar.Add(1)
	}
}

func main() {
	config := utils.ParseCommandLine()

	lineChan := make(chan string, 1000000)

	readLines(config.OAWorkJSONDirectoryPath, lineChan)

	_, numberOfWorkObjs := utils.ConvertToWorkObjs(lineChan)

	fmt.Println(numberOfWorkObjs)
}
