package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

/*
Open a file that already exists

Returns a pointer to that file
*/
func OpenFile(fp string) *os.File {
	var file *os.File
	var err error

	file, err = os.Open(fp)

	if err != nil {
		fmt.Println("Error opening:", fp)
		os.Exit(1)
	}

	return file
}

/*
Create a file that does not exist or empty a file that does exist

Returns a pointer to that file
*/
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

/*
Given a file, read each line in it

On error, break
*/
func ReadLines(file *os.File, outChannel chan string) {
	defer close(outChannel)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			if len(line) > 0 {
				outChannel <- line
			}
			break
		}

		if err != nil {
			panic(err)
		}

		outChannel <- line
	}
}

/*
Write JSON data to a file
*/
func writeJSONToFile(outputFP *os.File, data []byte) {
	writer := bufio.NewWriter(outputFP)
	_, err := writer.Write(data)
	if err != nil {
		panic(err)
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}

/*
Write Document objects from a channel into a JSON file
*/
func WriteDocumentObjsToFile(filePath string, inChannel chan Document) {
	var data []Document

	bar := progressbar.Default(-1, "Collecting objs...")

	for document := range inChannel {
		data = append(data, document)
		bar.Add(1)
	}
	bar.Exit()

	outputFP := CreateFile(filePath)
	defer outputFP.Close()

	jsonData, _ := json.MarshalIndent(data, "", "    ")

	writeJSONToFile(outputFP, jsonData)

	fmt.Println("Wrote to file:", filepath.Base(filePath))
}

/*
Write CitationRelationship objects from a channel into a JSON file
*/
func WriteCitationRelationshipObjsToFile(filePath string, inChannel chan CitationRelationship) {
	var data []CitationRelationship

	bar := progressbar.Default(-1, "Collecting objs...")

	for document := range inChannel {
		data = append(data, document)
		bar.Add(1)
	}
	bar.Exit()

	outputFP := CreateFile(filePath)
	defer outputFP.Close()

	jsonData, _ := json.MarshalIndent(data, "", "    ")

	writeJSONToFile(outputFP, jsonData)

	fmt.Println("Wrote to file:", filepath.Base(filePath))
}

/*
Write CitationRelationship objects from a channel into a JSON file
*/
func WriteODPObjsToFile(filePath string, inChannel chan ODP) {
	var data []ODP

	bar := progressbar.Default(-1, "Collecting objs...")

	for pair := range inChannel {
		data = append(data, pair)
		bar.Add(1)
	}
	bar.Exit()

	outputFP := CreateFile(filePath)
	defer outputFP.Close()

	jsonData, _ := json.MarshalIndent(data, "", "    ")

	writeJSONToFile(outputFP, jsonData)

	fmt.Println("Wrote to file:", filepath.Base(filePath))
}
