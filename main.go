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

func main() {
	var jsonLines []string
	var jsonObjs []map[string]any
	var workObjs []Work

	jsonFilePath, dbFilePath := parseCommandLine()

	// Connect to SQLite3 database and create tables
	dbConn := connectToDatabase(dbFilePath)
	defer dbConn.Close()
	createDBTables(dbConn)

	// Read in JSON data
	fileTime := time.Now()
	fmt.Println("Reading file:", filepath.Base(jsonFilePath))
	jsonFile := openFile(jsonFilePath)

	// Concurrent process for reading a file
	fileChannel := make(chan string, 100000)
	go readLines(jsonFile, fileChannel)
	for {
		line, ok := <-fileChannel

		if !ok {
			break
		}

		jsonLines = append(jsonLines, line)
	}

	jsonFile.Close()
	fmt.Println("Read file:", filepath.Base(jsonFilePath), time.Since(fileTime))

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
		obj, ok := <-workObjChannel

		if !ok {
			break
		}

		workObjs = append(workObjs, obj)
	}
	fmt.Println("Created Work objs", time.Since(workTime))

	// Write data to SQLite database
	sqlQuery := createDBQuery(workObjs)
	writeDataToDB(dbConn, sqlQuery)
}
