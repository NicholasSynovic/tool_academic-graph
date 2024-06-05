package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

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

	var jsonObjs []map[string]any

	jsonLinesStringChan := make(chan string)
	jsonObjsChan := make(chan map[string]any)

	inputFP := openFile(config.inputPath)
	defer inputFP.Close()

	go readLines(inputFP, jsonLinesStringChan)
	go createJSONObjs(jsonLinesStringChan, jsonObjsChan)

	bar := progressbar.Default(-1, "Collecting JSON objects...")
	for data := range jsonObjsChan {
		jsonObjs = append(jsonObjs, data)
		bar.Add(1)
	}
	bar.Exit()
}
