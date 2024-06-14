package main

import (
	"NicholasSynovic/utils"
	"fmt"
)

func main() {
	config := utils.ParseCommandLine()

	fmt.Println(config.OAWorkJSONDirectoryPath)
	fmt.Println(config.OutputDBPath)

}
