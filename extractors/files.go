package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

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

func readLines(file *os.File, channel chan string) {
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			if len(line) > 0 {
				channel <- line
			}
			break
		}

		if err != nil {
			fmt.Println("Error reading file:", filepath.Base(file.Name()))
			os.Exit(1)
		}

		channel <- line
	}
	close(channel)
}

func writeJSON(file *os.File, data Output) {
	outputJSON, _ := json.Marshal(data)
	writer := bufio.NewWriter(file)
	writer.Write(outputJSON)
}
