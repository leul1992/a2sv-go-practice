package task2

import (
	"strings"
)

func WordFrequency(text string) map[string]int {
	frequency := make(map[string]int)
	words := strings.Fields(text)
	for _, word := range words {
		frequency[word]++
	}
	return frequency
}
