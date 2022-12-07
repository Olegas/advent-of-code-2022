package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Olegas/goaocd"
)

type Range struct {
	left, right int
}

func parseRange(s string) Range {
	numbers := strings.Split(s, "-")
	a, _ := strconv.Atoi(numbers[0])
	b, _ := strconv.Atoi(numbers[1])
	return Range{a, b}
}

func isFullInside(a, b Range) bool {
	if a.left <= b.left && b.right <= a.right {
		return true
	}
	if b.left <= a.left && a.right <= b.right {
		return true
	}
	return false
}

func isOverlap(a, b Range) bool {
	if a.left <= b.left && b.left <= a.right {
		return true
	}
	if a.left <= b.right && b.right <= a.right {
		return true
	}
	if b.left <= a.left && a.left <= b.right {
		return true
	}
	if b.left <= a.right && a.right <= b.right {
		return true
	}
	return false
}

func main() {
	var accuA = 0
	var accuB = 0
	for _, line := range goaocd.Lines() {
		ranges := strings.Split(line, ",")
		r1 := parseRange(ranges[0])
		r2 := parseRange(ranges[1])
		if isFullInside(r1, r2) {
			accuA += 1
		}
		if isOverlap(r1, r2) {
			accuB += 1
		}
	}
	fmt.Printf("Part A: %d\n", accuA)
	fmt.Printf("Part B: %d\n", accuB)

	// aocd.Submit(1, accuA)
	// aocd.Submit(2, accuB)
}
