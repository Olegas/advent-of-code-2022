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
		switch i := i.(type) {
		case int:
			switch j := j.(type) {
			case int:
				if i < j {
					return 1
				} else if i > j {
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
			switch j := j.(type) {
			case int:
				cmp := make([]any, 1)
				cmp[0] = j
				res := Compare(i, cmp)
				if res > -1 {
					return res
				}
			case []interface{}:
				res := Compare(i, j)
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

func parse(r *strings.Reader, current *[]any) []any {
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
			return *current
		} else if b == '[' {
			if current == nil {
				ret := make([]any, 0)
				current = &ret
				continue
			}
			newItem := make([]any, 0)
			parse(r, &newItem)
			*current = append(*current, newItem)
		}
	}
}

func partA(lines []string, packets *[][]any) int {
	done := goaocd.Duration("Part A")
	defer done()

	var accu = 0
	var pairNum = 1
	var pair = make([][]any, 2)

	for idx, line := range lines {
		if idx%3 == 2 {
			// last line
			if Compare(pair[0], pair[1]) == 1 {
				accu += pairNum
			}
			// next pair
			pairNum++
		} else {
			newItem := parse(strings.NewReader(line), nil)
			pair[idx%3] = newItem
			*packets = append(*packets, newItem)
		}
	}
	return accu
}

func partB(packets *[][]any) int {
	done := goaocd.Duration("Part B")
	defer done()

	// Some typing shit ( Shame on me
	pkt1 := make([]any, 1)
	i := make([]any, 1)
	i[0] = 2
	pkt1[0] = i
	pkt2 := make([]any, 1)
	i = make([]any, 1)
	i[0] = 6
	pkt2[0] = i
	*packets = append(*packets, nil, nil)
	(*packets)[len(*packets)-2] = pkt1
	(*packets)[len(*packets)-1] = pkt2

	sort.SliceStable(*packets, func(i, j int) bool {
		a := (*packets)[i]
		b := (*packets)[j]
		return Compare(a, b) == 1
	})

	var accu = 1
	for idx, p := range *packets {
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

	return accu
}

func main() {
	data := goaocd.Lines()
	lines := append(data, "\n")

	packets := make([][]any, 0)
	fmt.Printf("Part A: %d\n", partA(lines, &packets))
	fmt.Printf("Part B: %d\n", partB(&packets))
}
