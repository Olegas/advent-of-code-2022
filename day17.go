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

type State struct {
	height  int
	stopped int
}

type Key struct {
	depth [7]int
	fig   int
	move  int
}

func sample() []string {
	d := `>>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>`
	return strings.Split(d, "")
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

func calcDepthMap(screen *[]uint8, top int) [7]int {
	res := [7]int{}
	for i := 0; i < 7; i++ {
		m := uint8(1 << i)
		for j := top; j >= 0; j-- {
			if (*screen)[j]&m != 0 {
				res[i] = top - j
				break
			}
		}
	}
	return res
}

func simulate(maxStoppedRocks int) int64 {
	checkPatterns := true
	highestPoint := 0
	width := 7
	heightInCycles := 0
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
	seen := make(map[Key]State)

	filledLineValue := (1 << width) - 1
	screen[0] = uint8(filledLineValue)

	figI := 0
	moveI := 0
	lenFig := len(figures)
	lenMoves := len(movesPos)

	stoppedCountLeft := maxStoppedRocks
	for {
		if stoppedCountLeft == 0 {
			break
		}
		fig := figures[figI]
		figI = (figI + 1) % lenFig

		fig.pos.X = newFigOffset.X
		fig.pos.Y = int(highestPoint) + newFigOffset.Y

		if len(screen) < fig.pos.Y+fig.h+1 {
			newScreen := make([]uint8, (fig.pos.Y+fig.h)*100000)
			copy(newScreen, screen)
			screen = newScreen
		}

		for {
			move := movesPos[moveI]
			moveI = (moveI + 1) % lenMoves

			// move
			if figCanMove(&fig, &move, width, &screen) {
				fig.pos.X += move.X
			}

			// fall
			if figCanMove(&fig, &moveDown, width, &screen) {
				fig.pos.Y--
			} else {
				drawFig(&fig, &screen)
				figHigh := fig.pos.Y + fig.h - 1
				highestPoint = util.Max(highestPoint, figHigh)
				stoppedCountLeft--
				if checkPatterns {
					depthMap := calcDepthMap(&screen, highestPoint)
					key := Key{fig: figI, move: moveI, depth: depthMap}
					prevS, ok := seen[key]
					if ok {
						diffHeight := highestPoint - prevS.height
						cycleLength := prevS.stopped - stoppedCountLeft

						countCylcesLeft := stoppedCountLeft / cycleLength
						heightInCycles = diffHeight * countCylcesLeft

						stoppedCountLeft = stoppedCountLeft % cycleLength
						checkPatterns = false
					} else {
						state := State{height: highestPoint, stopped: stoppedCountLeft}
						seen[key] = state
					}
				}
				break
			}
		}
	}

	return int64(heightInCycles) + int64(highestPoint)
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
