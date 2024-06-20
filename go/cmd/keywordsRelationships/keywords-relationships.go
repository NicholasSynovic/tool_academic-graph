package main

import (
	"NicholasSynovic/types"
	"NicholasSynovic/utils"
	"encoding/json"
	"fmt"
	"os"

	"github.com/schollz/progressbar/v3"
)

func CreateKeywordRelationships(inChannel chan types.FileLine) []types.KeywordRelationship {
	data := []types.KeywordRelationship{}

	idCounter := 0

	spinner := progressbar.Default(-1, "Creating types.KeywordRelationship...")

	for fl := range inChannel {
		var jsonObj map[string]any
		err := json.Unmarshal([]byte(fl.Line), &jsonObj)

		if err != nil {
			panic(err)
		}

		rawOAID, _ := jsonObj["id"].(string)
		oaid := utils.CleanOAID(rawOAID)

		keywords := jsonObj["keywords"].([]interface{})

		for idx := range keywords {
			keywordObject := keywords[idx].(map[string]any)
			rawKeywordOAID := keywordObject["id"].(string)

			keywordOAID := utils.CleanOAID(rawKeywordOAID)

			kr := types.KeywordRelationship{ID: idCounter, KEYWORD_OAID: keywordOAID, WORK_OAID: oaid}

			data = append(data, kr)
			idCounter += 1
		}
		spinner.Add(1)
	}
	return data
}

func WriteKeywordRelationshipsToFile(fp *os.File, data []types.KeywordRelationship) {
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

	cr := CreateKeywordRelationships(flChan)

	WriteKeywordRelationshipsToFile(outputFP, cr)

}