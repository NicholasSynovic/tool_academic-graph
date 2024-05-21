package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func readFile(fp string) ([]string, int64) {
	var file *os.File
	var err error
	var data []string
	var lineReader *bufio.Reader
	var bytes []byte
	var line string

	file, err = os.Open(fp)
	if err != nil {
		fmt.Println("Error reading", fp)
		os.Exit(1)
	}

	defer file.Close()

	bar := progressbar.Default(-1, ("Reading lines from " + fp))

	lineReader = bufio.NewReader(file)
	for {
		bytes, err = lineReader.ReadBytes('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading file bytes of", fp)
			os.Exit(1)
		} else {
			line = string(bytes)
			line = strings.TrimSpace(line)
			data = append(data, line)
		}

		bar.Add(1)

	}

	return data, int64(len(data))
}
