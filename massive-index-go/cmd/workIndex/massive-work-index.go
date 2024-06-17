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

	spinner := progressbar.Default(-1, "Creating JSON objs...")

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
	return data
}

func WriteWorkIndicesToFile(fp *os.File, data []types.WorkIndex) {
	dataBytes, _ := json.Marshal(data)
	utils.WriteJSONToFile(fp, dataBytes)
}

func main() {
	config := utils.ParseCommandLine(`Path to a directory containing OpenAlex "Works" JSON files`, "Path to a JSON file to store the output")

	if !utils.ValidateInputDirectory(config.InputDirectoryPath) {
		fmt.Printf("%s is not a directory\n", config.InputDirectoryPath)
		os.Exit(1)
	}

	if !utils.ValidateOutputFile(config.OutputJSONFilePath) {
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
