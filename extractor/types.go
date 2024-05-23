package main

import "time"

/*
Type to represent application configuration

Modified by user command line input
*/
type AppConfig struct {
	inputPath, worksOutputPath, citesOutputPath string
	processes                                   int
}

/*
Type to represent the citations of an academic work

Does not include all of the descriptors provided by OpenAlex, only the most
relevant ones for this application
*/
type Citation struct {
	Work_OA_ID string
	Ref_OA_ID  string
}

/*
Type to represent an academic work from OpenAlex

Does not include all of the descriptors provided by OpenAlex, only the most
relevant ones for this application
*/
type Work struct {
	OA_ID          string
	DOI            string
	Title          string
	Is_Paratext    bool
	Is_Retracted   bool
	Date_Published time.Time
	OA_Type        string
	CF_Type        string
}
