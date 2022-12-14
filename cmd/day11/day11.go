package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
)

type Monkey struct {
	items           []int
	op              string
	arg             int
	divisor         int
	trueTo          int
	falseTo         int
	madeInspections int
}

type PassTo struct {
	wl, to int
}

var supermodulo int

func (m *Monkey) Turn(breakThings bool) chan *PassTo {
	c := make(chan *PassTo)

	go func() {
		for ok := len(m.items) > 0; ok; ok = len(m.items) > 0 {
			item := m.items[0]
			m.items = m.items[1:]
			var wl int
			if breakThings {
				wl = m.inspect(item) % supermodulo
			} else {
				wl = m.getBored(m.inspect(item))
			}
			if m.test(wl) {
				c <- &PassTo{wl, m.trueTo}
			} else {
				c <- &PassTo{wl, m.falseTo}
			}
		}
		c <- nil
	}()

	return c
}

func (m *Monkey) inspect(item int) int {
	var worryLevel int = item
	switch m.op {
	case "+":
		worryLevel += m.arg
	case "pow":
		if m.arg != 2 {
			panic("Pow is not 2")
		}
		worryLevel = item * item
	case "*":
		worryLevel *= m.arg
	}
	m.madeInspections++
	return worryLevel
}

func (m *Monkey) getBored(wl int) int {
	fWorryLevel := float32(wl) / 3.0
	return int(math.Floor(float64(fWorryLevel)))
}

func (m *Monkey) test(wl int) bool {
	return wl%m.divisor == 0
}

func solve(monkeys []*Monkey, partB bool) int {
	part := "A"
	if partB {
		part = "B"
	}
	done := goaocd.Duration(fmt.Sprintf("Part %s", part))
	defer done()

	roundCount := 20
	if partB {
		roundCount = 10000
	}

	// I'm cheated =(
	// https://www.reddit.com/r/adventofcode/comments/zih7gf/2022_day_11_part_2_what_does_it_mean_find_another/
	supermodulo = 1
	for _, m := range monkeys {
		supermodulo *= m.divisor
	}

	for round := 0; round < roundCount; round++ {
		for i := 0; i < len(monkeys); i++ {
			monkey := monkeys[i]
			c := monkey.Turn(partB)
			for passTo := range c {
				if passTo == nil {
					break
				}
				monkeys[passTo.to].items = append(monkeys[passTo.to].items, passTo.wl)
			}
		}
		/*
			fmt.Print("Round completed\n")
			for idx, m := range monkeys {
				fmt.Printf("Monkey %d items: %v\n", idx, m.items)
			}
			fmt.Scanln()
		*/
	}

	var bl = make([]int, len(monkeys))
	for i, m := range monkeys {
		bl[i] = m.madeInspections
	}
	sort.Ints(bl)
	top2 := bl[len(bl)-2:]
	fmt.Printf("%v\n", top2)
	return top2[0] * top2[1]
}

func sample() []string {
	d := `
Monkey 0:
  Starting items: 79, 98
  Operation: new = old * 19
  Test: divisible by 23
    If true: throw to monkey 2
    If false: throw to monkey 3

Monkey 1:
  Starting items: 54, 65, 75, 74
  Operation: new = old + 6
  Test: divisible by 19
    If true: throw to monkey 2
    If false: throw to monkey 0

Monkey 2:
  Starting items: 79, 60, 97
  Operation: new = old * old
  Test: divisible by 13
    If true: throw to monkey 1
    If false: throw to monkey 3

Monkey 3:
  Starting items: 74
  Operation: new = old + 3
  Test: divisible by 17
    If true: throw to monkey 0
    If false: throw to monkey 1`

	d = strings.TrimLeft(d, "\n")
	return strings.Split(d, "\n")
}

func loadMonkeys() []*Monkey {
	monkeys := []*Monkey{}
	lines := goaocd.Lines()
	// lines := sample()
	var monkey *Monkey
	for idx, line := range lines {
		switch idx % 7 {
		case 0:
			monkey = &Monkey{}
			monkeys = append(monkeys, monkey)
		case 1:
			items := strings.Split(line[18:], ", ")
			monkey.items = make([]int, len(items))
			for i, v := range items {
				monkey.items[i] = util.Atoi(v)
			}
		case 2:
			opConfig := strings.Split(line[23:], " ")
			if opConfig[1] == "old" && opConfig[0] == "*" {
				monkey.op = "pow"
				monkey.arg = 2
			} else {
				monkey.op = opConfig[0]
				monkey.arg = util.Atoi(opConfig[1])
			}
		case 3:
			monkey.divisor = util.Atoi(line[21:])
		case 4:
			monkey.trueTo = util.Atoi(line[29:])
		case 5:
			monkey.falseTo = util.Atoi(line[30:])
		}
	}
	return monkeys
}

func main() {
	fmt.Printf("Part A: %d\n", solve(loadMonkeys(), false))
	fmt.Printf("Part B: %d\n", solve(loadMonkeys(), true))
}
