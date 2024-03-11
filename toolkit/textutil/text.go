package textutil

import "unicode/utf8"

func SplitText(text string, maxLength int) []string {
	var results []string
	var temp []rune
	var length int
	for _, char := range []rune(text) {
		var n = utf8.RuneLen(char)
		if length+n > maxLength {
			results = append(results, string(temp))
			temp = []rune{char}
			length = n
		} else {
			temp = append(temp, char)
			length += n
		}
	}
	if len(text) > 0 {
		results = append(results, string(temp))
	}
	return results
}

func RuneLength(text string) int {
	return utf8.RuneCountInString(text)
}
