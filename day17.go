package main

import (
	"fmt"
	"strings"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
)

type Fig struct {
	pos    goaocd.Pos
	points []goaocd.Pos
	h      int
	w      int
}

var figures = []Fig{
	{points: []goaocd.Pos{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, h: 1, w: 4},
	{points: []goaocd.Pos{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}}, h: 3, w: 3},
	{points: []goaocd.Pos{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}}, h: 3, w: 3},
	{points: []goaocd.Pos{{0, 0}, {0, 1}, {0, 2}, {0, 3}}, h: 4, w: 1},
	{points: []goaocd.Pos{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, h: 2, w: 2},
}

func sample() []string {
	d := `>>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>`
	return strings.Split(d, "")
}

type Stack struct {
	dumped       int64
	highestPoint goaocd.Pos
	w            int
}

func figCanMove(fig *Fig, move *goaocd.Pos, w int, screen *[]uint8) bool {
	px := fig.pos.X + move.X
	py := fig.pos.Y + move.Y
	for _, i := range fig.points {
		x := px + i.X
		if x < 0 || x >= w {
			return false
		}
		y := py + i.Y
		if (*screen)[y]&(1<<x) != 0 {
			return false
		}

	}
	return true
}

func drawFig(fig *Fig, screen *[]uint8) {
	for _, i := range fig.points {
		x := fig.pos.X + i.X
		y := fig.pos.Y + i.Y
		(*screen)[y] |= 1 << x
	}
}

func simulate(maxStoppedRocks int) int64 {
	s := Stack{w: 7}

	moves := sample()
	moves = strings.Split(strings.TrimRight(goaocd.Input(), "\n"), "")
	movesPos := make([]goaocd.Pos, len(moves))
	for idx, s := range moves {
		p := goaocd.Pos{X: -1, Y: 0}
		if s == ">" {
			p.X = 1
		}
		movesPos[idx] = p
	}
	newFigOffset := goaocd.Pos{X: 2, Y: 4}
	moveDown := goaocd.Pos{X: 0, Y: -1}
	screen := make([]uint8, 1)
	filledLineValue := (1 << s.w) - 1
	screen[0] = uint8(filledLineValue)

	figI := 0
	moveI := 0
	lenFig := len(figures)
	lenMoves := len(movesPos)

	stoppedCount := 0
	for {
		if stoppedCount == maxStoppedRocks {
			break
		}
		fig := figures[figI]
		figI = (figI + 1) % lenFig

		fig.pos.X = newFigOffset.X
		fig.pos.Y = s.highestPoint.Y + newFigOffset.Y

		if len(screen) < fig.pos.Y+fig.h+1 {
			newScreen := make([]uint8, (fig.pos.Y+fig.h)*100000)
			copy(newScreen, screen)
			screen = newScreen
		}

		for {
			move := movesPos[moveI]
			moveI = (moveI + 1) % lenMoves

			// move
			if figCanMove(&fig, &move, s.w, &screen) {
				fig.pos.X += move.X
			}

			// fall
			if figCanMove(&fig, &moveDown, s.w, &screen) {
				fig.pos.Y--
			} else {
				drawFig(&fig, &screen)
				figHigh := fig.pos.Y + fig.h - 1
				s.highestPoint.Y = util.Max(s.highestPoint.Y, figHigh)
				stoppedCount++

				// Detect filled lines to dump screen down
				for i := 0; i < fig.h; i++ {
					y := fig.pos.Y + i
					if screen[y] == uint8(filledLineValue) {
						screen = screen[y:]
						s.highestPoint.Y -= y
						s.dumped += int64(y)
						break
					}
				}

				if stoppedCount%100000000 == 0 {
					fmt.Printf("Stopped count %d (%d/%d)\n", stoppedCount, len(screen), cap(screen))
				}
				break
			}
		}
	}

	return s.dumped + int64(s.highestPoint.Y)
}

func partA() int64 {
	done := goaocd.Duration("Part A")
	defer done()

	return simulate(2022)
}

func partB() int64 {
	done := goaocd.Duration("Part B")
	defer done()

	return simulate(1000000000000)
}

func main() {
	fmt.Printf("Part A: %d\n", partA())
	fmt.Printf("Part B: %d\n", partB())
}
