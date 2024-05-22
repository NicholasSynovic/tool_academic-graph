package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/*
Type to represent an academic work from OpenAlex

Does not include all of the descriptors provided by OpenAlex, only the most
relevant ones for this application
*/
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

/*
Type to represent the citations of an academic work
*/
type Citation struct {
	Work_OA_ID string
	Ref_OA_ID  string
}

/*
A type to store an array of works to by marshelled into a JSON formatted string
*/
type CitationOutput []Citation

/*
A type to store an array of works to by marshelled into a JSON formatted string
*/
type WorkOutput []Work

func _cleanOAID(oa_id string) string {
	return strings.Replace(oa_id, "https://openalex.org/", "", -1)
}

/*
Parse a string representation of a JSON object into a map[string]any object

On a decode error, exit the application with code 1
*/
func createJSONObjs(jsonStrings []string, channel chan map[string]any) {
	var jsonBytes []byte

	size := int64(len(jsonStrings))
	bar := progressbar.Default(size, "Converting JSON strings to JSON objs...")

	for i := int64(0); i < size; i++ {
		var jsonObj map[string]any

		jsonBytes = []byte(jsonStrings[i])
		err := json.Unmarshal(jsonBytes, &jsonObj)

		if err != nil {
			fmt.Println("JSON decode error", i)
			os.Exit(1)
		}

		bar.Add(1)
		channel <- jsonObj

	}
	bar.Finish()
	bar.Exit()
	close(channel)
}

/*
Convert a map[string]any object into a Work object

On conversion error, continue to the next iteration of the for loop
*/
func jsonToWorkObjs(jsonObjs []map[string]any, channel chan Work) {
	var jsonObj map[string]any

	caser := cases.Title(language.AmericanEnglish)

	size := int64(len(jsonObjs))
	bar := progressbar.Default(size, "Converting JSON obs to Work objs...")

	for i := 0; i < len(jsonObjs); i++ {
		jsonObj = jsonObjs[i]

		id := _cleanOAID(jsonObj["id"].(string))

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

		bar.Add(1)
		channel <- workObj
	}
	bar.Finish()
	bar.Exit()
	close(channel)
}

/*
Convert a map[string]any object into a Work object

On conversion error, continue to the next iteration of the for loop
*/
func jsonToCitationObjs(jsonObjs []map[string]any, channel chan Citation) {
	var jsonObj map[string]any

	size := int64(len(jsonObjs))
	bar := progressbar.Default(size, "Converting JSON obs to Citation objs...")

	for i := 0; i < len(jsonObjs); i++ {
		jsonObj = jsonObjs[i]

		id := _cleanOAID(jsonObj["id"].(string))

		refs := jsonObj["referenced_works"].([]interface{})
		bar.ChangeMax(bar.GetMax() + len(refs))

		for idx := range refs {
			refID := _cleanOAID(refs[idx].(string))

			citationObj := Citation{
				Work_OA_ID: id,
				Ref_OA_ID:  refID,
			}

			bar.Add(1)
			channel <- citationObj

		}
	}
	bar.Finish()
	bar.Exit()
	close(channel)
}
