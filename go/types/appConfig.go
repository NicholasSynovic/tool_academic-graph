package types

type AppConfig struct {
	InputDirectoryPath, OutputJSONFilePath string
}

type AppConfig_GraphML struct {
	InputSQLite3DBPath, OutputXMLFilePath string
}
