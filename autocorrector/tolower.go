package autocorrector

func ToLower(s string) string {
	word := []rune(s)
	for i := 0; i < len(word); i++ {
		if word[i] >= 'A' && word[i] <= 'Z' {
			word[i] = word[i] + 32
		}
	}
	return string(word)
}
