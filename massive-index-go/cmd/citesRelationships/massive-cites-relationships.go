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
	"strings"

	"github.com/schollz/progressbar/v3"
)

func ValidateInputDirectory(directory string) bool {
	fi, fiErr := os.Stat(directory)
	if fiErr != nil {
		panic(fiErr)
	}

	return fi.IsDir()
}

func ValidateOutputFile(filepath string) bool {
	_, fiErr := os.Stat(filepath)
	if fiErr != nil {
		return true
	}

	return false
}

func GetFilesOfExt(directory string, ext string) []*os.File {
	data := []*os.File{}

	directoryReader, _ := os.ReadDir(directory)

	for idx := range directoryReader {
		fileName := directoryReader[idx].Name()
		fileExt := filepath.Ext(fileName)

		if fileExt == ext {
			fpString := filepath.Join(directory, fileName)
			fp, openErr := os.Open(fpString)

			if openErr != nil {
				panic(openErr)
			}

			data = append(data, fp)
		}
	}

	return data
}

func ReadLines(fps []*os.File, outChannel chan types.FileLine) {
	for idx := range fps {
		fpString := fps[idx].Name()
		reader := bufio.NewReader(fps[idx])

		for {
			line, err := reader.ReadString('\n')

			if err == io.EOF {
				if len(line) > 0 {
					outChannel <- types.FileLine{Line: line, Filepath: fpString}
				}
				break
			}

			if err != nil {
				panic(err)
			}

			outChannel <- types.FileLine{Line: line, Filepath: fpString}
		}
		fps[idx].Close()
	}
	close(outChannel)
}

func CreateCitesRelationships(inChannel chan types.FileLine) []types.CitesRelationship {
	data := []types.CitesRelationship{}

	idCounter := 0

	spinner := progressbar.Default(-1, "Creating JSON objs...")

	for fl := range inChannel {
		var jsonObj map[string]any
		err := json.Unmarshal([]byte(fl.Line), &jsonObj)

		if err != nil {
			panic(err)
		}

		rawOAID, _ := jsonObj["id"].(string)
		oaid := strings.Replace(rawOAID, "https://openalex.org/", "", -1)

		citedWorks := jsonObj["referenced_works"].([]interface{})

		for idx := range citedWorks {
			rawReferenceOAID := citedWorks[idx].(string)
			referenceOAID := strings.Replace(rawReferenceOAID, "https://openalex.org/", "", -1)

			cr := types.CitesRelationship{ID: idCounter, Work_OAID: oaid, Ref_OAID: referenceOAID}

			data = append(data, cr)
			idCounter += 1
		}
		spinner.Add(1)
	}
	return data
}

func CreateFile(fp string) *os.File {
	var file *os.File
	var err error

	file, err = os.Create(fp)

	if err != nil {
		fmt.Println("Error creating:", fp)
		os.Exit(1)
	}

	return file
}

func WriteCitesRelationshipsToFile(fp *os.File, data []types.CitesRelationship) {
	dataBytes, _ := json.Marshal(data)

	writer := bufio.NewWriter(fp)
	_, err := writer.Write(dataBytes)
	if err != nil {
		panic(err)
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}

func main() {
	config := utils.ParseCommandLine(`Path to a directory containing OpenAlex "Works" JSON files`, "Path to a JSON file to store the output")

	if !ValidateInputDirectory(config.InputDirectoryPath) {
		fmt.Printf("%s is not a directory\n", config.InputDirectoryPath)
		os.Exit(1)
	}

	if !ValidateOutputFile(config.OutputJSONFilePath) {
		fmt.Printf("%s is a file\n", config.OutputJSONFilePath)
		os.Exit(1)
	}

	outputFP := CreateFile(config.OutputJSONFilePath)

	flChan := make(chan types.FileLine)

	fps := GetFilesOfExt(config.InputDirectoryPath, ".json")

	go ReadLines(fps, flChan)

	cr := CreateCitesRelationships(flChan)

	WriteCitesRelationshipsToFile(outputFP, cr)

}
