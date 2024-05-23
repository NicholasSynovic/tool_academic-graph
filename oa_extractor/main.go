package main

import (
	"flag"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
)

/*
Parse the command line for relevant program flags

Returns AppConfig
*/
func parseCommandLine() AppConfig {
	config := AppConfig{inputPath: "works.json", worksOutputPath: "works_output.json", citesOutputPath: "citations_output.json"}

	flag.StringVar(&config.citesOutputPath, "c", config.citesOutputPath, "Path to output JSON file to store Citation relationship information")

	flag.StringVar(&config.inputPath, "i", config.inputPath, `Path to OpenAlex "Works" JSON Lines file`)

	flag.StringVar(&config.worksOutputPath, "w", config.worksOutputPath, "Path to output JSON file to store Works information")

	flag.Parse()

	config.inputPath, _ = filepath.Abs(config.inputPath)
	config.citesOutputPath, _ = filepath.Abs(config.citesOutputPath)
	config.worksOutputPath, _ = filepath.Abs(config.worksOutputPath)

	return config
}

func main() {
	// Variable to store objects
	var objs []map[string]any

	// Parse command line
	config := parseCommandLine()

	// Create channels
	jsonLinesChannel := make(chan string)
	jsonObjChannel := make(chan map[string]any)
	workObjChannel := make(chan Work)
	citationObjChannel := make(chan Citation)

	// Read in JSON data to channel
	inputFile := openFile(config.inputPath)
	// defer inputFile.Close()
	go readLines(inputFile, jsonLinesChannel)

	// Create JSON objs
	go createJSONObjs(jsonLinesChannel, jsonObjChannel)

	bar := progressbar.Default(-1, "Collecting JSON objects...")
	for data := range jsonObjChannel {
		objs = append(objs, data)
		bar.Add(1)
	}
	bar.Exit()

	// Create Work objs
	go jsonToWorkObj(objs, workObjChannel)

	// Create Citation objs
	go jsonToCitationObj(objs, citationObjChannel)

	// Write Work objs to file
	writeWorkToFile(config.worksOutputPath, workObjChannel)

	// Write Citations objs to file
	writeCitationToFile(config.citesOutputPath, citationObjChannel)
}
