package main

import (
	"fmt"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
)

func toKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func isOuter(x, y, width, height int) bool {
	return x == 0 || y == 0 || x == width-1 || y == height-1
}

func walk(mat *[][]uint8, visMap *map[string]bool, x, y, dx, dy, prevAcuuH int) {
	height := len(*mat)
	width := len((*mat)[0])
	h := int((*mat)[y][x])
	vis := h > prevAcuuH || isOuter(x, y, width, height)
	if vis {
		(*visMap)[toKey(x, y)] = true
	}
	if x+dx < width && x+dx >= 0 && y+dy < height && y+dy >= 0 {
		walk(mat, visMap, x+dx, y+dy, dx, dy, util.Max(h, prevAcuuH))
	}
}

func walkScore(mat *[][]uint8, cx, cy, dx, dy int) int {
	height := len(*mat)
	width := len((*mat)[0])
	ptH := (*mat)[cy][cx]
	score := 0
	x, y := cx+dx, cy+dy
	for {
		if 0 <= x && x < width && 0 <= y && y < height {
			h := (*mat)[y][x]
			score += 1
			if h >= ptH {
				break
			}
		} else {
			break
		}
		x, y = x+dx, y+dy
	}
	return score
}

func scoreForPoint(mat *[][]uint8, x, y int) int {
	score := walkScore(mat, x, y, 1, 0)
	score *= walkScore(mat, x, y, 0, 1)
	score *= walkScore(mat, x, y, -1, 0)
	score *= walkScore(mat, x, y, 0, -1)
	return score
}

func partA(mat *[][]uint8, width, height int) int {
	done := goaocd.Duration("Part A")
	defer done()

	visMap := make(map[string]bool)
	for y := 0; y < height; y++ {
		var start = []int{0, width - 1}
		for _, x := range start {
			dx := 1
			if x == width-1 {
				dx = -1
			}
			walk(mat, &visMap, x, y, dx, 0, -1)
		}
	}
	for x := 0; x < width; x++ {
		var start = []int{0, height - 1}
		for _, y := range start {
			dy := 1
			if y == height-1 {
				dy = -1
			}
			walk(mat, &visMap, x, y, 0, dy, -1)
		}
	}
	return len(visMap)
}

func partB(mat *[][]uint8, width, height int) int {
	done := goaocd.Duration("Part B")
	defer done()

	maxScore := 0
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			score := scoreForPoint(mat, x, y)
			if score > maxScore {
				maxScore = score
			}
		}
	}
	return maxScore
}

func main() {
	width, height, mat := goaocd.DigitMatrix()
	fmt.Printf("Part A: %d\n", partA(mat, width, height))
	fmt.Printf("Part B: %d\n", partB(mat, width, height))
}
