package main

import (
	"container/heap"
	"fmt"

	"github.com/Olegas/advent-of-code-2022/internal/day12"
	"github.com/Olegas/goaocd"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    goaocd.Pos // The value of the item; arbitrary.
	priority int        // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
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

func candidates(mat *[][]string, visited *map[goaocd.Pos]bool, from goaocd.Pos, width, height int, up bool) *[]goaocd.Pos {
	res := make([]goaocd.Pos, 0, 4)
	curH := (*mat)[from.Y][from.X]
	steps := []goaocd.Pos{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: 0, Y: -1}}
	for _, step := range steps {
		next := from.Mut(step)
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

func partA(mat *[][]string, edgePositions *map[string]goaocd.Pos, width, height int) int {
	done := goaocd.Duration("Part A")
	done()

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	start := (*edgePositions)["S"]
	end := (*edgePositions)["E"]
	heap.Push(&pq, &Item{value: start, priority: 0})
	visited := make(map[goaocd.Pos]bool)

	for {
		next := heap.Pop(&pq).(*Item)
		if next.value.Eq(end) {
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
		day12.Visualize(mat, &visited, edgePositions)
	}
}

func partB(mat *[][]string, edgePositions *map[string]goaocd.Pos, width, height int) int {
	done := goaocd.Duration("Part B")
	done()

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	start := (*edgePositions)["E"]

	heap.Push(&pq, &Item{value: start, priority: 0})
	visited := make(map[goaocd.Pos]bool)

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
		day12.Visualize(mat, &visited, edgePositions)
	}
}

func main() {
	searchFor := map[string]bool{"S": true, "E": true}
	width, height, mat, edgePositions := goaocd.CharMatrix(&searchFor)
	fmt.Printf("Part A: %d\n", partA(mat, edgePositions, width, height))
	fmt.Printf("Part B: %d\n", partB(mat, edgePositions, width, height))
}
