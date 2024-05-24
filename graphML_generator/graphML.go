package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"

	"github.com/schollz/progressbar/v3"
)

type GraphML struct {
	XMLName xml.Name `xml:"graphml"`
	Xmlns   string   `xml:"xmlns,attr"`
	Graphs  []Graph  `xml:"graph"`
}

type Graph struct {
	XMLName     xml.Name `xml:"graph"`
	ID          string   `xml:"id,attr"`
	Edgedefault string   `xml:"edgedefault,attr"`
	Nodes       []Node   `xml:"node"`
	Edges       []Edge   `xml:"edge"`
}

type Node struct {
	XMLName xml.Name `xml:"node"`
	ID      string   `xml:"id,attr"`
	Data    Data     `xml:"data"`
}

type Data struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",chardata"`
}

type Key struct {
	ID   string `xml:"id,attr"`
	For  string `xml:"for,attr"`
	Attr string `xml:"attr.name,attr"`
	Type string `xml:"attr.type,attr"`
}

type Edge struct {
	XMLName xml.Name `xml:"edge"`
	ID      string   `xml:"id,attr"`
	Source  string   `xml:"source,attr"`
	Target  string   `xml:"target,attr"`
}

func writeNodesToChannel(uniqueWorks *sql.Rows, outChannel chan Node) {
	defer uniqueWorks.Close()
	defer close(outChannel)

	counter := 0

	for uniqueWorks.Next() {
		var nodeID string
		uniqueWorks.Scan(&nodeID)
		outChannel <- Node{ID: fmt.Sprintf("n%d", counter), Data: Data{Key: "oa_id", Value: nodeID}}
		counter++
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
