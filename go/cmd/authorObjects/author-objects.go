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

func CreateAuthorObjects(inChannel chan types.FileLine) []types.AuthorObject {
	// ERROR_STRING := "!error"

	data := []types.AuthorObject{}

	idCounter := 0

	spinner := progressbar.Default(-1, "Creating types.AuthorObjects...")

	for fl := range inChannel {
		var jsonObj map[string]any
		err := json.Unmarshal([]byte(fl.Line), &jsonObj)

		if err != nil {
			panic(err)
		}

		// ===IDs===

		rawOAID, _ := jsonObj["id"].(string)
		oaid := utils.CleanOAID(rawOAID)

		rawORCID, _ := jsonObj["orcid"].(string)
		orcid := utils.CleanORCID(rawORCID)

		// ===Dates===

		updatedDateString, _ := jsonObj["updated_date"].(string)
		updatedDate, _ := time.Parse("2006-01-02T15:04:05.000000", updatedDateString)

		createdDateString, _ := jsonObj["created_date"].(string)
		createdDate, _ := time.Parse("2006-01-02", createdDateString)

		// ===Counts===

		affiliationCount := len(jsonObj["affiliations"].([]interface{}))

		citationCount := int(jsonObj["cited_by_count"].(float64))

		worksCount := int(jsonObj["works_count"].(float64))

		// ===Author Information===

		displayName := jsonObj["display_name"].(string)

		// ===Metrics===

		summaryStats := jsonObj["summary_stats"].(map[string]any)

		impactFactor := summaryStats["2yr_mean_citedness"].(float64)

		hIndex := int(summaryStats["h_index"].(float64))

		i10Index := int(summaryStats["i10_index"].(float64))

		// ===Work Object===
		authorObject := types.AuthorObject{
			// IDs
			ID:    idCounter,
			OAID:  oaid,
			ORCID: orcid,

			// Dates
			UPDATED: updatedDate,
			CREATED: createdDate,

			// Counts
			AFFILIATION_COUNT: affiliationCount,
			CITATION_COUNT:    citationCount,
			WORKS_COUNT:       worksCount,

			// Author Information
			DISPLAY_NAME: displayName,

			// Metrics
			IMPACT_FACTOR: impactFactor,
			H_INDEX:       hIndex,
			I10_INDEX:     i10Index,

			// Meta
			FILEPATH: fl.Filepath,
		}

		data = append(data, authorObject)

		idCounter += 1
		spinner.Add(1)
	}
	return data
}
func WriteAuthorObjectsToFile(fp *os.File, data []types.AuthorObject) {
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
	config := utils.ParseCommandLine(`Path to a directory containing OpenAlex "Authors" JSON files`, "Path to a JSON file to store the output")

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

	ao := CreateAuthorObjects(flChan)

	WriteAuthorObjectsToFile(outputFP, ao)

}
