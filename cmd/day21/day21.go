package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Olegas/goaocd"
)

func sample() []string {
	d := `root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32`
	return strings.Split(d, "\n")
}

type Node struct {
	name      string
	value     int
	operation string
}

type Graph struct {
	nodes map[string]*Node
	links map[string][]string
}

type Op func(int, int) int

var opMap = map[string]Op{
	"+": func(a, b int) int { return a + b },
	"-": func(a, b int) int { return a - b },
	"*": func(a, b int) int { return a * b },
	"/": func(a, b int) int { return a / b },
}

var revOpMap = map[string][]Op{
	"+": {func(a, m int) int { return m - a }, func(m, b int) int { return m - b }},
	"-": {func(a, m int) int { return -m + a }, func(m, b int) int { return m + b }},
	"*": {func(a, m int) int { return m / a }, func(m, b int) int { return m / b }},
	"/": {func(a, m int) int { return a / m }, func(m, b int) int { return m * b }},
}

func loadGraph(lines []string) *Graph {
	g := Graph{
		nodes: make(map[string]*Node, 0),
		links: make(map[string][]string),
	}
	for _, line := range lines {
		p := strings.Split(line, ": ")
		n := Node{name: p[0]}
		g.nodes[n.name] = &n

		i, err := strconv.Atoi(p[1])
		if err != nil {
			var a, b, op string
			i, err := fmt.Sscanf(p[1], "%s %s %s", &a, &op, &b)
			if err != nil || i != 3 {
				panic(fmt.Sprintf("Failed to parse op %s", p[1]))
			}
			if g.links[n.name] == nil {
				g.links[n.name] = make([]string, 0)
			}
			g.links[n.name] = append(g.links[n.name], a, b)
			n.operation = op
		} else {
			n.value = i
		}
	}

	return &g
}

func calc(g *Graph, node string) int {
	n := g.nodes[node]
	if n.operation == "" {
		return n.value
	} else {
		l, ok := g.links[node]
		if !ok {
			panic(fmt.Sprintf("No links for %s", node))
		}
		op := n.operation
		f := opMap[op]
		return f(calc(g, l[0]), calc(g, l[1]))
	}
}

func searchForHuman(g *Graph, node string) bool {
	links, ok := g.links[node]
	if ok {
		return searchForHuman(g, links[0]) || searchForHuman(g, links[1])
	}
	return node == "humn"
}

func partA(g *Graph) int {
	done := goaocd.Duration("Part A")
	defer done()

	return calc(g, "root")
}

func solve(g *Graph, mustEq int, branch string) int {
	links := g.links[branch]
	n := g.nodes[branch]
	if n.operation == "" {
		if n.name != "humn" {
			panic("We should not be here!")
		} else {
			return mustEq
		}
	} else {
		ops := revOpMap[n.operation]
		if searchForHuman(g, links[0]) {
			right := calc(g, links[1])
			op := ops[1]
			return solve(g, op(mustEq, right), links[0])
		} else {
			left := calc(g, links[0])
			op := ops[0]
			return solve(g, op(left, mustEq), links[1])
		}
	}
}

func partB(g *Graph) int {
	done := goaocd.Duration("Part B")
	defer done()

	links := g.links["root"]
	humanAtLeft := searchForHuman(g, links[0])

	var num int
	var branch string
	if humanAtLeft {
		num = calc(g, links[1])
		branch = links[0]
	} else {
		num = calc(g, links[0])
		branch = links[1]
	}

	return solve(g, num, branch)
}

func main() {
	lines := sample()
	lines = goaocd.Lines()
	graph := loadGraph(lines)

	fmt.Printf("Part A: %d\n", partA(graph))
	fmt.Printf("Part B: %d\n", partB(graph))
}
