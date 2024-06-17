package main

import (
	"NicholasSynovic/types"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

func ParseCommandLine() types.AppConfig {
	config := types.AppConfig{InputDirectoryPath: ".", OutputJSONFilePath: "output.json"}

	flag.StringVar(&config.InputDirectoryPath, "i", config.InputDirectoryPath, `Path to OpenAlex "Works" JSON directory`)

	flag.StringVar(&config.OutputJSONFilePath, "o", config.OutputJSONFilePath, "JSON file to write OA Works objects to")

	flag.Parse()

	config.InputDirectoryPath, _ = filepath.Abs(config.InputDirectoryPath)
	config.OutputJSONFilePath, _ = filepath.Abs(config.OutputJSONFilePath)

	return config
}

func ValidateInputDirectory(directory string) bool {
	fi, fiErr := os.Stat(directory)
	if fiErr != nil {
		panic(fiErr)
	}

	return fi.IsDir()
}

func ValidateOutputFile(filepath string) bool {
	_, fiErr := os.Stat(filepath)
	if fiErr != nil {
		return true
	}

	return false
}

func GetFilesOfExt(directory string, ext string) []*os.File {
	data := []*os.File{}

	directoryReader, _ := os.ReadDir(directory)

	for idx := range directoryReader {
		fileName := directoryReader[idx].Name()
		fileExt := filepath.Ext(fileName)

		if fileExt == ext {
			fpString := filepath.Join(directory, fileName)
			fp, openErr := os.Open(fpString)

			if openErr != nil {
				panic(openErr)
			}

			data = append(data, fp)
		}
	}

	return data
}

func ReadLines(fps []*os.File, outChannel chan types.FileLine) {
	for idx := range fps {
		fpString := fps[idx].Name()
		reader := bufio.NewReader(fps[idx])

		for {
			line, err := reader.ReadString('\n')

			if err == io.EOF {
				if len(line) > 0 {
					outChannel <- types.FileLine{Line: line, Filepath: fpString}
				}
				break
			}

			if err != nil {
				panic(err)
			}

			outChannel <- types.FileLine{Line: line, Filepath: fpString}
		}
		fps[idx].Close()
	}
	close(outChannel)
}

func CreateWorkObjects(inChannel chan types.FileLine) []types.WorkObject {
	ERROR_STRING := "!error"

	data := []types.WorkObject{}

	idCounter := 0

	spinner := progressbar.Default(-1, "Creating JSON objs...")

	for fl := range inChannel {
		var jsonObj map[string]any
		err := json.Unmarshal([]byte(fl.Line), &jsonObj)

		if err != nil {
			panic(err)
		}

		// ===IDs===

		rawOAID, _ := jsonObj["id"].(string)
		oaid := strings.Replace(rawOAID, "https://openalex.org/", "", -1)

		rawDOI, _ := jsonObj["doi"].(string)
		doi := strings.Replace(rawDOI, "https://doi.org/", "", -1)

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

		conceptsCount := len(jsonObj["concepts"].([]interface{}))

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
			CONCEPT_COUNT: conceptsCount,
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

func CreateFile(fp string) *os.File {
	var file *os.File
	var err error

	file, err = os.Create(fp)

	if err != nil {
		fmt.Println("Error creating:", fp)
		os.Exit(1)
	}

	return file
}

func WriteWorkObjectsToFile(fp *os.File, data []types.WorkObject) {
	dataBytes, _ := json.Marshal(data)

	writer := bufio.NewWriter(fp)
	_, err := writer.Write(dataBytes)
	if err != nil {
		panic(err)
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}

func main() {
	config := ParseCommandLine()

	if !ValidateInputDirectory(config.InputDirectoryPath) {
		fmt.Printf("%s is not a directory\n", config.InputDirectoryPath)
		os.Exit(1)
	}

	if !ValidateOutputFile(config.OutputJSONFilePath) {
		fmt.Printf("%s is a file\n", config.OutputJSONFilePath)
		os.Exit(1)
	}

	outputFP := CreateFile(config.OutputJSONFilePath)

	flChan := make(chan types.FileLine)

	fps := GetFilesOfExt(config.InputDirectoryPath, ".json")

	go ReadLines(fps, flChan)

	wo := CreateWorkObjects(flChan)

	WriteWorkObjectsToFile(outputFP, wo)

}
