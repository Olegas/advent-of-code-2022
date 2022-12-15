package main

import (
	"fmt"
	"strings"

	"github.com/Olegas/goaocd"
)

func sample() []string {
	d := ``
	return strings.Split(d, "\n")
}

func partA() int {
	done := goaocd.Duration("Part A")
	defer done()

	return 0
}

func partB() int {
	done := goaocd.Duration("Part B")
	defer done()

	return 0
}

func main() {
	lines := sample()

	fmt.Printf("Part A: %d\n", partA())
	fmt.Printf("Part B: %d\n", partB())
}
