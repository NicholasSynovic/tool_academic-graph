package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Work struct {
	OA_ID          string
	DOI            string
	Title          string
	Is_Paratext    bool
	Is_Retracted   bool
	Date_Published time.Time
	OA_Type        string
	CF_Type        string
}

func createJSONObjs(jsonStrings []string, channel chan map[string]any) {
	var jsonBytes []byte

	for i := 0; i < len(jsonStrings); i++ {
		var jsonObj map[string]any

		jsonBytes = []byte(jsonStrings[i])
		err := json.Unmarshal(jsonBytes, &jsonObj)

		if err != nil {
			fmt.Println("JSON decode error", i)
			os.Exit(1)
		}

		channel <- jsonObj

	}

	close(channel)
}

func jsonToWorkObjs(jsonObjs []map[string]any, channel chan Work) {
	var jsonObj map[string]any

	caser := cases.Title(language.AmericanEnglish)

	for i := 0; i < len(jsonObjs); i++ {
		jsonObj = jsonObjs[i]

		id := strings.Replace(jsonObj["id"].(string), "https://openalex.org/", "", -1)

		doi, ok := jsonObj["doi"].(string)
		if !ok {
			continue
		}
		doi = strings.Replace(doi, "https://doi.org/", "", -1)

		title, ok := jsonObj["title"].(string)
		if !ok {
			continue
		}
		title = caser.String(title)
		title = strings.Replace(title, "\"", "", -1)
		title = strings.Replace(title, `"`, `\"`, -1)

		publishedDateString, _ := jsonObj["publication_date"].(string)
		publishedDate, _ := time.Parse("2006-01-02", publishedDateString)

		workObj := Work{
			OA_ID:          id,
			DOI:            doi,
			Title:          title,
			Is_Paratext:    jsonObj["is_paratext"].(bool),
			Is_Retracted:   jsonObj["is_retracted"].(bool),
			Date_Published: publishedDate,
			OA_Type:        jsonObj["type"].(string),
			CF_Type:        jsonObj["type_crossref"].(string),
		}

		channel <- workObj
	}
	close(channel)
}
