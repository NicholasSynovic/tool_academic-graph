package main

import "ag/utils"

func main() {
	var jsonObjs []map[string]any

	config := utils.ParseCommandLine("citation relationships")

	inputFP := utils.OpenFile(config.InputPath)
	defer inputFP.Close()

	jsonLinesStringChan := make(chan string)
	jsonObjsChan := make(chan map[string]any)

	go utils.ReadLines(inputFP, jsonLinesStringChan)

}
