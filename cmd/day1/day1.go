package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Olegas/goaocd"
)

func main() {
	var accu int = 0
	var elves = []int{}
	for _, line := range goaocd.Lines(1) {
		if len(line) == 0 {
			elves = append(elves, accu)
			accu = 0
			continue
		}
		value, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		accu += value
	}
	sort.Ints(elves)
	top := elves[len(elves)-3:]
	accu = 0
	for _, v := range top {
		accu += v
	}
	fmt.Printf("%d", accu)
}
