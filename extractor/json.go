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
Parse a string representation of a JSON object into a map[string]any object and
pipe into a channel

On a decode error, break loop and close outChannel
*/
func createJSONObjs(data []string, outChannel chan map[string]any) {
	for i := range data {
		var jsonObj map[string]any
		err := json.Unmarshal([]byte(data[i]), &jsonObj)

		if err != nil {
			fmt.Println(err, data[i])
			break
		}

		outChannel <- jsonObj
	}

	close(outChannel)
}

/*
Convert a map[string]any object into a Work object

On conversion error, continue to the next iteration of the for loop
*/
func jsonToWorkObj(obj map[string]any, outChannel chan Work) {
	caser := cases.Title(language.AmericanEnglish)

	id := _cleanOAID(obj["id"].(string))

	doi, ok := obj["doi"].(string)
	if !ok {
		return
	}
	doi = strings.Replace(doi, "https://doi.org/", "", -1)

	title, ok := obj["title"].(string)
	if !ok {
		return
	}
	title = caser.String(title)
	title = strings.Replace(title, "\"", "", -1)
	title = strings.Replace(title, `"`, `\"`, -1)

	publishedDateString, _ := obj["publication_date"].(string)
	publishedDate, _ := time.Parse("2006-01-02", publishedDateString)

	workObj := Work{
		OA_ID:          id,
		DOI:            doi,
		Title:          title,
		Is_Paratext:    obj["is_paratext"].(bool),
		Is_Retracted:   obj["is_retracted"].(bool),
		Date_Published: publishedDate,
		OA_Type:        obj["type"].(string),
		CF_Type:        obj["type_crossref"].(string),
	}

	outChannel <- workObj
}

/*
Convert a map[string]any object into a Work object

On conversion error, continue to the next iteration of the for loop
*/
func jsonToCitationObj(obj map[string]any, outChannel chan Citation) {
	id := _cleanOAID(obj["id"].(string))

	refs := obj["referenced_works"].([]interface{})

	for idx := range refs {
		refID := _cleanOAID(refs[idx].(string))

		citationObj := Citation{
			Work_OA_ID: id,
			Ref_OA_ID:  refID,
		}

		outChannel <- citationObj

	}
}
