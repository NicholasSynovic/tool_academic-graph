package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	fmt.Println("Writing to file:", filepath.Base(fp))
	citesOutputFile := createFile(fp)
	writeJSON(citesOutputFile, data)
	citesOutputFile.Close()
	fmt.Println("Wrote to file:", filepath.Base(fp))
}

/*
Parse the command line for relevant program flags

Returns (string, string) where the first string is the absolute path of a
SQLite3 database and the second is the absolute path of a text file to output
queries to

On error, calls _printCommandLineParsingError()
*/
func parseCommandLine() (string, string, string) {
	var inputPath, worksOutputPath, citesOutputPath string

	flag.StringVar(&inputPath, "i", "", `Path to OpenAlex "Works" JSON Lines file`)
	flag.StringVar(&worksOutputPath, "works-output", "", "Path to output JSON file to store Works information")
	flag.StringVar(&citesOutputPath, "cites-output", "", "Path to output JSON file to store Citation relationship information")
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

	return absInputPath, absWorksOutputPath, absCitesOutputPath
}

/*
Code that is actually executed within the application
*/
func main() {
	var jsonLines []string
	var jsonObjs []map[string]any
	var workOutput WorkOutput
	var citationOutput CitationOutput

	inputPath, worksOutputPath, citesOutputPath := parseCommandLine()

	// Read in JSON data
	fileTime := time.Now()
	fmt.Println("Reading file:", filepath.Base(inputPath))
	inputFile := openFile(inputPath)

	// Concurrent process for reading a file
	fileChannel := make(chan string)
	go readLines(inputFile, fileChannel)
	for {
		line, ok := <-fileChannel

		if !ok {
			break
		}

		jsonLines = append(jsonLines, line)
	}

	inputFile.Close()
	fmt.Println("Read file:", filepath.Base(inputPath), time.Since(fileTime))

	// Create JSON objs
	jsonTime := time.Now()
	fmt.Println("Creating JSON objs...")

	// Concurrent channel for converting JSON strings to JSON objs
	jsonObjChannel := make(chan map[string]any)
	go createJSONObjs(jsonLines, jsonObjChannel)
	for {
		obj, ok := <-jsonObjChannel

		if !ok {
			break
		}

		jsonObjs = append(jsonObjs, obj)
	}

	fmt.Println("Created JSON objs", time.Since(jsonTime))

	// Create Work objs
	workTime := time.Now()
	fmt.Println("Converting JSON to Work objs...")

	// Concurrent channel for converting JSON objs to Work objs
	workObjChannel := make(chan Work)
	go jsonToWorkObjs(jsonObjs, workObjChannel)
	for {
		workObj, ok := <-workObjChannel

		if !ok {
			break
		}

		workOutput = append(workOutput, workObj)
	}
	fmt.Println("Created Work objs", time.Since(workTime))

	// Create Citation objs
	citationTime := time.Now()
	fmt.Println("Converting JSON to Citation objs...")

	// Concurrent channel for converting JSON objs to Citation objs
	citationObjChannel := make(chan Citation)
	go jsonToCitationObjs(jsonObjs, citationObjChannel)
	for {
		citationObj, ok := <-citationObjChannel

		if !ok {
			break
		}

		citationOutput = append(citationOutput, citationObj)
	}
	fmt.Println("Created Citation objs", time.Since(citationTime))

	// Write Works data to JSON file
	wrapper_WriteJSON(worksOutputPath, workOutput)

	// Write Citation data to JSON file
	wrapper_WriteJSON(citesOutputPath, citationOutput)

}
