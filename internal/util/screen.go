package util

import "fmt"

func ClearTerm() {
	fmt.Print("\033[H\033[2J")
}
