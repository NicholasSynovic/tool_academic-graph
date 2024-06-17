package types

import "time"

type WorkObject struct {
	// IDs
	ID   int    `json:"id"`
	DOI  string `json:"doi"`
	OAID string `json:"oaid"`

	// Dates
	UPDATED   time.Time `json:"updated"`
	CREATED   time.Time `json:"created"`
	PUBLISHED time.Time `json:"published"`

	// Authorships
	AUTHORSHIP_COUNT       int `json:"authorship_count"`
	INSTITUTION_COUNT      int `json:"institution_count"`
	DISTINCT_COUNTRY_COUNT int `json:"distinct_country_count"`

	// Categories
	CONCEPT_COUNT int    `json:"concept_count"`
	KEYWORD_COUNT int    `json:"keyword_count"`
	GRANT_COUNT   int    `json:"grant_count"`
	TOPIC_COUNT   int    `json:"topic_count"`
	IS_PARATEXT   bool   `json:"is_paratext"`
	IS_RETRACTED  bool   `json:"is_retracted"`
	LANGUAGE      string `json:"language"`
	LICENSE       string `json:"license"`

	// Publication Metrics
	CITED_BY_COUNT             int `json:"cited_by_count"`
	PUBLICATION_LOCATION_COUNT int `json:"publication_location_count"`

	// Document Metrics
	SUSTAINABLE_DEVELOPMENT_GOAL_COUNT int    `json:"sustainable_development_goal_count"`
	TITLE                              string `json:"title"`
	OA_TYPE                            string `json:"oa_type"`
	CR_TYPE                            string `json:"cr_type"`

	// Meta
	FILEPATH string `json:"filepath"`
}
