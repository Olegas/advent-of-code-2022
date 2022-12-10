package day10

import "fmt"

type CRT struct {
	cycle int
	line  string
}

func (c *CRT) pos() int {
	return c.cycle % 40
}

func (c *CRT) Inc() {
	c.cycle++
	if c.pos() == 0 {
		fmt.Printf("%s\n", c.line)
		c.line = ""
	}
}

func (c *CRT) Draw(spritePos int) {
	pos := c.pos()
	char := "."
	if spritePos-1 <= pos && pos <= spritePos+1 {
		char = "#"
	}
	c.line = fmt.Sprintf("%s%s", c.line, char)
}
