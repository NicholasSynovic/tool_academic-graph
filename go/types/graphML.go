package types

import "encoding/xml"

type GraphML struct {
	XMLName xml.Name `xml:"graphml"`
	Xmlns   string   `xml:"xmlns,attr"`

	KEYS   []Key `xml:"key`
	Graphs Graph `xml:"graph"`
}

type Key struct {
	ID             string `xml:"id,attr"`
	FOR            string `xml:"for,attr"`
	ATTRIBUTE_NAME string `xml:"attr.name,attr"`
}

type Graph struct {
	ID           string `xml:"id,attr"`
	EDGE_DEFAULT string `xml:"edgedefault,attr"`

	NODES []Node `xml:"node"`
	EDGES []Edge `xml:"edge"`
}

type Node struct {
	ID     string `xml:"id,attr"`
	LABELS string `xml:"labels,attr"`

	DATA []Data `xml:"data"`
}

type Edge struct {
	ID     string `xml:"id,attr"`
	SOURCE string `xml:"source,attr"`
	TARGET string `xml:"target,attr"`
	LABEL  string `xml:"label,attr"`

	DATA []Data `xml:"data"`
}

type Data struct {
	KEY string `xml:"key,attr"`
	VAL string `xml:",chardata"`
}
