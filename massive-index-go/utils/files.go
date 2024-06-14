package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CheckExtension(fp string, extension string) bool {
	return filepath.Ext(fp) == extension
}

func OpenFile(fp string) *os.File {
	var file *os.File
	var err error

	file, err = os.Open(fp)

	if err != nil {
		panic(err)
	}

	return file
}

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

func WriteJSONToFile(outputFP *os.File, data []byte) {
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
