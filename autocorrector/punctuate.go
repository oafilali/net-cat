package autocorrector

import (
	"strings"
)

func Punctuate(input string) string {
	// List of punctuation marks we want to format
	punctuationMarks := []string{".", ",", "!", "?", ":", ";"}

	// Remove space before punctuation
	for _, mark := range punctuationMarks {
		input = strings.ReplaceAll(input, " "+mark, mark)
	}

	// Add space after punctuation if needed
	for i := 0; i < len(input)-1; i++ {
		// Check if the current character is a punctuation mark
		if strings.ContainsAny(string(input[i]), ".,!?:;") {
			// If next character is not a punctuation or a space, add a space
			if !strings.ContainsAny(string(input[i+1]), ".,!?:; ") {
				input = input[:i+1] + " " + input[i+1:]
			}
		}
	}

	return input
}
