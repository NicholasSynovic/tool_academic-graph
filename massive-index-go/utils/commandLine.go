package utils

import (
	"NicholasSynovic/types"
	"flag"
	"os"
	"path/filepath"
)

func ParseCommandLine() types.MassiveWorkIndex_AppConfig {
	config := types.MassiveWorkIndex_AppConfig{OAWorkJSONDirectoryPath: ".", OutputJSONPath: "oa_works_index.json"}

	flag.StringVar(&config.OAWorkJSONDirectoryPath, "i", config.OAWorkJSONDirectoryPath, `Path to OpenAlex "Works" JSON directory`)

	flag.StringVar(&config.OutputJSONPath, "o", config.OutputJSONPath, "SQLite3 file to write OA Works Index to")

	flag.Parse()

	config.OAWorkJSONDirectoryPath, _ = filepath.Abs(config.OAWorkJSONDirectoryPath)
	config.OutputJSONPath, _ = filepath.Abs(config.OutputJSONPath)

	testValidInputs(config)

	return config
}

func testValidInputs(config types.MassiveWorkIndex_AppConfig) {
	_, err := os.Stat(config.OAWorkJSONDirectoryPath)
	if err != nil {
		panic(os.ErrNotExist)
	}

	_, err = os.Stat(config.OutputJSONPath)
	if err == nil {
		panic(os.ErrExist)
	}
}
