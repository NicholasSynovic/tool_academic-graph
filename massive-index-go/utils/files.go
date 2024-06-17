package utils

import (
	"NicholasSynovic/types"
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

func ReadLines(fp *os.File, filepath string, outChannel chan types.File_Lines) {
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
