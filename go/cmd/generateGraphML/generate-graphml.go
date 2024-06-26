package main

import (
	"NicholasSynovic/types"
	"NicholasSynovic/utils"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
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

		data := []types.Data{}

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

		data = append(data, labelData, oaidData, doiData, updatedData)

		outChannel <- types.Node{
			ID:     nodeID,
			LABELS: workNodeLabel,
			DATA:   data}

		counter += 1
	}
	close(outChannel)
}

func createEdgesFromCitesRows(citesRows *sql.Rows, outChannel chan types.Edge) {
	edgeLabel := "CITES"
	counter := 0

	for citesRows.Next() {
		var work_oaid string
		var ref_oaid string

		data := []types.Data{}

		edgeID := fmt.Sprintf("e%d", counter)

		citesRows.Scan(&work_oaid, &ref_oaid)

		labelData := types.Data{KEY: "label", VAL: edgeLabel}

		data = append(data, labelData)

		outChannel <- types.Edge{
			ID:     edgeID,
			SOURCE: work_oaid,
			TARGET: ref_oaid,
			LABEL:  edgeLabel,
			DATA:   data}

		counter += 1
	}
	close(outChannel)
}

func BufferWorkNodes(inChannel chan types.Node) []types.Node {
	data := []types.Node{}

	spinner := progressbar.Default(-1, "Buffering work nodes...")

	for node := range inChannel {
		data = append(data, node)
		spinner.Add(1)
	}
	spinner.Exit()
	return data
}

func BufferCitesEdges(inChannel chan types.Edge) []types.Edge {
	data := []types.Edge{}

	spinner := progressbar.Default(-1, "Buffering cites edges...")

	for node := range inChannel {
		data = append(data, node)
		spinner.Add(1)
	}
	spinner.Exit()
	return data
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

	graphMLKeys := utils.GenerateGraphMLKeys()

	nodeChannel := make(chan types.Node)
	edgeChannel := make(chan types.Edge)

	dbConn := utils.ConnectToSQLite3DB(config.InputSQLite3DBPath)
	defer dbConn.Close()

	workRows := getWorkRows(dbConn)
	defer workRows.Close()

	citesRows := getCitesRelationshipRows(dbConn)
	defer citesRows.Close()

	go createNodesFromWorkRows(workRows, nodeChannel)

	go createEdgesFromCitesRows(citesRows, edgeChannel)

	workNodes := BufferWorkNodes(nodeChannel)

	citesEdges := BufferCitesEdges(edgeChannel)

	graph := types.Graph{
		ID:           "G",
		EDGE_DEFAULT: "directed",
		NODES:        workNodes,
		EDGES:        citesEdges}

	graphML := types.GraphML{
		Xmlns: "http://graphml.graphdrawing.org/xmlns",
		KEYS:  graphMLKeys,
		GRAPH: graph,
	}

	fmt.Println("Writing data to ", config.OutputXMLFilePath)
	utils.WriteGraphMLToFile(outputFP, graphML)
}
