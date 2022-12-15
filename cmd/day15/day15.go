package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
)

func sample() []string {
	d := `Sensor at x=2, y=18: closest beacon is at x=-2, y=15
Sensor at x=9, y=16: closest beacon is at x=10, y=16
Sensor at x=13, y=2: closest beacon is at x=15, y=3
Sensor at x=12, y=14: closest beacon is at x=10, y=16
Sensor at x=10, y=20: closest beacon is at x=10, y=16
Sensor at x=14, y=17: closest beacon is at x=10, y=16
Sensor at x=8, y=7: closest beacon is at x=2, y=10
Sensor at x=2, y=0: closest beacon is at x=2, y=10
Sensor at x=0, y=11: closest beacon is at x=2, y=10
Sensor at x=20, y=14: closest beacon is at x=25, y=17
Sensor at x=17, y=20: closest beacon is at x=21, y=22
Sensor at x=16, y=7: closest beacon is at x=15, y=3
Sensor at x=14, y=3: closest beacon is at x=15, y=3
Sensor at x=20, y=1: closest beacon is at x=15, y=3`
	return strings.Split(d, "\n")
}

type Pair struct {
	sensor goaocd.Pos
	beacon goaocd.Pos
}

func partA(vis map[goaocd.Pos]string, pairs []Pair, watchLine int) int {
	done := goaocd.Duration("Part A")
	defer done()

	freeSpaces := make(map[goaocd.Pos]string)
	for _, pair := range pairs {
		sensorX := pair.sensor.X
		sensorY := pair.sensor.Y
		dist := pair.sensor.ManhattanDist(pair.beacon)
		if sensorY-dist < watchLine && watchLine < sensorY+dist {
			distToLine := util.AbsDiffInt(sensorY, watchLine)
			distFromEdge := dist - distToLine
			countFree := distFromEdge*2 + 1
			steps := int(math.Floor(float64(countFree) / 2.0))
			for dx := -steps; dx <= steps; dx++ {
				p := goaocd.Pos{X: sensorX + dx, Y: watchLine}
				// Check for beacon or sensor in this point
				_, ok := vis[p]
				if !ok {
					freeSpaces[p] = "#"
				}
			}
		}
	}

	return len(freeSpaces)
}

type Range struct {
	a, b int
}

func Clamp(a, b, to int) (int, int) {
	return util.Max(0, a), util.Min(b, to)
}

func partB(vis map[goaocd.Pos]string, pairs []Pair, maxCoord int) int64 {
	done := goaocd.Duration("Part B")
	defer done()

	var ranges = make([][]Range, maxCoord+1)

	putRange := func(y int, r Range) {
		p := ranges[y]
		if p == nil {
			p = make([]Range, 0)
		}
		ranges[y] = append(p, r)
	}

	for _, pair := range pairs {
		sensorX := pair.sensor.X
		sensorY := pair.sensor.Y
		dist := pair.sensor.ManhattanDist(pair.beacon)

		for dy := -dist; dy <= 0; dy++ {
			y := sensorY + dy
			dx := dy + dist
			a, b := Clamp(sensorX-dx, sensorX+dx, maxCoord)
			if 0 <= y && y <= maxCoord {
				putRange(y, Range{a, b})
			}

			if dy != 0 {
				y := sensorY - dy
				if 0 <= y && y <= maxCoord {
					putRange(y, Range{a, b})
				}
			}
		}
	}

	for y, r := range ranges {
		sort.SliceStable(r, func(i, j int) bool {
			a := r[i]
			b := r[j]
			return a.a < b.a
		})
		reduced := []*Range{&r[0]}
		for _, i := range r[1:] {
			t := reduced[len(reduced)-1]
			if (t.a <= i.a && i.a <= t.b) || (t.b == i.a-1) {
				t.b = util.Max(t.b, i.b)
			} else {
				reduced = append(reduced, &Range{i.a, i.b})
			}
		}
		if len(reduced) != 1 {
			return int64(reduced[0].b+1)*4000000 + int64(y)
		}
	}

	return 0
}

func main() {
	lines := sample()
	lines = goaocd.Lines()
	vis := make(map[goaocd.Pos]string)
	pairs := make([]Pair, 0)
	for _, line := range lines {
		items := strings.Split(line, ": ")
		sensor := items[0][10:]
		var x, y int
		n, err := fmt.Sscanf(sensor, "x=%d, y=%d", &x, &y)
		if err != nil {
			panic(err)
		}
		if n != 2 {
			panic(fmt.Sprintf("Failed to parse sensor position %s", sensor))
		}
		sensorPos := goaocd.Pos{X: x, Y: y}
		vis[sensorPos] = "S"

		beacon := items[1][21:]
		n, err = fmt.Sscanf(beacon, "x=%d, y=%d", &x, &y)
		if err != nil {
			panic(err)
		}
		if n != 2 {
			panic(fmt.Sprintf("Failed to parse sensor position %s", sensor))
		}
		beaconPos := goaocd.Pos{X: x, Y: y}
		vis[beaconPos] = "B"
		pairs = append(pairs, Pair{sensor: sensorPos, beacon: beaconPos})
	}

	// Sample: 10
	// Actual: 2000000
	fmt.Printf("Part A: %d\n", partA(vis, pairs, 2000000))

	// Sample: 20
	// Actual: 4000000
	fmt.Printf("Part A: %d\n", partB(vis, pairs, 4000000))
}
