package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

type AppConfig struct {
	InputDirectoryPath, OutputJSONFilePath string
}

type FileLine struct {
	Line, Filepath string
}

func ParseCommandLine() AppConfig {
	config := AppConfig{InputDirectoryPath: ".", OutputJSONFilePath: "output.json"}

	flag.StringVar(&config.InputDirectoryPath, "i", config.InputDirectoryPath, `Path to OpenAlex "Works" JSON directory`)

	flag.StringVar(&config.OutputJSONFilePath, "o", config.OutputJSONFilePath, "JSON file to write OA Works Index to")

	flag.Parse()

	config.InputDirectoryPath, _ = filepath.Abs(config.InputDirectoryPath)
	config.OutputJSONFilePath, _ = filepath.Abs(config.OutputJSONFilePath)

	return config
}

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

func ReadLines(fps []*os.File, outChannel chan FileLine) {
	for idx := range fps {
		fpString := fps[idx].Name()
		reader := bufio.NewReader(fps[idx])

		for {
			line, err := reader.ReadString('\n')

			if err == io.EOF {
				if len(line) > 0 {
					outChannel <- FileLine{Line: line, Filepath: fpString}
				}
				break
			}

			if err != nil {
				panic(err)
			}

			outChannel <- FileLine{Line: line, Filepath: fpString}
		}
		fps[idx].Close()
	}
	close(outChannel)
}

/*
Steps
*/

func main() {
	config := ParseCommandLine()

	if !ValidateInputDirectory(config.InputDirectoryPath) {
		fmt.Printf("%s is not a directory\n", config.InputDirectoryPath)
		os.Exit(1)
	}

	if !ValidateOutputFile(config.OutputJSONFilePath) {
		fmt.Printf("%s is a file\n", config.OutputJSONFilePath)
		os.Exit(1)
	}

	flChan := make(chan FileLine)

	fps := GetFilesOfExt(config.InputDirectoryPath, ".json")

	go ReadLines(fps, flChan)

	spinner := progressbar.Default(-1, "Iterating through files...")

	for fl := range flChan {
		_ = fl
		spinner.Add(1)
	}

}
