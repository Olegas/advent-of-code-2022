package day9

import "fmt"

type Rope struct {
	Head *Knot
}

func NewRope(len int) *Rope {
	rope := &Rope{Head: &Knot{Id: "H"}}
	cur := rope.Head
	for i := 1; i < len; i++ {
		knot := &Knot{Id: fmt.Sprintf("%d", i)}
		cur.Next = knot
		cur = knot
	}
	return rope
}

func (r *Rope) MoveHead(direction string) {
	switch direction {
	case "U":
		r.Head.Pos.Up()
	case "D":
		r.Head.Pos.Down()
	case "L":
		r.Head.Pos.Left()
	case "R":
		r.Head.Pos.Right()
	}
}

func (r *Rope) Tail() *Knot {
	var k *Knot
	for k = r.Head; k.Next != nil; k = k.Next {
	}
	return k
}

func (r *Rope) Adjust() {
	for k := r.Head; k != nil; k = k.Next {
		if k.AdjustNext() {
			break
		}
	}
}
