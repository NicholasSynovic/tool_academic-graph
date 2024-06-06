package utils

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

/*
Provide common command line input for executables.

storeObject string:	A string that is uesed for the "-o" flag to add to the help message what object is being written to a JSON file.

returns AppConfig:	An AppConfig struct of the "-i" and "-o" inputs
*/
func ParseCommandLine(storeObject string) AppConfig {
	config := AppConfig{inputPath: "oa_works.json", outputPath: "output.json"}

	flag.StringVar(&config.inputPath, "i", config.inputPath, `Path to OpenAlex "Works" JSON Lines file`)

	flag.StringVar(&config.outputPath, "o", config.outputPath, `Path to output JSON file to store `+storeObject)

	flag.Parse()

	config.inputPath, _ = filepath.Abs(config.inputPath)
	config.outputPath, _ = filepath.Abs(config.outputPath)

	testValidInputs(config)

	return config
}

func testValidInputs(config AppConfig) {
	_, err := os.Stat(config.inputPath)
	if err != nil {
		fmt.Println("ERROR: input doesn't exist:", config.inputPath)
		os.Exit(1)
	}

	_, err = os.Stat(config.outputPath)
	if err == nil {
		fmt.Println("ERROR: output exist:", config.outputPath)
		os.Exit(1)
	}
}
