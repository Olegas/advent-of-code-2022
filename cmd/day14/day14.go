package main

import (
	"fmt"
	"strings"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
)

func NewPos(s string) goaocd.Pos {
	coords := strings.Split(s, ",")
	x := util.Atoi(coords[0])
	y := util.Atoi(coords[1])
	return goaocd.Pos{X: x, Y: y}
}

type Dimensions struct {
	MinX, MaxX, MinY, MaxY int
}

func caveDimensions(cave map[goaocd.Pos]string) Dimensions {
	xCoords := []int{}
	yCoords := []int{}
	for c := range cave {
		xCoords = append(xCoords, c.X)
		yCoords = append(yCoords, c.Y)
	}
	minX, maxX := util.MinMax(xCoords)
	minY, maxY := util.MinMax(yCoords)

	return Dimensions{MinX: minX, MaxX: maxX, MinY: minY, MaxY: maxY}
}

func display(visited map[goaocd.Pos]string) {
	dim := caveDimensions(visited)
	padding := 2

	for y := dim.MinY - padding; y <= dim.MaxY+padding; y++ {
		for x := dim.MinX - padding; x <= dim.MaxX+padding; x++ {
			k := goaocd.Pos{X: x, Y: y}
			s, ok := visited[k]

			if !ok {
				s = "."
			}
			fmt.Print(s)
		}
		fmt.Print("\n")
	}
}

func sample() []string {
	d := `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`
	return strings.Split(d, "\n")
}

type Sand struct {
	goaocd.Pos
	Infinite bool
	Steady   bool
}

func (s *Sand) Fall(cave map[goaocd.Pos]string, dim Dimensions, floorPos *goaocd.Pos) {
	check := []goaocd.Pos{{X: 0, Y: 1}, {X: -1, Y: 1}, {X: 1, Y: 1}}
	for idx, mut := range check {
		next := s.Mut(mut)
		_, found := cave[next]
		if !found && floorPos != nil && next.Y == floorPos.Y {
			found = true
		}
		if !found {
			s.Pos = next
			break
		} else {
			if idx < len(check)-1 {
				continue
			}
			cave[s.Pos] = "o"
			s.Steady = true
			return
		}
	}

	if s.Y > dim.MaxY {
		s.Infinite = true
		return
	}
}

func emitNewSand() Sand {
	return Sand{Pos: goaocd.Pos{X: 500, Y: 0}}
}

func partA(cave map[goaocd.Pos]string, dim Dimensions) int {
	done := goaocd.Duration("Part A")
	defer done()

	steady := 0
	for {
		sand := emitNewSand()
		for {
			sand.Fall(cave, dim, nil)
			if sand.Steady {
				steady++
				break
			} else if sand.Infinite {
				return steady
			}
		}
	}
}

func partB(cave map[goaocd.Pos]string, dim Dimensions) int {
	done := goaocd.Duration("Part B")
	defer done()

	start := goaocd.Pos{X: 500, Y: 0}
	floor := goaocd.Pos{X: 0, Y: dim.MaxY + 2}
	steady := 0
	for {
		sand := emitNewSand()
		for {
			sand.Fall(cave, dim, &floor)
			if sand.Steady {
				steady++
				if sand.Pos.Eq(start) {
					return steady
				}
				break
			}
		}
	}
}

func loadCave(lines []string) (map[goaocd.Pos]string, Dimensions) {
	var cave = make(map[goaocd.Pos]string)
	cave[goaocd.Pos{X: 500, Y: 0}] = "+"
	for _, line := range lines {
		steps := strings.Split(line, " -> ")
		for idx, step := range steps {
			if idx > 0 {
				prev := NewPos(steps[idx-1])
				curr := NewPos(step)

				if prev.X == curr.X {
					yMut := 1
					if prev.Y > curr.Y {
						yMut = -1
					}
					mut := goaocd.Pos{Y: yMut, X: 0}
					for {
						cave[prev] = "#"
						prev = prev.Mut(mut)
						if prev.Y == curr.Y {
							cave[prev] = "#"
							break
						}
					}
				} else {
					xMut := 1
					if prev.X > curr.X {
						xMut = -1
					}
					mut := goaocd.Pos{Y: 0, X: xMut}
					for {
						cave[prev] = "#"
						prev = prev.Mut(mut)
						if prev.X == curr.X {
							cave[prev] = "#"
							break
						}
					}
				}
			}
		}
	}
	return cave, caveDimensions(cave)
}

func main() {
	lines := goaocd.Lines()
	// lines := sample()
	cave, dim := loadCave(lines)
	fmt.Printf("Part A: %d\n", partA(cave, dim))

	cave, dim = loadCave(lines)
	fmt.Printf("Part B: %d\n", partB(cave, dim))
}
