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

func CreateWorkObjects(inChannel chan types.FileLine) []types.WorkObject {
	ERROR_STRING := "!error"

	data := []types.WorkObject{}

	idCounter := 0

	spinner := progressbar.Default(-1, "Creating types.WorkObject...")

	for fl := range inChannel {
		var jsonObj map[string]any
		err := json.Unmarshal([]byte(fl.Line), &jsonObj)

		if err != nil {
			panic(err)
		}

		// ===IDs===

		rawOAID, _ := jsonObj["id"].(string)
		oaid := utils.CleanOAID(rawOAID)

		rawDOI, _ := jsonObj["doi"].(string)
		doi := utils.CleanDOI(rawDOI)

		// ===Dates===

		updatedDateString, _ := jsonObj["updated_date"].(string)
		updatedDate, _ := time.Parse("2006-01-02T15:04:05.000000", updatedDateString)

		createdDateString, _ := jsonObj["created_date"].(string)
		createdDate, _ := time.Parse("2006-01-02", createdDateString)

		publishedDateString, _ := jsonObj["publication_date"].(string)
		publishedDate, _ := time.Parse("2006-01-02", publishedDateString)

		// ===Authorships===

		authorshipCount := len(jsonObj["authorships"].([]interface{}))

		institutionCount := jsonObj["institutions_distinct_count"].(float64)

		distinctCountryCount := jsonObj["countries_distinct_count"].(float64)

		// ===Categories===

		keywordCount := len(jsonObj["keywords"].([]interface{}))

		grantCount := len(jsonObj["grants"].([]interface{}))

		isParatext := jsonObj["is_paratext"].(bool)
		isRetracted := jsonObj["is_retracted"].(bool)

		language, languageErr := jsonObj["language"].(string)
		if !languageErr {
			language = ERROR_STRING
		}

		license, licenseErr := jsonObj["license"].(string)
		if !licenseErr {
			license = ERROR_STRING
		}

		// ===Publication Metrics===

		citedByCount := jsonObj["cited_by_count"].(float64)
		publicationLocationCount := jsonObj["locations_count"].(float64)

		// ===Document Metrics===

		sdgCount := len(jsonObj["sustainable_development_goals"].([]interface{}))

		title, titleErr := jsonObj["title"].(string)
		if !titleErr {
			title = ERROR_STRING
		}

		oaType := jsonObj["type"].(string)
		crType := jsonObj["type_crossref"].(string)

		// ===Work Object===
		workObject := types.WorkObject{
			// IDs
			ID:   idCounter,
			OAID: oaid,
			DOI:  doi,

			// Dates
			UPDATED:   updatedDate,
			CREATED:   createdDate,
			PUBLISHED: publishedDate,

			// Authorships
			AUTHORSHIP_COUNT:       authorshipCount,
			INSTITUTION_COUNT:      int(institutionCount),
			DISTINCT_COUNTRY_COUNT: int(distinctCountryCount),

			// Categories
			KEYWORD_COUNT: keywordCount,
			GRANT_COUNT:   grantCount,
			IS_PARATEXT:   isParatext,
			IS_RETRACTED:  isRetracted,
			LANGUAGE:      language,
			LICENSE:       license,

			// Publication Metrics
			CITED_BY_COUNT:             int(citedByCount),
			PUBLICATION_LOCATION_COUNT: int(publicationLocationCount),

			// Document Metrics
			SUSTAINABLE_DEVELOPMENT_GOAL_COUNT: sdgCount,
			TITLE:                              title,
			OA_TYPE:                            oaType,
			CR_TYPE:                            crType,

			// Meta
			FILEPATH: fl.Filepath,
		}

		data = append(data, workObject)

		idCounter += 1
		spinner.Add(1)
	}
	return data
}
func WriteWorkObjectsToFile(fp *os.File, data []types.WorkObject) {
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

	if !utils.ValidateFileExistence(config.OutputJSONFilePath) {
		fmt.Printf("%s is a file\n", config.OutputJSONFilePath)
		os.Exit(1)
	}

	outputFP := utils.CreateFile(config.OutputJSONFilePath)

	flChan := make(chan types.FileLine)

	fps := utils.GetFilesOfExt(config.InputDirectoryPath, ".json")

	go utils.ReadLines(fps, flChan)

	wo := CreateWorkObjects(flChan)

	WriteWorkObjectsToFile(outputFP, wo)

}
