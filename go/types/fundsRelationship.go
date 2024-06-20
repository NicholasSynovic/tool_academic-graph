package types

type FunderRelationship struct {
	ID          int    `json:"id"`
	FUNDER_OAID string `json:"funder_oaid"`
	WORK_OAID   string `json:"work_oaid"`
}
