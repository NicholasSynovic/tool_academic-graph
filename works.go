package main

import (
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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Work struct {
	OA_ID          string
	DOI            string
	Title          string
	Is_Paratext    bool
	Is_Retracted   bool
	Date_Published time.Time
	OA_Type        string
	CF_Type        string
}

/*
Print common error that can occur during command line parsing if the user does
not input the required flags

Does not return anything and exits with code 1
*/
func _printCommandLineParsingError(parameter string) {
	var errorString string = "-%s is required\n"
	fmt.Printf(errorString, parameter)
	os.Exit(1)
}

/*
Parse the command line for relevant program flags

Returns (string, string) where the first string is the absolute path of a
SQLite3 database and the second is the absolute path of a text file to output
queries to

On error, calls _printCommandLineParsingError()
*/
func parseCommandLine() (string, string) {
	var oaWorksPath, dbPath string
	var err error

	flag.StringVar(&oaWorksPath, "i", "", "Path to OpenAlex 'Works' JSON Lines file")
	flag.StringVar(&dbPath, "o", "", "Path to SQLite3 database")
	flag.Parse()

	if oaWorksPath == "" {
		_printCommandLineParsingError("i")
	}

	if dbPath == "" {
		_printCommandLineParsingError("o")
	}

	absOAWorksPath, err := filepath.Abs(oaWorksPath)

	if err != nil {
		fmt.Println("Invalid input: ", oaWorksPath)
		_printCommandLineParsingError("i")
	}

	absDBPath, err := filepath.Abs(dbPath)

	if err != nil {
		fmt.Println("Invalid input: ", dbPath)
		_printCommandLineParsingError("o")
	}

	return absOAWorksPath, absDBPath
}

func readJSONLines(fp string) ([]string, int64) {
	var file *os.File
	var err error
	var data []string
	var lineReader *bufio.Reader
	var bytes []byte
	var line string

	file, err = os.Open(fp)
	if err != nil {
		fmt.Println("Error reading", fp)
		os.Exit(1)
	}

	defer file.Close()

	bar := progressbar.Default(-1, ("Reading lines from " + fp))

	lineReader = bufio.NewReader(file)
	for {
		bytes, err = lineReader.ReadBytes('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading file bytes of", fp)
			os.Exit(1)
		} else {
			line = string(bytes)
			line = strings.TrimSpace(line)
			data = append(data, line)
		}

		bar.Add(1)

	}

	return data, int64(len(data))
}

func createJSONObjs(jsonStrings []string, barSize int64) []map[string]any {
	var data []map[string]any
	var jsonBytes []byte

	bar := progressbar.Default(barSize, "Converting JSON strings to objects...")

	for i := 0; i < len(jsonStrings); i++ {
		var jsonObj map[string]any

		jsonBytes = []byte(jsonStrings[i])
		err := json.Unmarshal(jsonBytes, &jsonObj)

		if err != nil {
			fmt.Println("JSON decode error", i)
			os.Exit(1)
		}

		data = append(data, jsonObj)

		bar.Add(1)
	}

	return data
}

func createWorkArray(jsonObjs []map[string]any, barSize int64) []Work {
	var data []Work
	var jsonObj map[string]any

	caser := cases.Title(language.AmericanEnglish)

	bar := progressbar.Default(barSize, "Creating an array of Work objects...")

	for i := 0; i < len(jsonObjs); i++ {
		jsonObj = jsonObjs[i]

		id := strings.Replace(jsonObj["id"].(string), "https://openalex.org/", "", -1)

		doi, ok := jsonObj["doi"].(string)
		if !ok {
			bar.Add(1)
			continue
		}
		doi = strings.Replace(doi, "https://doi.org/", "", -1)

		title, ok := jsonObj["title"].(string)
		if !ok {
			bar.Add(1)
			continue
		}
		title = caser.String(title)

		publishedDateString, _ := jsonObj["publication_date"].(string)
		publishedDate, _ := time.Parse("2006-01-02", publishedDateString)

		workObj := Work{
			OA_ID:          id,
			DOI:            doi,
			Title:          title,
			Is_Paratext:    jsonObj["is_paratext"].(bool),
			Is_Retracted:   jsonObj["is_retracted"].(bool),
			Date_Published: publishedDate,
			OA_Type:        jsonObj["type"].(string),
			CF_Type:        jsonObj["type_crossref"].(string),
		}

		data = append(data, workObj)

		bar.Add(1)
	}

	return data
}

func main() {
	// var oaWorksPath, dbPath string
	var oaWorksPath string
	var jsonStrings []string
	var jsonObjs []map[string]any
	var jsonStringsCount int64
	var workObjs []Work

	oaWorksPath, _ = parseCommandLine()

	jsonStrings, jsonStringsCount = readJSONLines(oaWorksPath)
	jsonObjs = createJSONObjs(jsonStrings, jsonStringsCount)
	workObjs = createWorkArray(jsonObjs, int64(len(jsonObjs)))

	for i := 0; i < len(workObjs); i++ {
		fmt.Println(workObjs[i])
	}
}
