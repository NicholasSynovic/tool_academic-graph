package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

func writeGraphML(fp string, graphML GraphML) {
	fmt.Println("Writing data to", filepath.Base(fp))
	file, _ := os.Create(fp)
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(graphML); err != nil {
		fmt.Println("Error encoding XML:", err)
		return
	}
}
