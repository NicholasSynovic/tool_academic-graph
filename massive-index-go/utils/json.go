package utils

import (
	"NicholasSynovic/types"
	"encoding/json"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

func cleanOAID(oa_id string) string {
	return strings.Replace(oa_id, "https://openalex.org/", "", -1)
}

func ConvertToWorkObjs(inputChannel chan types.File_Lines) ([]types.Work_Index, int) {
	var workIndexObjs []types.Work_Index

	spinner := progressbar.Default(-1, "Creating JSON objs...")
	defer spinner.Exit()

	idCounter := 0

	for fl := range inputChannel {
		var jsonObj map[string]any
		err := json.Unmarshal([]byte(fl.Line), &jsonObj)

		if err != nil {
			panic(err)
		}

		rawOAID, _ := jsonObj["id"].(string)
		oaid := cleanOAID(rawOAID)

		updatedDateString, _ := jsonObj["updated_date"].(string)
		updatedDate, _ := time.Parse("2006-01-02T15:04:05.000000", updatedDateString)

		workIndexObj := types.Work_Index{
			ID:       idCounter,
			OAID:     oaid,
			UPDATED:  updatedDate,
			FILEPATH: fl.Filepath,
		}

		workIndexObjs = append(workIndexObjs, workIndexObj)

		idCounter += 1
		spinner.Add(1)
	}

	return workIndexObjs, idCounter
}
