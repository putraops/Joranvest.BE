package helper

import (
	"log"
	"regexp"
	"strings"
)

func StringToPathUrl(text string) string {
	var result string
	regNonAplhaNumeric, err := regexp.Compile("[^a-zA-Z0-9- ]")

	if err != nil {
		log.Fatal(err)
	}

	space := regexp.MustCompile(`\s+`)         //to remove double space
	result = space.ReplaceAllString(text, " ") //-- regex double space implementation

	result = regNonAplhaNumeric.ReplaceAllString(result, "") //-- remove non-alphanumeric
	result = strings.ReplaceAll(result, " ", "-")            //-- replace space into -
	result = strings.ToLower(result)
	return result
}
