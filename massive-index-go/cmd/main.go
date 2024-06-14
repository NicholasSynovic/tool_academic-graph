package main

import (
	"NicholasSynovic/types"
	"NicholasSynovic/utils"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func readLines(directory string, outChannel chan types.File_Lines) {
	files, fileCount := utils.ListFilesInDirectory(directory)

	barMessage := "Reading files from " + filepath.Base(directory)
	bar := progressbar.Default(int64(fileCount), barMessage)
	defer bar.Exit()

	for idx := range files {
		filepath := files[idx]
		if !utils.CheckExtension(filepath, ".json") {
			bar.Add(1)
			continue
		}

		fp := utils.OpenFile(filepath)

		// Handles closing outChannel
		go utils.ReadLines(fp, filepath, outChannel)

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

	lineChan := make(chan types.File_Lines, 1000000)

	outputJSONFP := utils.CreateFile(config.OutputJSONPath)
	defer outputJSONFP.Close()

	readLines(config.OAWorkJSONDirectoryPath, lineChan)

	workObjs, _ := utils.ConvertToWorkObjs(lineChan)

	writeToFile(config.OutputJSONPath, workObjs)

}
