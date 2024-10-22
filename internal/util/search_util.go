package util

import "unicode"

func RemoveSymbols(input string) string {
	var result []rune

	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			result = append(result, r)
		}
	}

	return string(result)
}

func MaxWorker(longData int) int {
	return (longData / 7) + 1
}
