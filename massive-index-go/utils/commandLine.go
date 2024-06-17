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
