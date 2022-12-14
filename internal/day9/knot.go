package day9

import "github.com/Olegas/advent-of-code-2022/internal/util"

type Knot struct {
	Id   string
	Pos  Coords
	Next *Knot
}

func (k *Knot) nextIsNear() bool {
	if k.Next == nil {
		return true
	}
	xDiff := util.AbsDiffInt(k.Pos.X, k.Next.Pos.X)
	yDiff := util.AbsDiffInt(k.Pos.Y, k.Next.Pos.Y)
	return xDiff <= 1 && yDiff <= 1
}

func (k *Knot) AdjustNext() bool {
	if k.nextIsNear() {
		return true
	}
	next := k.Next
	if k.Pos.Y > next.Pos.Y {
		next.Pos.Down()
	} else if k.Pos.Y < next.Pos.Y {
		next.Pos.Up()
	}
	if k.Pos.X > next.Pos.X {
		next.Pos.Right()
	} else if k.Pos.X < next.Pos.X {
		next.Pos.Left()
	}
	return false
}
