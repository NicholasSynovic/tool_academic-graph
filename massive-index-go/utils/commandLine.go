package utils

import (
	"NicholasSynovic/types"
	"flag"
	"os"
	"path/filepath"
)

func ParseCommandLine() types.MassiveWorkIndex_AppConfig {
	config := types.MassiveWorkIndex_AppConfig{InputDirectoryPath: ".", OutputFilePath: "oa_works_index.json"}

	flag.StringVar(&config.InputDirectoryPath, "i", config.InputDirectoryPath, `Path to OpenAlex "Works" JSON directory`)

	flag.StringVar(&config.OutputFilePath, "o", config.OutputFilePath, "SQLite3 file to write OA Works Index to")

	flag.Parse()

	config.InputDirectoryPath, _ = filepath.Abs(config.InputDirectoryPath)
	config.OutputFilePath, _ = filepath.Abs(config.OutputFilePath)

	testValidInputs(config)

	return config
}

func testValidInputs(config types.MassiveWorkIndex_AppConfig) {
	_, err := os.Stat(config.InputDirectoryPath)
	if err != nil {
		panic(os.ErrNotExist)
	}

	_, err = os.Stat(config.OutputFilePath)
	if err == nil {
		panic(os.ErrExist)
	}
}
