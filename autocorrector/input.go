package autocorrector

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

var badWords[]string

func Input(text string) string{
	umarshalBadWords()
	text = Manipulate(text) // Perform some manipulation on the text
	text = Punctuate(text) // Add punctuation to the text
	text = FixQuotes(text) // Fix quotes in the text
	text = isThisaBadWord(text)
	return text
}

func umarshalBadWords() {
	data, err := os.ReadFile("autocorrector/words.json")
	if err != nil {
		log.Println("error reading the json file", err)
	}

	err = json.Unmarshal(data, &badWords)
	if err != nil {
		log.Println("error unmarshaling the json file", err)
	}
}

func isThisaBadWord(text string) string {
	text = strings.ToLower(text)
	words := strings.Split(text, " ")
	for i := 0; i < len(words); i++ {
		for j := 0; j < len(badWords); j++ {
			if words[i] == badWords[j] {
				words[i] = string(words[i][0]) + strings.Repeat("*", len(words[i])-1)
			}
		}
	}
	text = strings.Join(words, " ")
	return text
}