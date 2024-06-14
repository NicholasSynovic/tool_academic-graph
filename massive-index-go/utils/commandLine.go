package utils

import (
	"NicholasSynovic/types"
	"flag"
	"os"
	"path/filepath"
)

func ParseCommandLine() types.MassiveWorkIndex_AppConfig {
	config := types.MassiveWorkIndex_AppConfig{OAWorkJSONDirectoryPath: ".", OutputDBPath: "oa_works_index.db"}

	flag.StringVar(&config.OAWorkJSONDirectoryPath, "i", config.OAWorkJSONDirectoryPath, `Path to OpenAlex "Works" JSON directory`)

	flag.StringVar(&config.OutputDBPath, "o", config.OutputDBPath, "SQLite3 file to write OA Works Index to")

	flag.Parse()

	config.OAWorkJSONDirectoryPath, _ = filepath.Abs(config.OAWorkJSONDirectoryPath)
	config.OutputDBPath, _ = filepath.Abs(config.OutputDBPath)

	testValidInputs(config)

	return config
}

func testValidInputs(config types.MassiveWorkIndex_AppConfig) {
	_, err := os.Stat(config.OAWorkJSONDirectoryPath)
	if err != nil {
		panic(os.ErrNotExist)
	}

	_, err = os.Stat(config.OutputDBPath)
	if err == nil {
		panic(os.ErrExist)
	}
}
