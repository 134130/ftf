package terminal

import "fmt"

type ColorCoder interface {
	FgCode() string
	BgCode() string
}

type Color3Bit struct {
	Value  int
	Bright bool
}

type Color8Bit struct {
	Value int
}

type Color24Bit struct {
	R, G, B int
}

func (c Color3Bit) FgCode() string {
	if c.Bright {
		return fmt.Sprint("9", c.Value)
	} else {
		return fmt.Sprint("3", c.Value)
	}
}

func (c Color3Bit) BgCode() string {
	if c.Bright {
		return fmt.Sprint("10", c.Value)
	} else {
		return fmt.Sprint("4", c.Value)
	}
}

func (c Color8Bit) FgCode() string {
	return fmt.Sprint("38;5;", c.Value)
}

func (c Color8Bit) BgCode() string {
	return fmt.Sprint("48;5;", c.Value)
}

func (c Color24Bit) FgCode() string {
	return fmt.Sprint("38;2;", c.R, ";", c.G, ";", c.B)
}

func (c Color24Bit) BgCode() string {
	return fmt.Sprint("48;2;", c.R, ";", c.G, ";", c.B)
}
