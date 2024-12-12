package autocorrector

func ToUpper(s string) string {
	word := []rune(s)
	for i := 0; i < len(word); i++ {
		if word[i] >= 97 && word[i] <= 122 {
			word[i] = word[i] - 32
		}
	}
	return string(word)
}
