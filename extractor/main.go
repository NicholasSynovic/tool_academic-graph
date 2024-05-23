package main

import (
	"flag"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
)

type AppConfig struct {
	inputPath, worksOutputPath, citesOutputPath string
	processes                                   int
}

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

func _writeToFile(fp string, data []interface{}) {
	outputFile := createFile(fp)
	defer outputFile.Close()
	writeJSON(outputFile, data)
	fmt.Println("Wrote to file:", filepath.Base(fp))

}

func writeWorkToFile(fp string, inChannel chan Work) {
	var output []interface{}

	bar := progressbar.Default(-1, "Collecting Work objs...")

	for data := range inChannel {
		output = append(output, data)
		bar.Add(1)
	}

	_writeToFile(fp, output)
}

func writeCitationToFile(fp string, inChannel chan Citation) {
	var output []interface{}

	bar := progressbar.Default(-1, "Collecting Work objs...")

	for data := range inChannel {
		output = append(output, data)
		bar.Add(1)
	}

	_writeToFile(fp, output)
}

func main() {
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

	// Create Work objs
	go jsonToWorkObj(jsonObjChannel, workObjChannel)

	// Write Work objs to file
	writeWorkToFile(config.worksOutputPath, workObjChannel)

	// Write Citations objs to file
	writeCitationToFile(config.worksOutputPath, citationObjChannel)
}
