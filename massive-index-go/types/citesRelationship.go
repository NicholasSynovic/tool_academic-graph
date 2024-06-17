package types

type CitesRelationship struct {
	ID        int    `json:"id"`
	Work_OAID string `json:"work_oaid"`
	Ref_OAID  string `json:"ref_oaid"`
}
