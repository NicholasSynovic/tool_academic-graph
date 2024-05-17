package test

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	var dbPath, outputFilepath string
	var err error

	flag.StringVar(&dbPath, "i", "", "Path to SQLite3 database")
	flag.StringVar(&outputFilepath, "o", "", "Path to store Neo4J queries")
	flag.Parse()

	if dbPath == "" {
		_printCommandLineParsingError("i")
	}

	if outputFilepath == "" {
		_printCommandLineParsingError("o")
	}

	absDBPath, err := filepath.Abs(dbPath)

	if err != nil {
		fmt.Println("Invalid input: ", dbPath)
		_printCommandLineParsingError("i")
	}

	absOutputFilepath, err := filepath.Abs(outputFilepath)

	if err != nil {
		fmt.Println("Invalid input: ", outputFilepath)
		_printCommandLineParsingError("o")
	}

	return absDBPath, absOutputFilepath
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

/*
Get the relevant rows from the 'cites' table in the SQLite3 database and close the larger *sql.DB connection

Returns (in64, *sql.Rows) where the int64 has the number of rows in the 'cites'
table and the *sql.Rows object contains the relevant information from teh
'cites' table

On error, exits with code 1
*/
func getRelevantRows(dbConn *sql.DB) (int64, *sql.Rows) {
	var rows *sql.Rows
	var err error
	var countQuery, rowQuery string
	var count int

	countQuery = "SELECT COUNT(*) FROM cites"
	err = dbConn.QueryRow(countQuery).Scan(&count)

	if err != nil {
		fmt.Println("Error counting table")
		fmt.Println(countQuery)
		os.Exit(1)
	}

	rowQuery = "SELECT work, reference FROM cites"
	rows, err = dbConn.Query(rowQuery)

	if err != nil {
		fmt.Println("Error retrieving data")
		fmt.Println(rowQuery)
		os.Exit(1)
	}

	dbConn.Close()
	return int64(count), rows
}

/*
From the SQLite3 relevant rows, create Neo4J Cypher queries for creating nodes
and relationships and closes the *sql.Rows object

Returns a []string object containing all of the queries

On error, exits with code 1
*/
func generateQueries(tableSize int64, rows *sql.Rows) []string {
	var query string
	var nodeCounter, referenceCounter int
	var data []string

	var queryFormat string = "MERGE (n%d:Work {oa_id: \"%s\"})-[r%d:Cites]->(n%d:Work {oa_id: \"%s\"})"

	bar := progressbar.Default(int64(tableSize), "Creating Neo4J Cypher queries...")

	for rows.Next() {
		var work, reference string

		if err := rows.Scan(&work, &reference); err != nil {
			fmt.Println("Error scanning data")
			os.Exit(1)
		}

		var referenceNodeCounter int = nodeCounter + 1

		query = fmt.Sprintf(queryFormat, nodeCounter, work, referenceCounter, referenceNodeCounter, reference)

		data = append(data, query)

		nodeCounter++
		referenceCounter++

		bar.Add(1)
	}

	rows.Close()

	return data
}

/*
Write Neo4J Cypher queries to a text file

# Does not return anything

On error, exit with code 1
*/
func writeQueriesToFile(queries []string, filepath string) {
	var file *os.File
	var err error
	var item string

	file, err = os.Create(filepath)

	if err != nil {
		fmt.Println("Error creating:", filepath)
		os.Exit(1)
	}

	defer file.Close()

	bar := progressbar.Default(int64(len(queries)), "Writing Neo4J Cypher queries...")

	for _, item = range queries {
		_, err = file.WriteString(item + "\n")

		if err != nil {
			fmt.Println("Error writing:", item)
			os.Exit(1)
		}

		bar.Add(1)
	}
}

func main() {
	var dbPath, outputFile string
	var tableSize int64
	var dbRows *sql.Rows

	dbPath, outputFile = parseCommandLine()
	var dbConn *sql.DB = connectToDatabase(dbPath)
	tableSize, dbRows = getRelevantRows(dbConn)
	var queries []string = generateQueries(tableSize, dbRows)
	writeQueriesToFile(queries, outputFile)
}
