package main

import (
	"fmt"
	"strings"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
)

type Item struct {
	Value int64
	next  *Item
	prev  *Item

	nextOrder *Item
}

func sample() []string {
	d := `1
2
-3
3
-2
0
4`
	return strings.Split(d, "\n")
}

func fwd(from *Item, steps int64) *Item {
	cur := from
	var i int64
	for i = 0; i < steps; i++ {
		cur = cur.next
	}
	return cur
}

func rev(from *Item, steps int64) *Item {
	cur := from
	var i int64
	for i = 0; i >= steps; i-- {
		cur = cur.prev
	}

	return cur
}

func find(start *Item, v int64) *Item {
	cur := start
	for {
		if cur.Value == v {
			return cur
		}
		cur = cur.next
		if cur == nil {
			panic("End of list")
		}
		if cur == start {
			panic("Not found")
		}
	}
}

func display(start *Item) {
	initStart := start
	for {
		fmt.Printf("%d, ", start.Value)
		start = start.next
		if start == initStart {
			break
		}
	}
	fmt.Print("\n")
}

func mix(start *Item, l int) {
	cur := start

	for {
		v := cur.Value % int64(l-1)
		if v != 0 {
			var pt *Item
			if v > 0 {
				pt = fwd(cur, v)
			} else if cur.Value < 0 {
				pt = rev(cur, v)
			}
			if pt == cur || pt == cur.prev {
				panic("Cycle!")
			} else {
				// oldPrev [cur] oldNext
				// nil [cur] oldNext
				// oldPrev [cur] nil
				oldPrev := cur.prev
				oldNext := cur.next

				oldPrev.next = cur.next
				oldNext.prev = cur.prev

				// pt [o] pt.next
				ptNext := pt.next

				// insert to the middle (start/end do not change)
				pt.next = cur
				cur.next = ptNext
				ptNext.prev = cur
				cur.prev = pt
			}
		}

		cur = cur.nextOrder
		if cur == nil {
			break
		}
	}
}

func calcGrove(start *Item) int64 {
	fmt.Printf("Calculating result\n")
	z := find(start, 0)
	fmt.Printf("Zero found\n")
	n1000 := fwd(z, 1000)
	fmt.Printf("1k found %d\n", n1000.Value)
	n2000 := fwd(z, 2000)
	fmt.Printf("2k found %d\n", n2000.Value)
	n3000 := fwd(z, 3000)
	fmt.Printf("3k found %d\n", n3000.Value)

	return n1000.Value + n2000.Value + n3000.Value
}

func partA(lines []string) int64 {
	done := goaocd.Duration("Part A")
	defer done()

	start, l := mkList(lines)
	mix(start, l)
	return calcGrove(start)
}

func forEach(start *Item, cb func(i *Item)) {
	cur := start
	for {
		cb(cur)
		cur = cur.next
		if cur == start {
			break
		}
	}
}

func partB(lines []string) int64 {
	done := goaocd.Duration("Part B")
	defer done()

	start, l := mkList(lines)

	forEach(start, func(i *Item) {
		i.Value *= 811589153
	})

	for i := 0; i < 10; i++ {
		mix(start, l)
	}

	return calcGrove(start)
}

func mkList(lines []string) (*Item, int) {
	var start, cur *Item
	for _, val := range lines {
		v := int64(util.Atoi(val))
		if start == nil {
			start = &Item{Value: v}
			cur = start
		} else {
			i := &Item{Value: v}
			cur.next = i
			cur.nextOrder = i
			i.prev = cur
			cur = i
		}
	}
	cur.nextOrder = nil
	cur.next = start
	start.prev = cur
	return start, len(lines)
}

func main() {
	lines := sample()
	lines = goaocd.Lines()

	fmt.Printf("Part A: %d\n", partA(lines))
	fmt.Printf("Part B: %d\n", partB(lines))
}
