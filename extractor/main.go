package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

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
Wraps around the writeJSON() function to provide output to the command line
*/
func wrapper_WriteJSON(fp string, data interface{}) {
	citesOutputFile := createFile(fp)
	writeJSON(citesOutputFile, data)
	citesOutputFile.Close()
	fmt.Println("Wrote to file:", filepath.Base(fp))
}

func wrapper_HandleJSON(inChannel chan map[string]any, workChannel chan Work, citationChannel chan Citation, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		_, ok := <-inChannel

		if !ok {
			break
		}

		// Create Work objects
		go jsonToWorkObjs(inChannel, workChannel)

		// Create Citation objects

		// Write Work objects

		// Write Citation objects

	}
}

/*
Parse the command line for relevant program flags

Returns (string, string) where the first string is the absolute path of a
SQLite3 database and the second is the absolute path of a text file to output
queries to

On error, calls _printCommandLineParsingError()
*/
func parseCommandLine() (string, string, string, int) {
	var inputPath, worksOutputPath, citesOutputPath string
	var processes int

	flag.StringVar(&inputPath, "i", "", `Path to OpenAlex "Works" JSON Lines file`)
	flag.StringVar(&worksOutputPath, "works-output", "", "Path to output JSON file to store Works information")
	flag.StringVar(&citesOutputPath, "cites-output", "", "Path to output JSON file to store Citation relationship information")
	flag.IntVar(&processes, "proc", runtime.NumCPU(), "Number of processors to use")
	flag.Parse()

	if inputPath == "" {
		_printCommandLineParsingError("i")
	}

	if worksOutputPath == "" {
		_printCommandLineParsingError("works-output")
	}

	if citesOutputPath == "" {
		_printCommandLineParsingError("cites-output")
	}

	absInputPath, _ := filepath.Abs(inputPath)
	absWorksOutputPath, _ := filepath.Abs(worksOutputPath)
	absCitesOutputPath, _ := filepath.Abs(citesOutputPath)

	return absInputPath, absWorksOutputPath, absCitesOutputPath, processes
}

/*
Code that is actually executed within the application
*/
func main() {
	var jsonObjs []map[string]any
	var workOutput WorkOutput
	var citationOutput CitationOutput
	var wg sync.WaitGroup

	// Parse command line
	inputPath, worksOutputPath, citesOutputPath, processes := parseCommandLine()

	// Create channels
	fileChannel := make(chan string)
	jsonObjChannel := make(chan map[string]any)
	workObjChannel := make(chan Work)
	citationObjChannel := make(chan Citation)

	// Read in JSON data
	inputFile := openFile(inputPath)
	defer inputFile.Close()
	go readLines(inputFile, fileChannel)

	/*
		Create JSON objs
		Concurrent channel for converting JSON strings to JSON objs
	*/
	for i := 0; i < processes; i++ {
		wg.Add(1)
		go wrapper_HandleJSON(jsonObjChannel, workObjChannel, citationObjChannel, &wg)
	}

	go createJSONObjs(fileChannel, jsonObjChannel)

	/*
		Create Citation objs
		Concurrent channel for converting JSON objs to Citation objs
	*/
	go jsonToCitationObjs(jsonObjs, citationObjChannel)
	for {
		citationObj, ok := <-citationObjChannel

		if !ok {
			break
		}

		citationOutput = append(citationOutput, citationObj)
	}

	// Write Works data to JSON file
	wrapper_WriteJSON(worksOutputPath, workOutput)

	// Write Citation data to JSON file
	wrapper_WriteJSON(citesOutputPath, citationOutput)

}
