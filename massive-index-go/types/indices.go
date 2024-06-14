package types

import "time"

type Work_Index struct {
	ID      int       `json:"id"`
	OAID    string    `json:"oaid"`
	UPDATED time.Time `json:"updated"`
}
