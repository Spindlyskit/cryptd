package util

import (
	"log"
	"regexp"
	"strings"
)

// Normalize makes a string uppercase and removes non alphabetic characters
func Normalize(str string) string {
	str = strings.ToUpper(str)
	// TODO Move this regex to a property on the TextScorer struct
	reg, err := regexp.Compile("[^A-Z]+")
	if err != nil {
		log.Fatalf("Error compiling regex: %s", err)
	}
	str = reg.ReplaceAllString(str, "")
	return str
}

// NormalizeB makes a byte slice uppercase and removes non alphabetic characters
func NormalizeB(text []byte) []byte {
	return []byte(Normalize(string(text)))
}
