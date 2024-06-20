package main

import (
	"NicholasSynovic/types"
	"NicholasSynovic/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config := utils.ParseCommandLine_GraphML()

	sqlQuery_GetUniqueWorks := "SELECT DISTINCT work FROM cites"
	sqlQuery_GetRows := `SELECT work, reference
	FROM cites
	WHERE reference IN (
		SELECT work FROM cites
		);`

	nodeChannel := make(chan types.Node)
	edgeChannel := make(chan types.Edge)

	dbConn := utils.ConnectToSQLite3DB(config.InputSQLite3DBPath)
	defer dbConn.Close()

	uniqueWorks := utils.QueryDB(dbConn, sqlQuery_GetUniqueWorks)
	rows := utils.QueryDB(dbConn, sqlQuery_GetRows)

	go utils.WriteNodesToChannel(uniqueWorks, nodeChannel)

	nodeMap := utils.BufferNodes(nodeChannel)
	nodes := utils.MapToNodeSlice(nodeMap)

	go utils.WriteEdgesToChannel(rows, nodeMap, edgeChannel)

	edges := utils.BufferEdges(edgeChannel)

	graphML := utils.CreateGraphML(nodes, edges)

	fp := utils.CreateFile(config.OutputXMLFilePath)
	defer fp.Close()

	utils.WriteGraphMLToFile(fp, graphML)
}
