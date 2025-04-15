package utility

import (
	"regexp"
	"strings"
	"unicode"
)

// SynthesizeString synthesizes a query by:
// 1. Trimming leading and trailing whitespaces
// 2. Removing duplicates (words)
// 3. Replacing separation symbols with a single space
// 4. Uppercasing the first letter of each word
// 5. Joining words with a comma
func SynthesizeString(input string) string {
	// Trim leading and trailing whitespaces
	trimmed := strings.TrimSpace(input)

	// Replace any separation symbols (non-letter/non-digit characters) with a single space
	re := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	normalized := re.ReplaceAllString(trimmed, " ")

	// Remove duplicate words
	words := strings.Fields(normalized)
	uniqueWords := make([]string, 0, len(words))
	seen := make(map[string]bool)

	for _, word := range words {
		if !seen[word] {
			seen[word] = true
			// Capitalize the first letter
			uniqueWords = append(uniqueWords, capitalize(word))
		}
	}

	// Join unique words with a comma
	result := strings.Join(uniqueWords, ", ")

	return result
}

// capitalize returns the string with the first letter in uppercase
func capitalize(word string) string {
	if len(word) == 0 {
		return ""
	}
	runes := []rune(word)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func IsCyrillic(text string) bool {
	var counter = 0.0
	var total = 0.0

	if len(text) > 100 {
		text = text[:50]
	}

	for _, v := range text {
		if unicode.Is(unicode.Cyrillic, v) {
			counter++
		}
		total++
	}

	return counter/total > 0.8
}

func IsEnglish(text string) bool {
	var counter = 0.0
	var total = 0.0

	if len(text) > 100 {
		text = text[:50]
	}

	for _, v := range text {
		if unicode.Is(unicode.Latin, v) {
			counter++
		}
		total++
	}

	return counter/total > 0.8
}
