package autocorrector

func isAlphanumeric(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

func Capitalize(s string) string {
	words := []byte(s)
	wordStart := true

	for i := 0; i < len(words); i++ {
		if isAlphanumeric(words[i]) {
			if wordStart {
				if words[i] >= 'a' && words[i] <= 'z' {
					words[i] -= 32
				}
				wordStart = false
			} else {
				if words[i] >= 'A' && words[i] <= 'Z' {
					words[i] += 32
				}
			}
		} else {
			wordStart = true
		}
	}

	return string(words)
}
