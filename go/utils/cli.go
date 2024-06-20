package utils

import (
	"NicholasSynovic/types"
	"flag"
	"path/filepath"
)

func ParseCommandLine(inputHelpString string, outputHelpString string) types.AppConfig {
	config := types.AppConfig{InputDirectoryPath: ".", OutputJSONFilePath: "output.json"}

	flag.StringVar(&config.InputDirectoryPath, "i", config.InputDirectoryPath, inputHelpString)

	flag.StringVar(&config.OutputJSONFilePath, "o", config.OutputJSONFilePath, outputHelpString)

	flag.Parse()

	config.InputDirectoryPath, _ = filepath.Abs(config.InputDirectoryPath)
	config.OutputJSONFilePath, _ = filepath.Abs(config.OutputJSONFilePath)

	return config
}

func ParseCommandLine_GraphML() types.AppConfig_GraphML {
	config := types.AppConfig_GraphML{InputSQLite3DBPath: "sqlite.db", OutputXMLFilePath: "graph.gml"}

	flag.StringVar(&config.InputSQLite3DBPath, "i", config.InputSQLite3DBPath, "Path to SQLite3 database")

	flag.StringVar(&config.OutputXMLFilePath, "o", config.OutputXMLFilePath, "Path to output GraphML file")

	flag.Parse()

	config.InputSQLite3DBPath, _ = filepath.Abs(config.InputSQLite3DBPath)
	config.OutputXMLFilePath, _ = filepath.Abs(config.OutputXMLFilePath)

	return config
}
