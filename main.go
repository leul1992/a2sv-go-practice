package main

import (
	"fmt"
	"go-basics/sum"
	"go-basics/task2"
)

func main() {
	// task 1: Sum of integers
    fmt.Println("Task 1: Sum of integers")
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println(sum.Sum(nums))

	// task 2: Word frequency
    fmt.Println("\nTask 2: Word Frequency")
	text := "hello world hello"
    res := task2.WordFrequency(text)
	for word, freq := range res {
        fmt.Printf("%s: %d\n", word, freq)
    }

    // task 2: Palindrome check
    fmt.Println("\nTask 2: Palindrome Check")
	textCheck := "law wal"
	fmt.Println(task2.IsPalindrome(textCheck))
}