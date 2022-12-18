package main

import (
	"fmt"
	"strings"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
	mapset "github.com/deckarep/golang-set"
)

var around = []Pos3D{
	{1, 0, 0}, {-1, 0, 0},
	{0, 1, 0}, {0, -1, 0},
	{0, 0, 1}, {0, 0, -1},
}

type Pos3D struct {
	x, y, z int
}

func NewPos3D(s string) Pos3D {
	p := strings.Split(s, ",")
	x := util.Atoi(p[0])
	y := util.Atoi(p[1])
	z := util.Atoi(p[2])
	return Pos3D{x, y, z}
}

func (p *Pos3D) Mut(b Pos3D) Pos3D {
	return Pos3D{x: p.x + b.x, y: p.y + b.y, z: p.z + b.z}
}

func (p *Pos3D) Eq(b Pos3D) bool {
	return p.x == b.x && p.y == b.y && p.z == b.z
}

type Cube struct {
	pos       Pos3D
	freeFaces int
}

func sample() []string {
	d := `2,2,2
1,2,2
3,2,2
2,1,2
2,3,2
2,2,1
2,2,3
2,2,4
2,2,6
1,2,5
3,2,5
2,1,5
2,3,5`
	return strings.Split(d, "\n")
}

func countFaces(data *[]*Cube) int {
	posMap := make(map[Pos3D]bool)
	for _, c := range *data {
		posMap[c.pos] = true
	}

	for _, c := range *data {
		for _, a := range around {
			n := c.pos.Mut(a)
			_, ok := posMap[n]
			if ok {
				c.freeFaces--
			}
		}
	}

	accu := 0
	for _, c := range *data {
		accu += c.freeFaces
	}

	return accu
}

func partA(data *[]*Cube) int {
	done := goaocd.Duration("Part A")
	defer done()

	return countFaces(data)
}

func partB(data *[]*Cube, partA int) int {
	done := goaocd.Duration("Part B")
	defer done()

	var minX, minY, minZ int = 100, 100, 100
	var maxX, maxY, maxZ int
	posMap := make(map[Pos3D]bool)
	figure := mapset.NewSet()
	for _, c := range *data {
		figure.Add(c.pos)

		posMap[c.pos] = true
		minX = util.Min(minX, c.pos.x)
		minY = util.Min(minY, c.pos.y)
		minZ = util.Min(minZ, c.pos.z)

		maxX = util.Max(maxX, c.pos.x)
		maxY = util.Max(maxY, c.pos.y)
		maxZ = util.Max(maxZ, c.pos.z)
	}

	allSpace := mapset.NewSet()
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				allSpace.Add(Pos3D{x, y, z})
			}
		}
	}

	// Here we have all cubes which is not present in initial set of lava
	// If we start from edges and collect clusters of cubes, we can
	// exclude them and get only those, which is inside lava droplet
	diff := allSpace.Difference(figure)

	// Collect outer shell and remove it from difference collected above
	// in diff wi will have only cubes, which was inside of initial lava droplet
	outerShell := mapset.NewSet()
	candidates := []Pos3D{{minX, minY, minZ}}
	for ok := true; ok; ok = len(candidates) > 0 {
		cur := candidates[0]
		candidates = candidates[1:]
		for _, v := range around {
			p := cur.Mut(v)
			if diff.Contains(p) && !outerShell.Contains(p) {
				diff.Remove(p)
				outerShell.Add(p)
				candidates = append(candidates, p)
			}
		}
	}

	// Collect clusters of air cubes and count partA for each of them
	surface := 0
	for ok := true; ok; ok = diff.Cardinality() > 0 {
		rand := diff.Pop().(Pos3D)
		cluster := mapset.NewSet()
		cluster.Add(rand)
		candidates = []Pos3D{rand}
		for ok := true; ok; ok = len(candidates) > 0 {
			cur := candidates[0]
			candidates = candidates[1:]
			for _, v := range around {
				p := cur.Mut(v)
				if diff.Contains(p) && !cluster.Contains(p) {
					diff.Remove(p)
					cluster.Add(p)
					candidates = append(candidates, p)
				}
			}
		}

		cubes := []*Cube{}
		for _, p := range cluster.ToSlice() {
			p := p.(Pos3D)
			cubes = append(cubes, &Cube{freeFaces: 6, pos: p})
		}
		surface += countFaces(&cubes)
	}

	/*
		// Debug with OpenSCAD ;)
		builder := strings.Builder{}
		for _, p := range figure.ToSlice() {
			p := p.(Pos3D)
			builder.WriteString(fmt.Sprintf("translate([%d,%d,%d]) cube(1);\n", p.x, p.y, p.z))
		}
		os.WriteFile("./cubes.scad", []byte(builder.String()), os.ModePerm)
	*/

	// Result is part A result of initial data minus accumulated surface of inner cubes
	return partA - surface
}

func main() {
	lines := sample()
	lines = goaocd.Lines()
	data := make([]*Cube, len(lines))
	for idx, line := range lines {
		data[idx] = &Cube{pos: NewPos3D(line), freeFaces: 6}
	}

	a := partA(&data)
	fmt.Printf("Part A: %d\n", a)
	fmt.Printf("Part B: %d\n", partB(&data, a))
}
