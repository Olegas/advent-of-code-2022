package day10

import "strconv"

type Machine struct {
	channel   chan *Signal
	cycles    int
	registers map[string]int
}

var cyclesPerCommand = map[string]int{"noop": 1, "addx": 2}

func NewMachine(c chan *Signal) *Machine {
	reg := map[string]int{"x": 1}
	m := Machine{channel: c, registers: reg}
	return &m
}

func (m *Machine) Execute(op string, args []string) {
	switch op {
	case "noop":
		m.waitCommandCompleted(op, func() {})
	case "addx":
		v, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}
		m.waitCommandCompleted(op, func() { m.registers["x"] += v })
	}

}

func (m *Machine) notify() {
	sig := m.cycles * m.registers["x"]
	m.channel <- &Signal{m.cycles, sig}
}

func (m *Machine) waitCommandCompleted(op string, done func()) {
	for i := 0; i < cyclesPerCommand[op]; i++ {
		m.cycles++
		m.notify()
	}
	done()
}
