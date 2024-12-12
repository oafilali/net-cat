package autocorrector

import (
	"strings"
)

func FixQuotes(text string) string {
	result := ""       // Stores the final output
	inQuote := false   // Tracks whether we're inside a quote
	quoteContent := "" // Stores content inside quotes

	for _, char := range text {
		if char == '\'' { // Check if the character is a single quote
			if !inQuote { // If we are outside a quote, start a new quote
				inQuote = true
			} else { // If we are inside a quote, close the current quote
				trimmedContent := strings.TrimSpace(quoteContent) // Trim spaces inside the quote
				result += "'" + trimmedContent + "'"              // Add the formatted quoted content to the result
				quoteContent = ""                                 // Reset quoteContent for future quotes
				inQuote = false                                   // Mark that we are now outside the quote
			}
		} else if inQuote { // If we are inside a quote, accumulate characters
			quoteContent += string(char)
		} else { // If we are outside a quote, add characters directly to result
			result += string(char)
		}
	}

	// Check if we are still inside a quote (unclosed quote case)
	if inQuote {
		trimmedContent := strings.TrimSpace(quoteContent) // Trim spaces and close the unclosed quote
		result += "'" + trimmedContent + "'"              // Add the final unclosed quote content to result
	}

	return strings.TrimSpace(result) // Return the final result, trimmed of extra spaces
}
