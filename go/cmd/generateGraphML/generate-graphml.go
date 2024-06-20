package main

import (
	"NicholasSynovic/types"
	"NicholasSynovic/utils"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func getWorkRows(dbConn *sql.DB) *sql.Rows {
	query := "SELECT oaid, doi, updated FROM works;"
	return utils.QueryDB(dbConn, query)
}

func getCitesRelationshipRows(dbConn *sql.DB) *sql.Rows {
	query := `SELECT work_oaid, ref_oaid
	FROM relationship_cites
	WHERE EXISTS (
		SELECT 1
		FROM works
		WHERE relationship_cites.work_oaid = works.oaid
	);`
	return utils.QueryDB(dbConn, query)
}

func createNodesFromWorkRows(workRows *sql.Rows, outChannel chan types.Node) {
	workNodeLabel := ":Work"
	counter := 0

	for workRows.Next() {
		var oaid string
		var doi string
		var updated string

		nodeID := fmt.Sprintf("n%d", counter)

		workRows.Scan(&oaid, &doi, &updated)

		oaidData := types.Data{KEY: "oaid", VAL: oaid}
		doiData := types.Data{KEY: "doi", VAL: doi}

		labelData := types.Data{
			KEY: "labels",
			VAL: workNodeLabel}

		updatedData := types.Data{
			KEY: "updated",
			VAL: updated}

		outChannel <- types.Node{
			ID:         nodeID,
			LABELS:     workNodeLabel,
			LABEL_DATA: labelData,
			OAID_DATA:  oaidData,
			DOI_DATA:   doiData,
			UPDATED:    updatedData}

		counter += 1
	}
	close(outChannel)
}

func main() {
	config := utils.ParseCommandLine_GraphML()

	if utils.ValidateFileExistence(config.InputSQLite3DBPath) {
		fmt.Printf("%s is not a file\n", config.InputSQLite3DBPath)
		os.Exit(1)
	}

	if !utils.ValidateFileExistence(config.OutputXMLFilePath) {
		fmt.Printf("%s is a file\n", config.OutputXMLFilePath)
		os.Exit(1)
	}

	outputFP := utils.CreateFile(config.OutputXMLFilePath)
	defer outputFP.Close()

	nodeChannel := make(chan types.Node)
	// edgeChannel := make(chan types.Edge)

	dbConn := utils.ConnectToSQLite3DB(config.InputSQLite3DBPath)
	defer dbConn.Close()

	workRows := getWorkRows(dbConn)
	defer workRows.Close()

	citesRows := getCitesRelationshipRows(dbConn)
	defer citesRows.Close()

	go createNodesFromWorkRows(workRows, nodeChannel)

	for node := range nodeChannel {
		fmt.Println(node.ID)
	}

	// nodeMap := utils.BufferNodes(nodeChannel)
	// nodes := utils.MapToNodeSlice(nodeMap)

	// go utils.WriteEdgesToChannel(rows, nodeMap, edgeChannel)

	// edges := utils.BufferEdges(edgeChannel)

	// graphML := utils.CreateGraphML(nodes, edges)

	// utils.WriteGraphMLToFile(outputFP, graphML)
}
