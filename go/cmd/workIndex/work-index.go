package main

import (
	"NicholasSynovic/types"
	"NicholasSynovic/utils"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

func CreateWorkIndices(inChannel chan types.FileLine) []types.WorkIndex {
	data := []types.WorkIndex{}

	idCounter := 0

	spinner := progressbar.Default(-1, "Creating types.WorkIndex...")

	for fl := range inChannel {
		var jsonObj map[string]any
		err := json.Unmarshal([]byte(fl.Line), &jsonObj)

		if err != nil {
			panic(err)
		}

		rawOAID, _ := jsonObj["id"].(string)
		oaid := utils.CleanOAID(rawOAID)

		updatedDateString, _ := jsonObj["updated_date"].(string)
		updatedDate, _ := time.Parse("2006-01-02T15:04:05.000000", updatedDateString)

		workIndexObj := types.WorkIndex{
			ID:       idCounter,
			OAID:     oaid,
			UPDATED:  updatedDate,
			FILEPATH: fl.Filepath,
		}

		data = append(data, workIndexObj)

		idCounter += 1
		spinner.Add(1)
	}
	spinner.Exit()
	return data
}

func WriteWorkIndicesToFile(fp *os.File, data []types.WorkIndex) {
	bar := progressbar.Default(int64(len(data)), "Writing to file...")

	for _, record := range data {
		jsonBytes, err := json.Marshal(record)
		if err != nil {
			continue
		}
		jsonBytes = append(jsonBytes, '\n')
		_, err = fp.Write(jsonBytes)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		bar.Add(1)
	}
	bar.Exit()
}

func main() {
	config := utils.ParseCommandLine(`Path to a directory containing OpenAlex "Works" JSON files`, "Path to a JSON file to store the output")

	if !utils.ValidateInputDirectory(config.InputDirectoryPath) {
		fmt.Printf("%s is not a directory\n", config.InputDirectoryPath)
		os.Exit(1)
	}

	if !utils.ValidateFileExistence(config.OutputJSONFilePath) {
		fmt.Printf("%s is a file\n", config.OutputJSONFilePath)
		os.Exit(1)
	}

	outputFP := utils.CreateFile(config.OutputJSONFilePath)

	flChan := make(chan types.FileLine)

	fps := utils.GetFilesOfExt(config.InputDirectoryPath, ".json")

	go utils.ReadLines(fps, flChan)

	wis := CreateWorkIndices(flChan)

	WriteWorkIndicesToFile(outputFP, wis)

}
