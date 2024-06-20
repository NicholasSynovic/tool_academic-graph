package utils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func ConnectToSQLite3DB(fp string) *sql.DB {
	db, err := sql.Open("sqlite3", fp)

	if err != nil {
		fmt.Println("Error connecting to:", filepath.Base(fp))
		os.Exit(1)
	}

	return db
}

func QueryDB(dbConn *sql.DB, sqlQuery string) *sql.Rows {
	rows, err := dbConn.Query(sqlQuery)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return rows
}
