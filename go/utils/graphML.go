package utils

import "NicholasSynovic/types"

func GenerateGraphMLKeys() []types.Key {
	data := []types.Key{}

	citesLabel := types.Key{
		ID:             "label",
		FOR:            "edge",
		ATTRIBUTE_NAME: "label"}

	workLabel := types.Key{
		ID:             "label",
		FOR:            "node",
		ATTRIBUTE_NAME: "label"}

	workOAID := types.Key{
		ID:             "oaid",
		FOR:            "node",
		ATTRIBUTE_NAME: "oaid"}

	workDOI := types.Key{
		ID:             "doi",
		FOR:            "node",
		ATTRIBUTE_NAME: "doi"}

	workUpdated := types.Key{
		ID:             "updated",
		FOR:            "node",
		ATTRIBUTE_NAME: "updated"}

	data = append(data, citesLabel, workLabel, workOAID, workDOI, workUpdated)

	return data
}
