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
			fmt.Println("Error reading file:", filepath.Base(file.Name()), err)
			break
		}

		outChannel <- line
	}
}

/*
Write Document objects from a channel into a JSON file
*/
func WriteObjsToFile(fp string, inChannel chan interface{}) {
	var output []interface{}

	bar := progressbar.Default(-1, "Collecting objs...")

	for data := range inChannel {
		output = append(output, data)
		bar.Add(1)
	}
	bar.Exit()

	outputFP := CreateFile(fp)
	defer outputFP.Close()

	outputJSON, _ := json.MarshalIndent(output, "", "    ")

	writer := bufio.NewWriter(outputFP)
	_, err := writer.Write(outputJSON)
	if err != nil {
		panic(err)
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}

	fmt.Println("Wrote to file:", filepath.Base(fp))
}
