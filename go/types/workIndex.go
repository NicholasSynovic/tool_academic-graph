package types

import "time"

type WorkIndex struct {
	ID       int       `json:"id"`
	OAID     string    `json:"oaid"`
	UPDATED  time.Time `json:"updated"`
	FILEPATH string    `json:"filepath"`
}
