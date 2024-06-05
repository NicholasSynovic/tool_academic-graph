package main

import (
	"encoding/json"
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

/*
Remove OpenAlex URI from string

Returns a formatted string
*/
func _cleanOAID(oa_id string) string {
	return strings.Replace(oa_id, "https://openalex.org/", "", -1)
}

/*
Parse a string representation of a JSON object into a map[string]any object and
pipe into a channel

On a decode error, break loop and close outChannel
*/
func createJSONObjs(inChannel chan string, outChannel chan map[string]any) {
	defer close(outChannel)

	set := mapset.NewSet[string]()

	for data := range inChannel {
		var jsonObj map[string]any
		err := json.Unmarshal([]byte(data), &jsonObj)

		if err != nil {
			fmt.Println(err, data)
			break
		}

		doi, ok := jsonObj["doi"].(string)
		if !ok {
			continue
		}

		_, ok = jsonObj["title"].(string)
		if !ok {
			continue
		}

		if set.Contains(doi) {
			continue
		} else {
			set.Add(doi)
		}

		outChannel <- jsonObj
	}
}

/*
Convert a map[string]any object into a Citation object
*/
func jsonToCitationtObj(data []map[string]any, outChannel chan Citation) {
	defer close(outChannel)

	for idx := range data {
		json := data[idx]

		sourceID := _cleanOAID(json["id"].(string))
		refIDs := json["referenced_works"].([]any)

		for refIDX := range refIDs {
			outChannel <- Citation{
				SOURCE: sourceID,
				DEST:   _cleanOAID(refIDs[refIDX].(string)),
			}
		}

	}
}
