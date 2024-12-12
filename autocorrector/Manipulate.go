package autocorrector

import (
	"strconv"
	"strings"
)

func Manipulate(text string) string {
	new := ""
	vowels := []rune{'a', 'e', 'i', 'o', 'u', 'h', 'A', 'E', 'I', 'O', 'U', 'H'}
	concatenated := ""
	separated := strings.Split(text, " ") // Split the input text into words
	toRemove := []int{}                   // Initialize a slice to keep track of indices to remove

	for i := 0; i < len(separated); i++ {

		// Handle "a" or "A" followed by a word starting with a vowel
		if separated[i] == "a" || separated[i] == "A" {
			if i+1 != len(separated) {
				next := separated[i+1]
				for _, c := range vowels {
					if rune(next[0]) == c {
						separated[i] += "n" // Replace "a" with "an" or "A" with "An"
					}
				}
			}
		}

		// Handle hexadecimal to decimal conversion "(hex)" command
		if separated[i] == "(hex)" {
			if i > 0 {
				new1, _ := strconv.ParseInt(separated[i-1], 16, 64) // Parse previous word as hexadecimal
				separated[i-1] = strconv.Itoa(int(new1))            // Convert decimal to string
				toRemove = append(toRemove, i)                      // Mark "(hex)" for removal
			}
		}

		// Handle binary to decimal conversion "(bin)" command
		if separated[i] == "(bin)" {
			if i > 0 {
				new1, _ := strconv.ParseInt(separated[i-1], 2, 64) // Parse previous word as binary
				separated[i-1] = strconv.Itoa(int(new1))           // Convert decimal to string
				toRemove = append(toRemove, i)                     // Mark "(bin)" for removal
			}
		}

		// Handle capitalization "(cap)" command
		if separated[i] == "(cap)" {
			if i > 0 {
				new = Capitalize(separated[i-1]) // Capitalize the previous word
				separated[i-1] = new
				toRemove = append(toRemove, i) // Mark "(cap)" for removal
			}
		}

		// Handle uppercase transformation "(up)" command
		if separated[i] == "(up)" {
			if i > 0 {
				new = ToUpper(separated[i-1]) // Convert the previous word to uppercase
				separated[i-1] = new
				toRemove = append(toRemove, i) // Mark "(up)" for removal
			}
		}

		// Handle lowercase transformation "(low)" command
		if separated[i] == "(low)" {
			if i > 0 {
				new = ToLower(separated[i-1]) // Convert the previous word to lowercase
				separated[i-1] = new
				toRemove = append(toRemove, i) // Mark "(low)" for removal
			}
		}

		// Handle capitalization for previous 'n' words "(cap," command
		if separated[i] == "(cap," {
			if i > 0 {
				num := "" // Initialize num to collect digits
				for _, v := range separated[i+1] {
					if v == ')' {
						break // Stop adding when closing parenthesis is found
					} else {
						num += string(v) // Collect digits for number
					}
				}
				n, _ := strconv.Atoi(num) // Convert collected number to integer
				for j := 1; j <= n && i-j >= 0; j++ {
					new = Capitalize(separated[i-j]) // Apply Capitalize to each of the previous 'n' words
					separated[i-j] = new
				}
				toRemove = append(toRemove, i, i+1) // Mark both "(cap," and the number for removal
			}
		}

		// Handle uppercase transformation for previous 'n' words "(up," command
		if separated[i] == "(up," {
			if i > 0 {
				num := "" // Initialize num to collect digits
				for _, v := range separated[i+1] {
					if v == ')' {
						break // Stop adding when closing parenthesis is found
					} else {
						num += string(v) // Collect digits for number
					}
				}
				n, _ := strconv.Atoi(num) // Convert collected number to integer
				for j := 1; j <= n && i-j >= 0; j++ {
					new = ToUpper(separated[i-j]) // Apply ToUpper to each of the previous 'n' words
					separated[i-j] = new
				}
				toRemove = append(toRemove, i, i+1) // Mark both "(up," and the number for removal
			}
		}

		// Handle lowercase transformation for previous 'n' words "(low," command
		if separated[i] == "(low," {
			if i > 0 {
				num := "" // Initialize num to collect digits
				for _, v := range separated[i+1] {
					if v == ')' {
						break // Stop adding when closing parenthesis is found
					} else {
						num += string(v) // Collect digits for number
					}
				}
				n, _ := strconv.Atoi(num) // Convert collected number to integer
				for j := 1; j <= n && i-j >= 0; j++ {
					new = ToLower(separated[i-j]) // Apply ToLower to each of the previous 'n' words
					separated[i-j] = new
				}
				toRemove = append(toRemove, i, i+1) // Mark both "(low," and the number for removal
			}
		}

	}

	// Create a new slice excluding the marked indices for removal
	sepShorter := []string{}
	for i, s := range separated {
		found := false
		for _, num := range toRemove {
			if i == num {
				found = true
			}
		}
		if !found {
			sepShorter = append(sepShorter, s)
		}
	}

	separated = sepShorter

	// Join the remaining words back into a single string
	concatenated = strings.Join(sepShorter, " ")
	concatenated = strings.TrimSpace(concatenated) // Remove any leading or trailing spaces

	return concatenated
}
