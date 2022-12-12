package main

import (
	"container/heap"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
	"github.com/fatih/color"
)

var vis bool

// An Item is something we manage in a priority queue.
type Item struct {
	value    Pos // The value of the item; arbitrary.
	priority int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type Pos struct {
	X, Y int
}

func (p *Pos) mut(m Pos) Pos {
	return Pos{p.X + m.X, p.Y + m.Y}
}

func (p *Pos) eq(m Pos) bool {
	return p.X == m.X && p.Y == m.Y
}

func (p *Pos) dist(to Pos) int {
	return util.AbsDiffInt(p.X, to.X) + util.AbsDiffInt(p.Y, to.Y)
}

func CharMatrix(lines []string, lookup map[string]bool) (int, int, *[][]string, *map[string]Pos) {
	height := len(lines)
	width := len(lines[0])
	res := make([][]string, height)
	resolve := make(map[string]Pos)
	for y := 0; y < height; y++ {
		res[y] = make([]string, width)
	}
	for y, line := range lines {
		for x, s := range line {
			s := string(s)
			res[y][x] = s
			if ok := lookup[s]; ok {
				resolve[s] = Pos{x, y}
			}
		}
	}
	return width, height, &res, &resolve
}

func canStep(from, to string, up bool) bool {
	if to == "E" {
		to = "z"
	}
	if from == "S" {
		from = "a"
	}
	if from == "E" {
		from = "z"
	}
	delta := int(to[0]) - int(from[0])
	if up {
		return delta <= 1
	} else {
		return delta == -1 || delta == 0 || delta > 0
	}

}

func candidates(mat *[][]string, visited *map[Pos]bool, from Pos, width, height int, up bool) *[]Pos {
	res := make([]Pos, 0, 4)
	curH := (*mat)[from.Y][from.X]
	steps := []Pos{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	for _, step := range steps {
		next := from.mut(step)
		if 0 <= next.X && next.X < width && 0 <= next.Y && next.Y < height {
			candidateH := (*mat)[next.Y][next.X]
			wasHere := (*visited)[next]
			if canStep(curH, candidateH, up) && !wasHere {
				res = append(res, next)
			}
		}
	}
	return &res
}

func sample() []string {
	d := `Sabqponm
abcryxxl
accszExk
acctuvwj
abdefghi`
	return strings.Split(d, "\n")
}

var visitedPlace = color.New(color.FgWhite, color.Bold).SprintFunc()
var normalPlace = color.New(color.FgCyan).SprintFunc()

func visualize(mat *[][]string, visited *map[Pos]bool, edgePositions *map[string]Pos) {
	if !vis {
		return
	}
	util.ClearTerm()
	for y, line := range *mat {
		lineStr := ""
		for x, s := range line {
			pos := Pos{x, y}
			var char string
			if pos.eq((*edgePositions)["S"]) {
				char = "S"
			} else if pos.eq((*edgePositions)["E"]) {
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
		fmt.Printf("%s\n", lineStr)
	}
	time.Sleep(25 * time.Millisecond)
}

func partA(mat *[][]string, edgePositions *map[string]Pos, width, height int) int {
	done := goaocd.Duration("Part A")
	done()

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	start := (*edgePositions)["S"]
	end := (*edgePositions)["E"]
	heap.Push(&pq, &Item{value: start, priority: 0})
	visited := make(map[Pos]bool)

	for {
		next := heap.Pop(&pq).(*Item)
		if next.value.eq(end) {
			return next.priority
		}
		c := candidates(mat, &visited, next.value, width, height, true)
		if len(*c) == 0 && pq.Len() == 0 {
			panic("Nowhere to go!")
		}
		for _, i := range *c {
			visited[i] = true
			heap.Push(&pq, &Item{value: i, priority: next.priority + 1})
		}
		visualize(mat, &visited, edgePositions)
	}
}

func partB(mat *[][]string, edgePositions *map[string]Pos, width, height int) int {
	done := goaocd.Duration("Part B")
	done()

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	start := (*edgePositions)["E"]

	heap.Push(&pq, &Item{value: start, priority: 0})
	visited := make(map[Pos]bool)

	for {
		next := heap.Pop(&pq).(*Item)
		pos := next.value
		char := (*mat)[pos.Y][pos.X]
		if char == "a" {
			return next.priority
		}
		c := candidates(mat, &visited, next.value, width, height, false)
		if len(*c) == 0 && pq.Len() == 0 {
			panic("Nowhere to go!")
		}
		for _, i := range *c {
			visited[i] = true
			heap.Push(&pq, &Item{value: i, priority: next.priority + 1})
		}
		visualize(mat, &visited, edgePositions)
	}
}

func main() {
	vis = os.Getenv("VIS") != ""
	// lines := sample()
	lines := goaocd.Lines()
	searchFor := map[string]bool{"S": true, "E": true}

	width, height, mat, edgePositions := CharMatrix(lines, searchFor)
	fmt.Printf("Part A: %d\n", partA(mat, edgePositions, width, height))
	fmt.Printf("Part B: %d\n", partB(mat, edgePositions, width, height))
}
