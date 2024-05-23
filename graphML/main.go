package main

import (
	"database/sql"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
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
		fmt.Println(`Error getting the "work" and "reference" column from table "cites":`, err)
		os.Exit(1)
	}

	return rows
}

func writeNodesToChannel(uniqueWorks *sql.Rows, outChannel chan Node) {
	defer uniqueWorks.Close()
	defer close(outChannel)

	for uniqueWorks.Next() {
		var nodeID string
		uniqueWorks.Scan(&nodeID)
		outChannel <- Node{ID: nodeID}
	}
}

func writeEdgesToChannel(rows *sql.Rows, outChannel chan Edge) {
	defer rows.Close()
	defer close(outChannel)

	counter := 0
	for rows.Next() {
		var s, t string

		edgeID := fmt.Sprintf("e%d", counter)

		rows.Scan(&s, &t)
		outChannel <- Edge{
			ID:     edgeID,
			Source: s,
			Target: t,
		}
	}
}

func bufferNodes(nodeChannel chan Node) []Node {
	var nodes []Node

	bar := progressbar.Default(-1, "Buffering nodes...")
	for node := range nodeChannel {
		nodes = append(nodes, node)
		bar.Add(1)
	}
	bar.Exit()

	return nodes
}

func bufferEdges(edgeChannel chan Edge) []Edge {
	var edges []Edge

	bar := progressbar.Default(-1, "Buffering edges...")
	for edge := range edgeChannel {
		edges = append(edges, edge)
		bar.Add(1)
	}
	bar.Exit()

	return edges
}

func createGraphML(nodes []Node, edges []Edge) GraphML {
	graph := Graph{
		ID:          "G",
		Edgedefault: "directed",
		Nodes:       nodes,
		Edges:       edges,
	}

	graphML := GraphML{
		Xmlns:  "http://graphml.graphdrawing.org/xmlns",
		Graphs: []Graph{graph},
	}

	return graphML
}

func writeGraphML(fp string, graphML GraphML) {
	fmt.Println("Writing data to", filepath.Base(fp))
	file, _ := os.Create(fp)
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(graphML); err != nil {
		fmt.Println("Error encoding XML:", err)
		return
	}
}

func main() {
	config := parseCommandLine()

	nodeChannel := make(chan Node)
	edgeChannel := make(chan Edge)

	sqlQuery_GetUniqueWorks := "SELECT DISTINCT work FROM cites"
	sqlQuery_GetRows := `SELECT work, reference
	FROM cites
	WHERE reference IN (
		SELECT work FROM cites
	);`

	dbConn := connectToDatabase(config.inputPath)
	defer dbConn.Close()

	uniqueWorks := queryDB(dbConn, sqlQuery_GetUniqueWorks)
	rows := queryDB(dbConn, sqlQuery_GetRows)

	go writeNodesToChannel(uniqueWorks, nodeChannel)
	go writeEdgesToChannel(rows, edgeChannel)

	nodes := bufferNodes(nodeChannel)
	edges := bufferEdges(edgeChannel)

	grpahML := createGraphML(nodes, edges)

	writeGraphML(config.outputPath, grpahML)
}
