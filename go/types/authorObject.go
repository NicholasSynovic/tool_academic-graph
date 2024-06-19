package types

import "time"

type AuthorObject struct {
	// IDs
	ID    int    `json:"id"`
	OAID  string `json:"oaid"`
	ORCID string `json:"orcid"`

	// Dates
	UPDATED time.Time `json:"updated"`
	CREATED time.Time `json:"created"`

	// Counts
	AFFILIATION_COUNT int `json:"affiliation_count"`
	CITATION_COUNT    int `json:"citation_count"`
	WORKS_COUNT       int `json:"works_count"`

	// Author Information
	DISPLAY_NAME string `json:"display_name"`

	// Metrics
	IMPACT_FACTOR float64 `json:"imapct_factor"`
	H_INDEX       int     `json:"h_index"`
	I10_INDEX     int     `json:"i10_index"`

	// Meta
	FILEPATH string `json:"filepath"`
}
