package types

type AuthorshipRelationship struct {
	ID          int    `json:"id"`
	AUTHOR_OAID string `json:"author_oaid"`
	WORK_OAID   string `json:"work_oaid"`
}
