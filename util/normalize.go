package util

import (
	"log"
	"regexp"
	"strings"
	"unicode"
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

// RepairB reverts the changes made by a normalize by fixing case and adding punctuation
func RepairB(text, oldText []byte) []byte {
	var repairedText []byte
	oldTextCounter := 0
	for _, c := range oldText {
		char := rune(text[oldTextCounter])
		if c >= 'A' && c <= 'Z' {
			repairedText = append(repairedText, byte(unicode.ToUpper(char)))
			oldTextCounter++
		} else if c >= 'a' && c <= 'z' {
			repairedText = append(repairedText, byte(unicode.ToLower(char)))
			oldTextCounter++
		} else {
			repairedText = append(repairedText, c)
		}
	}
	return repairedText
}
