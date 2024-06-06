package utils

import (
	"flag"
	"os"
	"path/filepath"
)

/*
Provide common command line input for executables.

storeObject string:	A string that is uesed for the "-o" flag to add to the help message what object is being written to a JSON file.

returns AppConfig:	An AppConfig struct of the "-i" and "-o" inputs
*/
func ParseCommandLine(storeObject string) AppConfig {
	config := AppConfig{InputPath: "oa_works.json", OutputPath: "output.json"}

	flag.StringVar(&config.InputPath, "i", config.InputPath, `Path to OpenAlex "Works" JSON Lines file`)

	flag.StringVar(&config.OutputPath, "o", config.OutputPath, `Path to output JSON file to store `+storeObject)

	flag.Parse()

	config.InputPath, _ = filepath.Abs(config.InputPath)
	config.OutputPath, _ = filepath.Abs(config.OutputPath)

	testValidInputs(config)

	return config
}

/*
Ensure that AppConfig values:

	exist in the file system (InputPath),
	do not exist in the file system (OutputPath)
*/
func testValidInputs(config AppConfig) {
	_, err := os.Stat(config.InputPath)
	if err != nil {
		panic(os.ErrNotExist)
	}

	_, err = os.Stat(config.OutputPath)
	if err == nil {
		panic(os.ErrExist)
	}
}
