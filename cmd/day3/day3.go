package main

import (
	"fmt"

	"github.com/Olegas/goaocd"
	mapset "github.com/deckarep/golang-set"
)

func isSmall(a byte) bool {
	return 'a' <= a && a <= 'z'
}

func priority(sym byte) int {
	if isSmall(sym) {
		return int(sym - 'a' + 1)
	} else {
		return int(sym - 'A' + 27)
	}
}

func main() {
	text := goaocd.Lines(3)
	var accu = 0
	for i := 0; i <= len(text)-3; i += 3 {
		lines := text[i : i+3]
		var groupSet mapset.Set
		for idx, line := range lines {
			set := mapset.NewSet()
			for _, s := range line {
				set.Add(s)
			}
			if idx == 0 {
				groupSet = set
			} else {
				groupSet = groupSet.Intersect(set)
			}
		}
		if groupSet.Cardinality() != 1 {
			panic("To many elements")
		}
		sym := byte(groupSet.Pop().(int32))
		accu += priority(sym)
	}

	fmt.Printf("%d\n", accu)
}
