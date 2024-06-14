package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

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
