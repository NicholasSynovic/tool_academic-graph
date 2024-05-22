package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
)

type AppConfig struct {
	inputPath, worksOutputPath, citesOutputPath string
	processes                                   int
}

/*
Wraps around the writeJSON() function to provide output to the command line
*/
func wrapper_WriteJSON(fp string, data interface{}) {
	citesOutputFile := createFile(fp)
	writeJSON(citesOutputFile, data)
	citesOutputFile.Close()
	fmt.Println("Wrote to file:", filepath.Base(fp))
}

func wrapper_HandleJSON(inChannel chan map[string]any, workChannel chan Work, citationChannel chan Citation, processes int, wg *sync.WaitGroup) {
	defer wg.Done()

	bar := progressbar.Default(-1, "Creating objects...")
	// TODO: Make sure that we are not coupling/ involving states on obj
	for {
		obj, ok := <-inChannel

		if !ok {
			break
		}

		// Create Work objects
		go jsonToWorkObj(obj, workChannel)

		// Create Citation objects
		go jsonToCitationObj(obj, citationChannel)

		bar.Add(processes)
	}
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

func main() {
	var workOutput WorkOutput
	var citationOutput CitationOutput
	var wg sync.WaitGroup

	// Parse command line
	config := parseCommandLine()

	// Create channels
	jsonObjChannel := make(chan map[string]any)
	workObjChannel := make(chan Work)
	citationObjChannel := make(chan Citation)

	// Read in JSON data
	inputFile := openFile(config.inputPath)
	jsonLines := readLines(inputFile)
	inputFile.Close()

	/*
		Create JSON objs
		Concurrent channel for converting JSON strings to JSON objs
	*/
	for i := 0; i < config.processes; i++ {
		wg.Add(1)
		go wrapper_HandleJSON(jsonObjChannel, workObjChannel, citationObjChannel, config.processes, &wg)
	}

	go createJSONObjs(jsonLines, jsonObjChannel)

	wg.Wait()

	os.Exit(1)

	/*
		Create Citation objs
		Concurrent channel for converting JSON objs to Citation objs
	*/

	// Write Works data to JSON file
	wrapper_WriteJSON(config.worksOutputPath, workOutput)

	// Write Citation data to JSON file
	wrapper_WriteJSON(config.citesOutputPath, citationOutput)

}
