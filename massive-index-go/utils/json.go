package utils

import (
	"NicholasSynovic/types"
	"encoding/json"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func cleanOAID(oa_id string) string {
	return strings.Replace(oa_id, "https://openalex.org/", "", -1)
}

func ConvertToWorkObjs(inputChannel chan string) ([]types.Work_Index, int) {
	var workIndexObjs []types.Work_Index

	spinner := progressbar.Default(-1, "Creating JSON objs...")
	defer spinner.Exit()

	idCounter := 0

	for document := range inputChannel {
		var jsonObj map[string]any
		err := json.Unmarshal([]byte(document), &jsonObj)

		if err != nil {
			panic(err)
		}

		rawOAID := jsonObj["id"].(string)
		oaid := cleanOAID(rawOAID)

		workIndexObj := types.Work_Index{
			ID:      idCounter,
			OAID:    oaid,
			UPDATED: "test",
		}

		workIndexObjs = append(workIndexObjs, workIndexObj)

		idCounter += 1
		spinner.Add(1)
	}

	return workIndexObjs, idCounter
}