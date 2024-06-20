package main

import (
	"NicholasSynovic/types"
	"NicholasSynovic/utils"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config := utils.ParseCommandLine_GraphML()

	if utils.ValidateFileExistence(config.InputSQLite3DBPath) == false{
		fmt.Printf("%s is not a file\n", config.InputSQLite3DBPath)
		os.Exit(1)

	if !utils.ValidateFileExistence(config.OutputXMLFilePath) {
		fmt.Printf("%s is a file\n", config.OutputXMLFilePath)
		os.Exit(1)
	}

	sqlQuery_GetUniqueWorks := "SELECT DISTINCT work_oaid FROM relationship_cites"
	sqlQuery_GetRows := `SELECT work_oaid, ref_oaid
	FROM relationship_cites
	WHERE ref_oaid IN (
		SELECT work_oaid FROM relationship_cites
		);`

	outputFP := utils.CreateFile(config.OutputXMLFilePath)
	defer outputFP.Close()

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

	defer fp.Close()

	utils.WriteGraphMLToFile(fp, graphML)
}
