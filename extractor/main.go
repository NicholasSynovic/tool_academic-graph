package main

import (
	"flag"

	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
)

/*
Parse the command line for relevant program flags

Returns (string, string) where the first string is the absolute path of a
SQLite3 database and the second is the absolute path of a text file to output
queries to

On error, calls _printCommandLineParsingError()
*/
func parseCommandLine() AppConfig {
	config := AppConfig{inputPath: "works.json", worksOutputPath: "works_output.json", citesOutputPath: "citations_output.json", processes: 1}

	flag.StringVar(&config.inputPath, "i", config.inputPath, `Path to OpenAlex "Works" JSON Lines file`)

	flag.StringVar(&config.worksOutputPath, "works-output", config.worksOutputPath, "Path to output JSON file to store Works information")

	flag.StringVar(&config.citesOutputPath, "cites-output", config.citesOutputPath, "Path to output JSON file to store Citation relationship information")

	flag.IntVar(&config.processes, "proc", config.processes, "Number of processors to use")
	flag.Parse()

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
