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

		if set.Contains(doi) {
			continue
		} else {
			set.Add(doi)
		}

		outChannel <- jsonObj
	}
}

/*
Convert a map[string]any object into a Dcoument object
*/
func jsonToDocumentObj(data []map[string]any, outChannel chan Document) {
	defer close(outChannel)

	for idx := range data {
		// id := _cleanOAID(data[idx]["id"].(string))

		doi, ok := data[idx]["doi"].(string)
		if !ok {
			doi = "!error"
		}
		doi = strings.Replace(doi, "https://doi.org/", "", -1)

		outChannel <- Document{DOI: doi}
	}
}
