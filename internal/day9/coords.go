package day9

import "fmt"

type Coords struct {
	X, Y int
}

func (c *Coords) ToString() string {
	return fmt.Sprintf("%d,%d", c.X, c.Y)
}

func (c *Coords) Up() {
	c.Y--
}

func (c *Coords) Down() {
	c.Y++
}

func (c *Coords) Left() {
	c.X--
}

func (c *Coords) Right() {
	c.X++
}
