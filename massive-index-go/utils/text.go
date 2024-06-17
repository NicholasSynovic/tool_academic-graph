package utils

import (
	"strings"
)

func cleanText(text string, url string) string {
	text = strings.TrimSpace(text)
	text = strings.Replace(text, url, "", -1)

	if len(text) > 0 {
		return text
	} else {
		return "!error"
	}
}

func CleanOAID(oaid string) string {
	return cleanText(oaid, "https://openalex.org/")
}

func CleanDOI(doi string) string {
	return cleanText(doi, "https://openalex.org/")
}
