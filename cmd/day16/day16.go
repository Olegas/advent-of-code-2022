package main

// NOT WORKING...
// Correct solution for sample but incorrect for actual input

import (
	"container/heap"
	"fmt"
	"regexp"
	"strings"

	"github.com/Olegas/goaocd"
)

type Valve struct {
	name string
	rate int
}

type Graph struct {
	nodes map[string]Valve
	links map[string][]string
}

type Item struct {
	pathLen int
	node    string
	index   int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].pathLen > pq[j].pathLen
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

func sample() []string {
	d := `Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
Valve BB has flow rate=13; tunnels lead to valves CC, AA
Valve CC has flow rate=2; tunnels lead to valves DD, BB
Valve DD has flow rate=20; tunnels lead to valves CC, AA, EE
Valve EE has flow rate=3; tunnels lead to valves FF, DD
Valve FF has flow rate=0; tunnels lead to valves EE, GG
Valve GG has flow rate=0; tunnels lead to valves FF, HH
Valve HH has flow rate=22; tunnel leads to valve GG
Valve II has flow rate=0; tunnels lead to valves AA, JJ
Valve JJ has flow rate=21; tunnel leads to valve II`
	return strings.Split(d, "\n")
}

func partA(graph *Graph, nonEmptyNodes *[]string) int {
	done := goaocd.Duration("Part A")
	defer done()

	max := 0

	type i struct {
		timeLeft      int
		totalPressure int
		current       string
		seen          map[string]bool
		path          string
	}

	leftToVisit := func(visited map[string]bool) []string {
		res := make([]string, 0, len(*nonEmptyNodes)-len(visited))
		for _, n := range *nonEmptyNodes {
			ok := visited[n]
			if !ok {
				res = append(res, n)
			}
		}
		return res
	}

	candidates := []i{{
		seen:     map[string]bool{"AA": true},
		current:  "AA",
		timeLeft: 30,
		path:     "AA",
	}}

	for {
		newItemsAdded := false
		for j := 0; j < len(candidates); j++ {
			c := candidates[j]
			variants := leftToVisit(c.seen)
			for _, n := range variants {
				leg := fmt.Sprintf("%s,%s", c.current, n)
				d, ok := distances[leg]
				if !ok {
					panic(fmt.Sprintf("Distance for leg %s is not calculated", leg))
				}
				timeLeftAfterGoAndOpen := c.timeLeft - d - 1
				if timeLeftAfterGoAndOpen >= 0 {
					newSeen := make(map[string]bool)
					for k, v := range c.seen {
						newSeen[k] = v
					}
					newSeen[n] = true
					ni := i{
						seen:          newSeen,
						current:       n,
						timeLeft:      timeLeftAfterGoAndOpen,
						path:          c.path + "-" + n,
						totalPressure: c.totalPressure + timeLeftAfterGoAndOpen*graph.nodes[n].rate,
					}
					candidates = append(candidates, ni)
					if ni.totalPressure > max {
						max = ni.totalPressure
					}
				}
			}
		}
		if !newItemsAdded {
			break
		}
	}

	return max
}

func partB() int {
	done := goaocd.Duration("Part B")
	defer done()

	return 0
}

func loadPuzzle() *Graph {
	lines := sample()
	// lines = goaocd.Lines()
	graph := Graph{nodes: make(map[string]Valve), links: map[string][]string{}}
	re := regexp.MustCompile("[A-Z]{2,2}")
	for _, line := range lines {
		var name string
		var rate int
		pts := strings.Split(line, "; ")
		n, err := fmt.Sscanf(pts[0], "Valve %2s has flow rate=%d", &name, &rate)
		if err != nil {
			panic(err)
		}
		if n != 2 {
			panic(fmt.Sprintf("Failed to parse %s", line))
		}
		node := Valve{name: name, rate: rate}
		graph.nodes[name] = node
		if graph.links[name] == nil {
			graph.links[name] = make([]string, 0)
		}

		links := re.FindAllString(pts[1], -1)
		graph.links[name] = append(graph.links[name], links...)
	}
	return &graph
}

var distances = map[string]int{}

func path(graph *Graph, from, to string) int {

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{pathLen: 0, node: from})

	seen := map[string]bool{}
	seen[from] = true

	for {
		item := pq.Pop().(*Item)
		next := graph.links[item.node]
		for _, n := range next {
			if n == to {
				return item.pathLen + 1
			}
			_, ok := seen[n]
			if ok {
				continue
			}
			seen[n] = true
			heap.Push(&pq, &Item{pathLen: item.pathLen + 1, node: n})
		}
	}
}

func precacheDistances(graph *Graph, nonZeroNodes *[]string) {
	for _, a := range *nonZeroNodes {
		for _, b := range *nonZeroNodes {
			if a == b || b == "AA" {
				continue
			}
			leg := fmt.Sprintf("%s,%s", a, b)
			_, ok := distances[leg]
			if ok {
				continue
			}
			d := path(graph, a, b)
			distances[leg] = d
			distances[leg] = d
		}
	}
}

func getNonZeroNodes(graph *Graph) *[]string {
	nonZeroNodes := []string{"AA"}
	for name, node := range graph.nodes {
		if node.rate != 0 {
			nonZeroNodes = append(nonZeroNodes, name)
		}
	}
	return &nonZeroNodes
}

func main() {
	graph := loadPuzzle()
	nonZeroNodes := getNonZeroNodes(graph)
	precacheDistances(graph, nonZeroNodes)

	fmt.Printf("Part A: %d\n", partA(graph, nonZeroNodes))
	//fmt.Printf("Part B: %d\n", partB())
}
