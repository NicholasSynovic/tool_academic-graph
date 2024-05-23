package main

import "encoding/xml"

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
}

// Edge structure
type Edge struct {
	XMLName xml.Name `xml:"edge"`
	ID      string   `xml:"id,attr"`
	Source  string   `xml:"source,attr"`
	Target  string   `xml:"target,attr"`
}
