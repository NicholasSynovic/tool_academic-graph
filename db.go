package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

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

func createDBTables(dbConn *sql.DB) {
	var sqlQuery string
	var err error

	sqlQuery = `CREATE TABLE IF NOT EXISTS works (
		oa_id TEXT NOT NULL PRIMARY KEY,
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

func createDBQuery(workObjs []Work) string {
	queries := []string{"BEGIN TRANSACTION;\n"}

	for i := range workObjs {
		queries = append(queries, fmt.Sprintf(`INSERT INTO works (oa_id, doi, title, paratext, retracted, published, oa_type, cf_type) VALUES ("%s", "%s", "%s", %t, %t, "%s", "%s", "%s");`, workObjs[i].OA_ID, workObjs[i].DOI, workObjs[i].Title, workObjs[i].Is_Paratext, workObjs[i].Is_Retracted, workObjs[i].Date_Published, workObjs[i].OA_Type, workObjs[i].CF_Type))
	}
	queries = append(queries, "COMMIT;")

	return strings.Join(queries, "\n")
}

func writeDataToDB(dbConn *sql.DB, sqlQuery string) {
	_, err := dbConn.Exec(sqlQuery)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Error writing data to the database")
		os.Exit(1)
	}
}
