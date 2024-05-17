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

func readJSONLines(fp string) []string {
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

	fmt.Println("Reading", fp)

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

	}

	return data
}

func createJSONObjs(jsonStrings []string) {
	// var data []string
	var jsonBytes []byte
	var jsonObj map[string]any

	fmt.Println("Converting JSON strings to objects...")
	for i := 0; i < len(jsonStrings); i++ {
		jsonBytes = []byte(jsonStrings[i])
		err := json.Unmarshal(jsonBytes, &jsonObj)

		if err != nil {
			fmt.Println("JSON decode error", i)
			os.Exit(1)
		}
	}

}

func main() {
	// var oaWorksPath, dbPath string
	var oaWorksPath string
	var jsonStrings []string

	oaWorksPath, _ = parseCommandLine()

	jsonStrings = readJSONLines(oaWorksPath)
	createJSONObjs(jsonStrings)

}
