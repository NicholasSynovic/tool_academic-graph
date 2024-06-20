package types

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
	XML_NAME xml.Name `xml:"node"`
	ID       string   `xml:"id,attr"`

	LABELS     string `xml:"labels,attr"`
	LABEL_DATA Data   `xml:"data"`

	OAID_DATA Data `xml:"data"`
	DOI_DATA  Data `xml:"data"`
	UPDATED   Data `xml:"data"`
}

type Data struct {
	KEY string `xml:"key,attr"`
	VAL string `xml:",chardata"`
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
