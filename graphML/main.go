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

func main() {
	connectToDatabase("../test2.db")
}
