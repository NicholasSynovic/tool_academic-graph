package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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
Convert a map[string]any object into a Document object
*/
func jsonToDocumentObj(data []map[string]any, outChannel chan Document) {
	defer close(outChannel)

	caser := cases.Title(language.AmericanEnglish)

	for idx := range data {
		json := data[idx]
		openAccessObj := json["open_access"].(map[string]any)

		doi := json["doi"].(string)
		doi = strings.Replace(doi, "https://doi.org/", "", -1)

		title := caser.String(json["title"].(string))
		title = strings.Replace(title, "\"", "", -1)
		title = strings.Replace(title, `"`, `\"`, -1)

		publishedDateString, _ := json["publication_date"].(string)
		publication_date, _ := time.Parse("2006-01-02", publishedDateString)

		outChannel <- Document{
			DOI:              doi,
			TITLE:            title,
			PUBLICATION_DATE: publication_date,
			OA_TYPE:          json["type"].(string),
			CR_TYPE:          json["type_crossref"].(string),
			CITED_BY_COUNT:   int(json["cited_by_count"].(float64)),
			RETRACTED:        json["is_retracted"].(bool),
			OPEN_ACCESS:      openAccessObj["is_oa"].(bool),
		}
	}
}
