package utils

/*
Type to represent application configuration

Modified by user command line input
*/
type AppConfig struct {
	inputPath, outputPath string
}

/*
Type of a Citation
*/
type Citation struct {
	SOURCE string `json:"source"`
	DEST   string `json:"dest"`
}
