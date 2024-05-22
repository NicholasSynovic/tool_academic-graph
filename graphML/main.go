package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func connectToDatabase(fp string) *sql.DB {
	db, err := sql.Open("sqlite3", fp)

	if err != nil {
		fmt.Println("Error connecting to:", filepath.Base(fp))
		os.Exit(1)
	}

	return db
}

func getCites(dbConn *sql.DB) {
	sqlQuery := "SELECT work, reference FROM cites"

	result, err := dbConn.Exec(sqlQuery)

	if err != nil {
		fmt.Println(`Error getting the "work" and "reference" column from table "cites":`, err)
		os.Exit(1)
	}

	fmt.Println(result)
}

func main() {
	var dbConn *sql.DB

	// TODO: Make this a command line flag
	absFilepath, _ := filepath.Abs("../test2.db")

	dbConn = connectToDatabase(absFilepath)
	defer dbConn.Close()

	getCites(dbConn)
}
