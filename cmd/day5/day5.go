package main

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
)

func main() {
	lines := goaocd.Lines(5)
	/*
	   	sample := `
	       [D]
	   [N] [C]
	   [Z] [M] [P]
	    1   2   3

	   move 1 from 2 to 1
	   move 3 from 1 to 3
	   move 2 from 2 to 1
	   move 1 from 1 to 2`
	   	sample = strings.TrimLeft(sample, "\n")
	   	lines := strings.Split(sample, "\n")
	*/

	// Set to true to Part 2
	var moveMultipleCratesAtOnce = true
	var stacks []*list.List
	var indexLine string
	var indexPos int
	for idx, line := range lines {
		if line[:2] == " 1" {
			indexLine = line
			indexPos = idx
			break
		}
	}

	for idx, s := range indexLine {
		if util.IsDigit(byte(s)) {
			curStack := list.New()
			stacks = append(stacks, curStack)
			for i := indexPos - 1; i >= 0; i-- {
				crateName := lines[i][idx]
				if crateName == ' ' {
					break
				}
				curStack.PushBack(crateName)
			}
		}
	}

	for idx, commandLine := range lines[indexPos+2:] {
		var count, from, to int8
		n, err := fmt.Sscanf(commandLine, "move %d from %d to %d", &count, &from, &to)
		if err != nil {
			panic(err)
		}
		if n != 3 {
			panic("Not enought data")
		}
		fromStack := stacks[from-1]
		toStack := stacks[to-1]
		if moveMultipleCratesAtOnce {
			// Part 2
			elt := fromStack.Back()
			for i := count; i > 1; i-- {
				elt = elt.Prev()
			}
			for ok := true; ok; ok = elt != nil {
				toStack.PushBack(elt.Value)
				t := elt
				elt = elt.Next()
				fromStack.Remove(t)
			}
		} else {
			// Part 1
			for i := count; i > 0; i-- {
				elt := fromStack.Back()
				if elt == nil {
					panic(fmt.Sprintf("Stack %d is empty when doing command %s at line %d", from, commandLine, idx))
				}
				toStack.PushBack(fromStack.Remove(elt))
			}
		}

	}

	var accu []string
	for _, stack := range stacks {
		elt := stack.Back().Value
		accu = append(accu, string(elt.(byte)))
	}
	fmt.Print(strings.Join(accu, ""))
}
