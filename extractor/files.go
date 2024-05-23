package main

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
func openFile(fp string) *os.File {
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
func createFile(fp string) *os.File {
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
func readLines(file *os.File, outChannel chan string) {
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

func writeFile(fp string, data []interface{}) {
	outputFile := createFile(fp)
	defer outputFile.Close()

	outputJSON, _ := json.MarshalIndent(data, "", "    ")

	writer := bufio.NewWriter(outputFile)
	writer.Write(outputJSON)

	fmt.Println("Wrote to file:", filepath.Base(fp))
}

func writeWorkToFile(fp string, inChannel chan Work) {
	var output []interface{}

	bar := progressbar.Default(-1, "Collecting Work objs...")

	for data := range inChannel {
		output = append(output, data)
		bar.Add(1)
	}

	writeFile(fp, output)
}

func writeCitationToFile(fp string, inChannel chan Citation) {
	var output []interface{}

	bar := progressbar.Default(-1, "Collecting Work objs...")

	for data := range inChannel {
		output = append(output, data)
		bar.Add(1)
	}

	writeFile(fp, output)
}
