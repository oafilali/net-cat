package basic

import (
	"fmt"
	"os"

	"github.com/atouba/piscine"
)

// printBasic generates ascii art for the function Basic()
func printBasic(inLineStr, banner string) string {
	asciiArtChars, err := os.ReadFile("./banners/" + banner + ".txt")
	if err != nil {
		fmt.Println("Error reading banner file")
		return ""
	}

	out := ""

	linesArt := piscine.Split(clearCarReturns(string(asciiArtChars)), "\n")
	for indexLine := range 8 {
		for _, char := range inLineStr {
			out += linesArt[(int(char)-32)*8+indexLine]
		}
		out += fmt.Sprintln()
	}

	return out
}

// Basic returns an ascii art text string from a string str
func Basic(str, banner string) string {
	var newLineI int
	i := 0
	out := ""
	prevIsNL := true

	for i < len(str) {
		newLineI = index(str[i:], "\\n")
		if newLineI == 0 {
			if prevIsNL || i == len(str)-2 {
				out += fmt.Sprintln()
			}
			i += 2
			prevIsNL = true
		} else {
			out += printBasic(str[i:i+newLineI], banner)
			i += newLineI
			prevIsNL = false
		}
	}

	return out
}

// index returns index of subStr, if not found
// returns the length of str
func index(str, subStr string) int {
	iSubStr := piscine.Index(str, subStr)
	if iSubStr == -1 {
		return len(str)
	}
	return iSubStr
}

// clearCarReturns returns the input string without carriage returns
func clearCarReturns(s string) (out string) {
	for _, r := range s {
		if r != 13 {
			out += string(r)
		}
	}
	return
}
