package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// 2: Open JSON file

// 3: Read each line in the JSON file to channel

// 4: For each line in the channel, convert to JSON object
// progress spinner

func parseCommandLine() AppConfig {
	config := AppConfig{inputPath: "oa_works.json", outputPath: "works_output.json"}

	flag.StringVar(&config.inputPath, "i", config.inputPath, `Path to OpenAlex "Works" JSON Lines file`)

	flag.StringVar(&config.outputPath, "o", config.outputPath, `Path to output JSON file to store "Works" information`)

	flag.Parse()

	config.inputPath, _ = filepath.Abs(config.inputPath)
	config.outputPath, _ = filepath.Abs(config.outputPath)

	return config
}

func main() {
	config := parseCommandLine()

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
