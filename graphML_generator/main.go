package main

import (
	"flag"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

/*
Parse the command line for relevant program flags

Returns AppConfig
*/
func parseCommandLine() AppConfig {
	config := AppConfig{inputPath: "sqlite.db", outputPath: "graph.gml"}

	flag.StringVar(&config.inputPath, "i", config.inputPath, "Path to SQLite3 database")

	flag.StringVar(&config.outputPath, "o", config.outputPath, "Path to output GraphML file")

	flag.Parse()

	config.inputPath, _ = filepath.Abs(config.inputPath)
	config.outputPath, _ = filepath.Abs(config.outputPath)

	return config
}

func main() {
	config := parseCommandLine()

	sqlQuery_GetUniqueWorks := "SELECT DISTINCT work FROM cites"
	sqlQuery_GetRows := `SELECT work, reference
	FROM cites
	WHERE reference IN (
		SELECT work FROM cites
		);`

	nodeChannel := make(chan Node)
	edgeChannel := make(chan Edge)

	dbConn := connectToDatabase(config.inputPath)
	defer dbConn.Close()

	uniqueWorks := queryDB(dbConn, sqlQuery_GetUniqueWorks)
	rows := queryDB(dbConn, sqlQuery_GetRows)

	go writeNodesToChannel(uniqueWorks, nodeChannel)

	nodeMap := bufferNodes(nodeChannel)
	nodes := mapToNodeSlice(nodeMap)

	go writeEdgesToChannel(rows, nodeMap, edgeChannel)

	edges := bufferEdges(edgeChannel)

	graphML := createGraphML(nodes, edges)

	writeGraphML(config.outputPath, graphML)
}
