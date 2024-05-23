package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func connectToDatabase(fp string) *sql.DB {
	db, err := sql.Open("sqlite3", fp)

	if err != nil {
		fmt.Println("Error connecting to:", filepath.Base(fp))
		os.Exit(1)
	}

	return db
}

func queryDB(dbConn *sql.DB, sqlQuery string) *sql.Rows {
	fmtQuery := strings.ReplaceAll(sqlQuery, "\t", "")
	fmtQuery = strings.ReplaceAll(fmtQuery, "\n", " ")
	fmtQuery = strings.TrimSpace(fmtQuery)

	fmt.Println("SQL Query:", fmtQuery)

	rows, err := dbConn.Query(sqlQuery)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return rows
}
