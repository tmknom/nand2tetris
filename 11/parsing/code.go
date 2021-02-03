package parsing

import "fmt"

const DebugCode = true

type Code struct {
	Lines      []string
	DebugLines []string
}

func NewCode() *Code {
	return &Code{
		Lines: []string{},
	}
}

func (c *Code) AddCode(subroutineDec *SubroutineDec) {
	lines := subroutineDec.ToCode()
	c.Lines = append(c.Lines, lines...)
	c.Lines = append(c.Lines, "")
	c.addDebugCode(lines)
}

func (c *Code) addDebugCode(lines []string) {
	c.DebugLines = append(c.DebugLines, "")
	c.DebugLines = append(c.DebugLines, "==================")
	c.DebugLines = append(c.DebugLines, lines...)
}

func (c *Code) PrintDebugCode() {
	if DebugCode {
		for _, line := range c.DebugLines {
			fmt.Println(line)
		}
	}
}

func (c *Code) CodeLines() []string {
	return c.Lines
}
