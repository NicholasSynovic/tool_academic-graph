package main

import (
	"NicholasSynovic/types"
	"NicholasSynovic/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func readLines(fp *os.File, outChannel chan types.File_Lines) {
	filepath := fp.Name()
	reader := bufio.NewReader(fp)

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			if len(line) > 0 {
				outChannel <- types.File_Lines{Line: line, Filepath: filepath}
			}
			break
		}

		if err != nil {
			panic(err)
		}

		outChannel <- types.File_Lines{Line: line, Filepath: filepath}
	}

	close(outChannel)
}

func writeToFile(filepathString string, data []types.Work_Index) {
	fp := utils.CreateFile(filepathString)
	defer fp.Close()

	filename := filepath.Base(filepathString)

	jsonData, err := json.MarshalIndent(data, "", "    ")

	if err != nil {
		panic(err)
	}

	utils.WriteJSONToFile(fp, jsonData)

	fmt.Println("Wrote to " + filename)
}

func main() {
	config := utils.ParseCommandLine()

	lineChan := make(chan types.File_Lines)

	inputFP := utils.OpenFile(config.InputFilePath)
	defer inputFP.Close()

	outputJSONFP := utils.CreateFile(config.OutputFilePath)
	defer outputJSONFP.Close()

	go readLines(inputFP, lineChan)

	workObjs, _ := utils.ConvertToWorkObjs(lineChan)

	writeToFile(config.OutputFilePath, workObjs)
}
