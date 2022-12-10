package main

import (
	"fmt"
	"strings"

	"github.com/Olegas/advent-of-code-2022/internal/day10"
	"github.com/Olegas/goaocd"
)

func readSignal(signal chan *day10.Signal, commands *[]string) {
	m := day10.NewMachine(signal)
	for _, line := range *commands {
		a := strings.Split(line, " ")
		op := a[0]
		args := a[1:]
		m.Execute(op, args)
	}
	signal <- nil
}

func partA(commands *[]string) int {
	done := goaocd.Duration("Part A")
	defer done()

	c := make(chan *day10.Signal)
	go readSignal(c, commands)

	accu := 0
	for signal := range c {
		if signal == nil {
			break
		}
		if signal.Cycle == 20 || (signal.Cycle-20)%40 == 0 {
			fmt.Printf("Cycle %d, Register: %d, Value %d\n", signal.Cycle, signal.Value/signal.Cycle, signal.Value)
			accu += signal.Value
		}
	}

	return accu
}

func partB(commands *[]string) {
	done := goaocd.Duration("Part B")
	defer done()

	c := make(chan *day10.Signal)
	go readSignal(c, commands)

	crt := day10.CRT{}

	for signal := range c {
		if signal == nil {
			break
		}
		crt.Draw(signal.Value / signal.Cycle)
		crt.Inc()
	}
}

func main() {
	lines := goaocd.Lines()
	a := partA(&lines)
	fmt.Printf("Part A: %d\n", a)

	partB(&lines)
}
