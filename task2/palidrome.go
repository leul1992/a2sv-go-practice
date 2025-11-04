package task2

import (
	"unicode"
)

func IsPalindrome(s string) bool {
	n := len(s)
	for i := 0; i < n/2; i++ {
		if unicode.IsLetter(rune(s[i])) == false || unicode.IsLetter(rune(s[n-1-i])) == false {
			continue
		}

		if unicode.ToLower(rune(s[i])) != unicode.ToLower(rune(s[n-1-i])) {
			return false
		}
	}
	return true
}
