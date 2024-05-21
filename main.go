package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
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

/*
Connect to a SQLite3 database

# Returns a connection to the database

On error, exits the application with code 1
*/
func connectToDatabase(dbPath string) *sql.DB {
	var dbConn *sql.DB
	var err error

	dbConn, err = sql.Open("sqlite3", dbPath)

	if err != nil {
		fmt.Println("Could not open database:", dbPath)
		os.Exit(1)
	}

	return dbConn
}

func createTable(dbConn *sql.DB) {
	var sqlQuery string
	var err error

	sqlQuery = `CREATE TABLE IF NOT EXISTS works (
		oa_id TEXT NOT NULL,
		doi TEXT,
		title TEXT,
		paratext BOOL,
		retracted BOOL,
		published DATE,
		oa_type TEXT,
		cf_type TEXT
	);`

	_, err = dbConn.Exec(sqlQuery)

	if err != nil {
		fmt.Println("Error creating table")
		os.Exit(1)
	}
}

func writeDataToDB(dbConn *sql.DB, workObjs []Work, barSize int64) {
	var queries []string

	bar := progressbar.Default(barSize, "Creating SQL queries...")

	for i := 0; i < int(barSize); i++ {
		sqlQuery := fmt.Sprintf(`INSERT INTO works (oa_id, doi, title, paratext, retracted, published, oa_type, cf_type)VALUES (%s, %s, %s, %t, %t, %s, %s, %s);`, workObjs[i].OA_ID, workObjs[i].DOI, workObjs[i].Title, workObjs[i].Is_Paratext, workObjs[i].Is_Retracted, workObjs[i].Date_Published, workObjs[i].OA_Type, workObjs[i].CF_Type)

		queries = append(queries, sqlQuery)

		bar.Add(1)
	}

	bar2 := progressbar.Default(barSize, "Writing data to database...")
	for i := 0; i < len(queries); i++ {
		_, err := dbConn.Exec(queries[i])

		if err != nil {
			fmt.Println(err.Error())
			fmt.Printf("Error writing line %d to database\n", i)
			os.Exit(1)
		}

		bar2.Add(1)
	}

}

func main() {
	jsonFilePath, _ := parseCommandLine()

	// Read in JSON data
	fileTime := time.Now()
	fmt.Println("Reading file:", filepath.Base(jsonFilePath))
	jsonFile := openFile(jsonFilePath)
	jsonLines := readLines(jsonFile)
	jsonFile.Close()
	fmt.Println("Read file:", filepath.Base(jsonFilePath), time.Since(fileTime))

	// Create JSON objs
	jsonTime := time.Now()
	fmt.Println("Creating JSON objs...")
	jsonObjs := createJSONObjs(jsonLines)
	fmt.Println("Created JSON objs", time.Since(jsonTime))

	// Create Work objs
	workTime := time.Now()
	fmt.Println("Converting JSON to Work objs...")
	jsonToWorkObjs(jsonObjs)
	fmt.Println("Created Work objs", time.Since(workTime))

}
