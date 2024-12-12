package autocorrector

func Input(text string) string{
	text = Manipulate(text) // Perform some manipulation on the text
	text = Punctuate(text) // Add punctuation to the text
	text = FixQuotes(text) // Fix quotes in the text
	return text
}
