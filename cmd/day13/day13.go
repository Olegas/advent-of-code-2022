package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
)

// 1 - ok
// 0 - failed
// -1 - proceed
func Compare(l, r []any) int {
	for idx, i := range l {
		if idx > len(r)-1 {
			return 0
		}

		j := r[idx]
		switch i.(type) {
		case int:
			switch j := j.(type) {
			case int:
				if i.(int) < j {
					return 1
				} else if i.(int) > j {
					return 0
				}
			case []interface{}:
				cmp := make([]any, 1)
				cmp[0] = i
				res := Compare(cmp, j)
				if res > -1 {
					return res
				}
			}
		case []interface{}:
			switch j.(type) {
			case int:
				cmp := make([]any, 1)
				cmp[0] = j
				res := Compare(i.([]any), cmp)
				if res > -1 {
					return res
				}
			case []interface{}:
				res := Compare(i.([]any), j.([]any))
				if res > -1 {
					return res
				}
			}
		}
	}
	if len(l) < len(r) {
		return 1
	} else if len(l) == len(r) {
		return -1
	}
	return 0
}

func parse(r *strings.Reader, current *[]any) {
	for {
		b, err := r.ReadByte()
		if err != nil {
			panic(err)
		}
		if util.IsDigit(b) {
			accu := string(b)
			for {
				b, err := r.ReadByte()
				if err != nil {
					panic(err)
				}
				if util.IsDigit(b) {
					accu += string(b)
				} else {
					*current = append(*current, util.Atoi(accu))
					r.Seek(-1, io.SeekCurrent)
					break
				}
			}
		} else if b == ']' {
			return
		} else if b == '[' {
			newItem := make([]any, 0)
			parse(r, &newItem)
			*current = append(*current, newItem)
		}
	}
}

func printInner(packet []any) {
	for idx, i := range packet {
		if idx != 0 {
			fmt.Print(",")
		}
		switch i := i.(type) {
		case int:
			fmt.Printf("%d", i)
		case []interface{}:
			fmt.Printf("[")
			printInner(i)
			fmt.Printf("]")
		}

	}
}

func print(packets [][]any) {
	for _, p := range packets {
		fmt.Print("[")
		printInner(p)
		fmt.Print("]\n")
	}
}

func main() {
	data := goaocd.Lines()
	lines := append(data, "\n")
	var pair = make([][]any, 2)
	var accu = 0
	var pairNum = 1
	packets := make([][]any, 0)
	for idx, line := range lines {
		if idx%3 == 2 {
			// lastLint
			if Compare(pair[0], pair[1]) == 1 {
				accu += pairNum
			}
			// next pair
			pairNum++
		} else {
			newItem := make([]any, 0)
			r := strings.NewReader(line)
			// Skip first char
			r.Seek(1, io.SeekStart)
			parse(r, &newItem)
			pair[idx%3] = newItem
			packets = append(packets, newItem)
		}
	}
	fmt.Printf("Part A: %d\n", accu)

	// Some typing shit ( Shame on me
	pkt1 := make([]any, 1)
	i := make([]any, 1)
	i[0] = 2
	pkt1[0] = i
	pkt2 := make([]any, 1)
	i = make([]any, 1)
	i[0] = 6
	pkt2[0] = i
	packets = append(packets, nil, nil)
	packets[len(packets)-2] = pkt1
	packets[len(packets)-1] = pkt2

	sort.SliceStable(packets, func(i, j int) bool {
		a := packets[i]
		b := packets[j]
		return Compare(a, b) == 1
	})

	accu = 1
	for idx, p := range packets {
		if len(p) == 1 {
			i := p[0]
			switch i := i.(type) {
			case []interface{}:
				if len(i) == 1 {
					j := i[0]
					if j == 2 || j == 6 {
						accu *= idx + 1
					}
				}
			}
		}
	}
	// print(packets)

	fmt.Printf("Part B: %d\n", accu)
}
