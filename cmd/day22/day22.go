package main

import (
	"fmt"
	"github.com/Olegas/advent-of-code-2022/internal/util"
	"strings"

	"github.com/Olegas/goaocd"
)

type Me struct {
	pos    goaocd.Pos
	facing goaocd.Pos
}

type Move struct {
	steps int
	rot   string
}

/**
  A B
  C
D E
F
           /+
      A   / |
  +-----+/  |
E |  C  | B/+
  |     | /
  +-----+/
     D
*/

type Transition struct {
	face   string
	facing goaocd.Pos
}

//R: (1, 0)
//R: (0, 1), L: (0, -1)

//L: (-1, 0)
//R: (0, -1), L: (0, 1)

//U (0, -1)
//R: (1, 0), L: (-1, 0)

//D (0, 1)
//R: (1, 0), L: (-1, 0)

func (m *Me) Rotate(dir string) {
	if dir == "" {
		return
	}
	f := m.facing
	v := f.Y != 0
	mul := 1
	if !v && dir == "L" || v && dir == "R" {
		mul = -1
	}
	m.facing = goaocd.Pos{X: f.Y * mul, Y: f.X * mul}
	if util.AbsDiffInt(m.facing.X, m.facing.Y) != 1 {
		panic("Broken facing")
	}
}

var wrapAround = make(map[Me]goaocd.Pos)

type Dimensions struct {
	MinX, MaxX, MinY, MaxY int
}

func dimensions(world map[goaocd.Pos]string) Dimensions {
	xCoords := []int{}
	yCoords := []int{}
	for c := range world {
		xCoords = append(xCoords, c.X)
		yCoords = append(yCoords, c.Y)
	}
	minX, maxX := util.MinMax(xCoords)
	minY, maxY := util.MinMax(yCoords)

	return Dimensions{MinX: minX, MaxX: maxX, MinY: minY, MaxY: maxY}
}

func display(world map[goaocd.Pos]string, me Me) {
	dim := Dimensions{
		MinX: me.pos.X - 10,
		MaxX: me.pos.X + 10,
		MinY: me.pos.Y - 10,
		MaxY: me.pos.Y + 10,
	}
	padding := 2

	for y := dim.MinY - padding; y <= dim.MaxY+padding; y++ {
		for x := dim.MinX - padding; x <= dim.MaxX+padding; x++ {
			k := goaocd.Pos{X: x, Y: y}
			s, ok := world[k]

			if !ok {
				s = " "
			}
			if me.pos.Eq(k) {
				_, s = facingToInt(me.facing)
			}
			fmt.Print(s)
		}
		fmt.Print("\n")
	}

	fmt.Scanln()
}

func sample() []string {
	d := `....
....
....
...#

0R0R0R0R0L0L0L0L`
	return strings.Split(d, "\n")
}

func findNextWrapPos(world map[goaocd.Pos]string, me Me) goaocd.Pos {
	r, ok := wrapAround[me]
	if ok {
		return r
	}
	pos := me.pos
	facing := me.facing
	opposite := goaocd.Pos{X: -facing.X, Y: -facing.Y}
	for {
		n := pos.Mut(opposite)
		t, ok := world[n]
		if !ok || t == " " {
			wrapAround[me] = pos
			return pos
		}
		pos = n
	}
}

func (m *Me) Move(world map[goaocd.Pos]string, steps int) {
	for i := 0; i < steps; i++ {
		n := m.pos.Mut(m.facing)
		t, ok := world[n]
		if !ok || t == " " {
			// wrap around
			n = findNextWrapPos(world, *m)
			t, ok := world[n]
			if !ok {
				panic("Wraparound is void!!!")
			}
			if t == "#" {
				// wall after warp around, stop
				break
			}
		} else if t == "#" {
			break
		}
		// step
		m.pos = n
		// display(world, *m)
	}
}

func facingToInt(pos goaocd.Pos) (int, string) {
	if pos.X == 1 && pos.Y == 0 {
		return 0, ">"
	} else if pos.X == 0 && pos.Y == 1 {
		return 1, "v"
	} else if pos.X == -1 && pos.Y == 0 {
		return 2, "<"
	} else if pos.X == 0 && pos.Y == -1 {
		return 3, "^"
	}
	panic("Incorrect facing")
}

func partA(world map[goaocd.Pos]string, moves []Move, me Me) int {
	done := goaocd.Duration("Part A")
	defer done()

	for _, m := range moves {
		fmt.Printf("%d%s\n", m.steps, m.rot)
		me.Move(world, m.steps)
		me.Rotate(m.rot)
		// display(world, me)
	}
	fi, _ := facingToInt(me.facing)
	fmt.Printf("%d * 1000 + %d * 4 + %d", me.pos.Y+1, me.pos.X+1, fi)
	return (me.pos.Y+1)*1000 + (me.pos.X+1)*4 + fi
}

func partB() int {
	done := goaocd.Duration("Part B")
	defer done()

	return 0
}

func loadMap(lines []string) (map[goaocd.Pos]string, []Move, Me) {
	ret := make(map[goaocd.Pos]string)
	me := Me{}
	moves := []Move{}
	mapMode := true
	meFound := false
	for y, line := range lines {
		if line == "" {
			mapMode = false
			continue
		}
		if mapMode {
			for x, s := range line {
				p := goaocd.Pos{X: x, Y: y}
				if !meFound && string(s) == "." {
					me.pos = p
					me.facing = goaocd.Pos{X: 1, Y: 0}
					meFound = true
				}
				ret[p] = string(s)
			}
		} else {
			dAccu := ""
			for _, s := range line {
				if util.IsDigit(byte(s)) {
					dAccu += string(s)
				} else {
					m := Move{steps: util.Atoi(dAccu), rot: string(s)}
					moves = append(moves, m)
					dAccu = ""
				}
			}
			if dAccu != "" {
				m := Move{steps: util.Atoi(dAccu), rot: ""}
				moves = append(moves, m)
			}
		}

	}
	return ret, moves, me
}

func main() {
	lines := sample()
	lines = goaocd.Lines()
	world, moves, me := loadMap(lines)

	fmt.Printf("Part A: %d\n", partA(world, moves, me))
	fmt.Printf("Part B: %d\n", partB())
}
