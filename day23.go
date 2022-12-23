package main

import (
	"fmt"
	"strings"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/fatih/color"

	"github.com/Olegas/goaocd"
)

type Move struct {
	dir      goaocd.Pos
	adjacent []goaocd.Pos
}

func NewMoveMap() []Move {
	return []Move{
		{
			// north
			goaocd.Pos{0, -1},
			[]goaocd.Pos{{-1, -1}, {0, -1}, {1, -1}},
		},
		{
			// south
			goaocd.Pos{0, 1},
			[]goaocd.Pos{{-1, 1}, {0, 1}, {1, 1}},
		},
		{
			// west
			goaocd.Pos{-1, 0},
			[]goaocd.Pos{{-1, -1}, {-1, 0}, {-1, 1}},
		},
		{
			// east
			goaocd.Pos{1, 0},
			[]goaocd.Pos{{1, -1}, {1, 0}, {1, 1}},
		},
	}
}

var around = []goaocd.Pos{
	{-1, -1}, {0, -1}, {1, -1},
	{-1, 0}, {1, 0},
	{-1, 1}, {0, 1}, {1, 1},
}

type Elve struct {
	pos          goaocd.Pos
	moveMap      []Move
	moveIdx      int
	proposedMove *goaocd.Pos
}

func NewElve(at goaocd.Pos) Elve {
	return Elve{pos: at, moveMap: NewMoveMap()}
}

func (e *Elve) commit() {
	e.pos = goaocd.Pos{X: e.proposedMove.X, Y: e.proposedMove.Y}
}

func (e *Elve) prepareToNextRound() {
	e.proposedMove = nil
	e.moveIdx = (e.moveIdx + 1) % len(e.moveMap)
}

func (e *Elve) proposeMove(world *map[goaocd.Pos]*Elve) *goaocd.Pos {
	nM := len(e.moveMap)
	for i := 0; i < nM; i++ {
		j := (e.moveIdx + i) % nM
		m := e.moveMap[j]
		free := true
		for _, p := range m.adjacent {
			n := e.pos.Mut(p)
			_, occupied := (*world)[n]
			if occupied {
				free = false
				break
			}
		}
		if free {
			np := e.pos.Mut(m.dir)
			e.proposedMove = &np
			return e.proposedMove
		}
	}
	return nil
}

func (e *Elve) needToMove(world *map[goaocd.Pos]*Elve) bool {
	for _, p := range around {
		n := e.pos.Mut(p)
		_, occupied := (*world)[n]
		if occupied {
			return true
		}
	}
	return false
}

func sample() []string {
	d := `....#..
..###.#
#...#.#
.#...##
#.###..
##.#.##
.#..#..`
	return strings.Split(d, "\n")
}

func display(world *map[goaocd.Pos]*Elve) {
	dim := dimensions(world)
	padding := 2

	for y := dim.MinY - padding; y <= dim.MaxY+padding; y++ {
		for x := dim.MinX - padding; x <= dim.MaxX+padding; x++ {
			k := goaocd.Pos{X: x, Y: y}
			_, occupied := (*world)[k]
			s := "#"
			if !occupied {
				s = "."
			}
			fmt.Print(s)
		}
		fmt.Print("\n")
	}
}

var nextPlace = color.New(color.FgWhite, color.Bold).PrintFunc()

func display2(world *map[goaocd.Pos]*Elve, world2 *map[goaocd.Pos][]*Elve) {
	dim := dimensions(world)
	padding := 2

	for y := dim.MinY - padding; y <= dim.MaxY+padding; y++ {
		for x := dim.MinX - padding; x <= dim.MaxX+padding; x++ {
			k := goaocd.Pos{X: x, Y: y}
			_, occupied := (*world)[k]
			s := "#"
			if !occupied {
				s = "."
			}

			l, occupied := (*world2)[k]
			if occupied {
				nextPlace(fmt.Sprintf("%d", len(l)))
			} else {
				fmt.Print(s)
			}
		}
		fmt.Print("\n")
	}
}

func simulate(world *map[goaocd.Pos]*Elve, maxRounds *int) int {
	i := 0
	for {
		nextWorld := make(map[goaocd.Pos][]*Elve)
		//display(world)
		//fmt.Scanln()

		someOneNeedToMove := false
		for _, e := range *world {
			if e.needToMove(world) {
				someOneNeedToMove = true
				np := e.proposeMove(world)
				if np != nil {
					anp, present := nextWorld[*np]
					if !present {
						anp = make([]*Elve, 0)
					}
					nextWorld[*np] = append(anp, e)
				}
			}
		}

		//display2(world, &nextWorld)
		//fmt.Scanln()

		for _, es := range nextWorld {
			if len(es) == 1 {
				e := es[0]
				// remove Elve from old position
				_, present := (*world)[e.pos]
				if !present {
					panic("Not present")
				}
				delete(*world, e.pos)
				// move
				e.commit()
				(*world)[e.pos] = e
			}
		}

		for p := range *world {
			k := (*world)[p]
			k.prepareToNextRound()
		}

		i++

		if (maxRounds == nil && !someOneNeedToMove) || (maxRounds != nil && *maxRounds == i) {
			// display(world)
			return i
		}

		//display(world)
		//fmt.Scanln()
	}
}

type Dimensions struct {
	MinX, MaxX, MinY, MaxY int
}

func dimensions(world *map[goaocd.Pos]*Elve) Dimensions {
	xCoords := []int{}
	yCoords := []int{}
	for c := range *world {
		xCoords = append(xCoords, c.X)
		yCoords = append(yCoords, c.Y)
	}
	minX, maxX := util.MinMax(xCoords)
	minY, maxY := util.MinMax(yCoords)

	return Dimensions{MinX: minX, MaxX: maxX, MinY: minY, MaxY: maxY}
}

func partA(world *map[goaocd.Pos]*Elve) int {
	done := goaocd.Duration("Part A")
	defer done()

	mR := 10
	simulate(world, &mR)

	dim := dimensions(world)
	accu := 0
	for x := dim.MinX; x <= dim.MaxX; x++ {
		for y := dim.MinY; y <= dim.MaxY; y++ {
			p := goaocd.Pos{x, y}
			_, occupied := (*world)[p]
			if !occupied {
				accu++
			}
		}
	}

	return accu
}

func partB(world *map[goaocd.Pos]*Elve) int {
	done := goaocd.Duration("Part A")
	defer done()

	return simulate(world, nil)
}

func loadWorld(lines []string) *map[goaocd.Pos]*Elve {
	res := make(map[goaocd.Pos]*Elve)
	for y, line := range lines {
		for x, s := range line {
			if string(s) == "#" {
				p := goaocd.Pos{X: x, Y: y}
				e := NewElve(p)
				res[p] = &e
			}
		}
	}
	return &res
}

func main() {
	lines := sample()
	lines = goaocd.Lines()

	world := loadWorld(lines)
	fmt.Printf("Part A: %d\n", partA(world))

	world = loadWorld(lines)
	fmt.Printf("Part B: %d\n", partB(world))
}
