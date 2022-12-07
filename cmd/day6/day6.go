package main

import (
	"fmt"

	"github.com/Olegas/goaocd"
)

func allDiff(s string) bool {
	var chars map[rune]bool = make(map[rune]bool)
	for _, c := range s {
		_, ok := chars[c]
		if ok {
			return false
		}
		chars[c] = true
	}
	return true
}

func detectSignal(s string, numChars int) int {
	for i := 0; i < len(s)-numChars; i++ {
		part := s[i : i+numChars]
		if allDiff(part) {
			return i + numChars
		}
	}
	return -1
}

func main() {
	line := goaocd.Input()
	fmt.Printf("Part A: %d\n", detectSignal(line, 4))
	fmt.Printf("Part B: %d\n", detectSignal(line, 14))
}
