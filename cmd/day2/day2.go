package main

import (
	"fmt"
	"strings"

	"github.com/Olegas/goaocd"
)

var scoresMap = map[string]int{"A": 1, "B": 2, "C": 3}
var outcomeMap = map[string]int{"X": 0, "Y": 3, "Z": 6}

func scoreByPick(you string) int {
	return scoresMap[you]
}

func scoreBuOutcome(outcome string) int {
	return outcomeMap[outcome]
}

func pickMove(opponent, outcome string) string {
	if outcome == "Y" {
		return opponent
	}
	if opponent == "A" { // Rock
		switch outcome {
		case "X":
			// Lose
			return "C"
		case "Z":
			// Win
			return "B"
		}
	} else if opponent == "B" { // Paper
		switch outcome {
		case "X":
			// Lose
			return "A"
		case "Z":
			// Win
			return "C"
		}
	} else if opponent == "C" { // Scissors
		switch outcome {
		// Lose
		case "X":
			return "B"
		// Win
		case "Z":
			return "A"
		}
	}
	panic(fmt.Sprintf("Impossible! %s %s", opponent, outcome))
}

func score(you, outcome string) int {
	return scoreBuOutcome(outcome) + scoreByPick(you)
}

func main() {
	accu := 0
	for _, line := range goaocd.Lines(2) {
		steps := strings.Split(line, " ")
		if len(steps) != 2 {
			continue
		}
		opponent := steps[0]
		outcome := steps[1]
		you := pickMove(opponent, outcome)
		accu += score(you, outcome)
	}
	fmt.Printf("%d", accu)
}
