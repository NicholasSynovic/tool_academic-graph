package types

type KeywordRelationship struct {
	ID           int    `json:"id"`
	KEYWORD_OAID string `json:"keyword_oaid"`
	WORK_OAID    string `json:"work_oaid"`
}
