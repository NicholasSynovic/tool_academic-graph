package main

/*
Type to represent application configuration

Modified by user command line input
*/
type AppConfig struct {
	inputPath, outputPath string
}

/*
Type to a pairing between OA IDs and DOIs
*/
type Pair struct {
	OAID string `json:"oaid"`
	DOI  string `json:"doi"`
}
