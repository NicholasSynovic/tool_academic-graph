package main

import (
	"NicholasSynovic/types"
	"NicholasSynovic/utils"
	"encoding/json"
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

func writeToFile(filepathString string, data []types.Work_Index) {
	fp := utils.CreateFile(filepathString)
	defer fp.Close()

	filename := filepath.Base(filepathString)

	fmt.Println("Writing to " + filename)

	jsonData, err := json.MarshalIndent(data, "", "    ")

	if err != nil {
		panic(err)
	}

	utils.WriteJSONToFile(fp, jsonData)

	fmt.Println("Wrote to " + filename)
}

func main() {
	config := utils.ParseCommandLine()

	lineChan := make(chan string, 1000000)

	readLines(config.OAWorkJSONDirectoryPath, lineChan)

	workObjs, _ := utils.ConvertToWorkObjs(lineChan)

	writeToFile(config.OutputJSONPath, workObjs)
}
