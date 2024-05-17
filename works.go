package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
	var oaWorksPath, dbPath string

	oaWorksPath, dbPath = parseCommandLine()
	fmt.Println(oaWorksPath, dbPath)
}
