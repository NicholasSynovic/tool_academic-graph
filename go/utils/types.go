package utils

/*
Type to represent application configuration

Modified by user command line input
*/
type AppConfig struct {
	InputPath, OutputPath string
}

/*
Type of a Citation
*/
type Citation struct {
	SOURCE string `json:"source"`
	DEST   string `json:"dest"`
}
