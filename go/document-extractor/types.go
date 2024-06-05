package main

import "time"

/*
Type to represent application configuration

Modified by user command line input
*/
type AppConfig struct {
	inputPath, outputPath string
}

/*
Type of a Document
*/
type Document struct {
	DOI              string    `json:"doi"`
	TITLE            string    `json:"title"`
	PUBLICATION_DATE time.Time `json:"publication_date"`
	OA_TYPE          string    `json:"openalex_type"`
	CR_TYPE          string    `json:"crossref_type"`
	CITED_BY_COUNT   int       `json:"cited_by_count"`
	RETRACTED        bool      `json:"is_retracted"`
	OPEN_ACCESS      bool      `json:"is_open_access"`
}
