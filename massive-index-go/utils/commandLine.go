package utils

import (
	"NicholasSynovic/types"
	"flag"
	"os"
	"path/filepath"
)

func ParseCommandLine() types.MassiveWorkIndex_AppConfig {
	config := types.MassiveWorkIndex_AppConfig{InputFilePath: "part_000.json", OutputFilePath: "oa_works_000.json"}

	flag.StringVar(&config.InputFilePath, "i", config.InputFilePath, `Path to OpenAlex "Works" JSON directory`)

	flag.StringVar(&config.OutputFilePath, "o", config.OutputFilePath, "SQLite3 file to write OA Works Index to")

	flag.Parse()

	config.InputFilePath, _ = filepath.Abs(config.InputFilePath)
	config.OutputFilePath, _ = filepath.Abs(config.OutputFilePath)

	testValidInputs(config)

	return config
}

func testValidInputs(config types.MassiveWorkIndex_AppConfig) {
	_, err := os.Stat(config.InputFilePath)
	if err != nil {
		panic(os.ErrNotExist)
	}

	_, err = os.Stat(config.OutputFilePath)
	if err == nil {
		panic(os.ErrExist)
	}
}
