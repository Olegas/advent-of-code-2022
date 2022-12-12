package day12

import (
	"fmt"
	"os"
	"time"

	"github.com/Olegas/goaocd"
	"github.com/fatih/color"
)

var visitedPlace = color.New(color.FgWhite, color.Bold).SprintFunc()
var normalPlace = color.New(color.FgCyan).SprintFunc()

func Visualize(mat *[][]string, visited *map[goaocd.Pos]bool, edgePositions *map[string]goaocd.Pos) {
	if os.Getenv("VIS") == "" {
		return
	}
	// util.ClearTerm()
	lineStr := "\033[H\033[2J"
	for y, line := range *mat {
		for x, s := range line {
			pos := goaocd.Pos{X: x, Y: y}
			var char string
			if pos.Eq((*edgePositions)["S"]) {
				char = "S"
			} else if pos.Eq((*edgePositions)["E"]) {
				char = "E"
			} else {
				char = s
			}
			if (*visited)[pos] {
				lineStr += visitedPlace(char)
			} else {
				lineStr += normalPlace(char)
			}
		}
		lineStr += "\n"
	}
	fmt.Print(lineStr)
	time.Sleep(10 * time.Millisecond)
}
