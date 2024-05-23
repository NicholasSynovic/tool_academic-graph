package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

	for data := range inChannel {
		var jsonObj map[string]any
		err := json.Unmarshal([]byte(data), &jsonObj)

		if err != nil {
			fmt.Println(err, data)
			break
		}

		outChannel <- jsonObj
	}
}

/*
Convert a map[string]any object into a Work object
*/
func jsonToWorkObj(inChannel chan map[string]any, outChannel chan Work) {
	defer close(outChannel)

	for data := range inChannel {
		caser := cases.Title(language.AmericanEnglish)

		id := _cleanOAID(data["id"].(string))

		doi, ok := data["doi"].(string)
		if !ok {
			doi = "!error"
		}
		doi = strings.Replace(doi, "https://doi.org/", "", -1)

		title, ok := data["title"].(string)
		if !ok {
			title = "!error"

		}
		title = caser.String(title)
		title = strings.Replace(title, "\"", "", -1)
		title = strings.Replace(title, `"`, `\"`, -1)

		publishedDateString, _ := data["publication_date"].(string)
		publishedDate, _ := time.Parse("2006-01-02", publishedDateString)

		workObj := Work{
			OA_ID:          id,
			DOI:            doi,
			Title:          title,
			Is_Paratext:    data["is_paratext"].(bool),
			Is_Retracted:   data["is_retracted"].(bool),
			Date_Published: publishedDate,
			OA_Type:        data["type"].(string),
			CF_Type:        data["type_crossref"].(string),
		}
		outChannel <- workObj
	}
}

/*
Convert a map[string]any object into a Citation object
*/
func jsonToCitationObj(inChannel chan map[string]any, outChannel chan Citation) {
	defer close(outChannel)

	for data := range inChannel {
		id := _cleanOAID(data["id"].(string))

		refs := data["referenced_works"].([]interface{})

		for idx := range refs {
			refID := _cleanOAID(refs[idx].(string))

			citationObj := Citation{
				Work_OA_ID: id,
				Ref_OA_ID:  refID,
			}

			outChannel <- citationObj
		}
	}
}
