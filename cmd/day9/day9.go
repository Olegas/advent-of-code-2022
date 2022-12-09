package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Olegas/advent-of-code-2022/internal/day9"
	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
)

var visualize bool
var pause bool

func display(r *day9.Rope, visited *map[string]bool, showVisited ...bool) {
	if !visualize {
		return
	}
	fmt.Print("\033[H\033[2J")
	xCoords := []int{}
	yCoords := []int{}
	for c := range *visited {
		var x, y int
		fmt.Sscanf(c, "%d,%d", &x, &y)
		xCoords = append(xCoords, x)
		yCoords = append(yCoords, y)
	}
	minX, maxX := util.MinMax(xCoords)
	minY, maxY := util.MinMax(yCoords)

	for y := minY - 5; y <= maxY+5; y++ {
		for x := minX - 5; x <= maxX+5; x++ {
			k := fmt.Sprintf("%d,%d", x, y)
			ok := (*visited)[k]
			s := "."

			for k := r.Head; k != nil; k = k.Next {
				pos := k.Pos
				if pos.X == x && pos.Y == y {
					s = k.Id
					break
				}
			}
			if ok && s == "." {
				if len(showVisited) == 1 && showVisited[0] {
					s = "#"
				}
			}
			fmt.Print(s)
		}
		fmt.Print("\n")
	}
	if pause {
		fmt.Scanln()
	} else {
		time.Sleep(50 * time.Millisecond)
	}
}

func simulate(lines *[]string, ropeLen int) int {
	visited := map[string]bool{}
	rope := day9.NewRope(ropeLen)

	// At start count current tail position as visited
	visited[rope.Head.Pos.ToString()] = true

	for _, line := range *lines {
		var direction string
		var count int
		n, err := fmt.Sscanf(line, "%s %d", &direction, &count)
		if err != nil || n != 2 {
			panic(fmt.Sprintf("Unexpected line %s", line))
		}
		for i := 0; i < count; i++ {
			rope.MoveHead(direction)
			rope.Adjust()
			display(rope, &visited, true)

			visited[rope.Tail().Pos.ToString()] = true
		}
	}

	return len(visited)
}

func partA(lines *[]string) int {
	done := goaocd.Duration("Part A")
	defer done()

	return simulate(lines, 2)
}

func partB(lines *[]string) int {
	done := goaocd.Duration("Part B")
	defer done()

	return simulate(lines, 10)
}

func main() {
	visualize = os.Getenv("VIS") != ""
	pause = os.Getenv("PAUSE") != ""

	lines := goaocd.Lines()
	a := partA(&lines)
	fmt.Printf("Part A: %d\n", a)

	b := partB(&lines)
	fmt.Printf("Part B: %d\n", b)
}
