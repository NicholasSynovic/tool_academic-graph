package utils

import (
	"encoding/json"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/*
Given a list of keys, ensure that the JSON object has a string value for each key.

Return false if not.
*/
func checkJSONKeyExistence(obj map[string]any, keys []string) bool {
	for idx := range keys {
		_, ok := obj[keys[idx]].(string)
		if !ok {
			return false
		}
	}

	return true
}

/*
Remove OpenAlex URI from string

Returns a formatted string
*/
func cleanOAID(oa_id string) string {
	return strings.Replace(oa_id, "https://openalex.org/", "", -1)
}

/*
Remove DOI URI from string

Returns a formatted string
*/
func cleanDOI(doi string) string {
	return strings.Replace(doi, "https://doi.org/", "", -1)
}

/*
Parse a string representation of a JSON object into a map[string]any object and
pipe into a channel

On a decode error, break loop and close outChannel
*/
func CreateJSONObjs(inChannel chan string, outChannel chan map[string]any) {
	defer close(outChannel)

	set := mapset.NewSet[string]()

	for data := range inChannel {
		var jsonObj map[string]any
		err := json.Unmarshal([]byte(data), &jsonObj)

		if err != nil {
			panic(err)
		}

		// If false, continue
		if !(checkJSONKeyExistence(jsonObj, []string{"doi", "title"})) {
			continue
		}

		doi, _ := jsonObj["doi"].(string)

		// If the doi has already been identified, continue
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
func JSONToDocumentObj(data []map[string]any, outChannel chan Document) {
	defer close(outChannel)

	caser := cases.Title(language.AmericanEnglish)

	for idx := range data {
		json := data[idx]
		openAccessObj := json["open_access"].(map[string]any)

		doi := cleanDOI(json["doi"].(string))

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

/*
Convert a map[string]any object into a CitationRelationship object
*/
func JSONToCitationRelationshipObj(data []map[string]any, outChannel chan CitationRelationship) {
	defer close(outChannel)

	for idx := range data {
		json := data[idx]

		id := cleanOAID(json["id"].(string))
		refs := json["referenced_works"].([]interface{})

		for idx := range refs {
			ref := cleanOAID(refs[idx].(string))
			outChannel <- CitationRelationship{SOURCE: id, DEST: ref}
		}

	}
}

func JSONToODPObj(data []map[string]any, outChannel chan ODP) {
	defer close(outChannel)

	for idx := range data {
		json := data[idx]
		id := cleanOAID(json["id"].(string))
		doi := cleanDOI(json["doi"].(string))

		outChannel <- ODP{OAID: id, DOI: doi}
	}
}
