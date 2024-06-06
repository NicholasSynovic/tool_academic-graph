package main

import (
	"ag/utils"

	"github.com/schollz/progressbar/v3"
)

func main() {
	var jsonObjs []map[string]any

	config := utils.ParseCommandLine("citation relationships")

	inputFP := utils.OpenFile(config.InputPath)
	defer inputFP.Close()

	jsonLinesStringChan := make(chan string)
	jsonObjsChan := make(chan map[string]any)
	documentObjsChan := make(chan utils.Document)

	go utils.ReadLines(inputFP, jsonLinesStringChan)
	go utils.CreateJSONObjs(jsonLinesStringChan, jsonObjsChan)

	bar := progressbar.Default(-1, "Collecting JSON objects...")
	for data := range jsonObjsChan {
		jsonObjs = append(jsonObjs, data)
		bar.Add(1)
	}
	bar.Exit()

	go utils.JSONToDocumentObj(jsonObjs, documentObjsChan)

	utils.WriteDocumentObjsToFile(config.OutputPath, documentObjsChan)
}
