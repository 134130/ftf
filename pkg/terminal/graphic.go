package terminal

import "strings"

type Graphic struct {
	FgColor ColorCoder
	BgColor ColorCoder

	Bold    bool
	Reverse bool
}

func (g *Graphic) ToEscapeCode() string {
	codes := []string{}
	if g.Bold {
		codes = append(codes, bold)
	}
	if g.Reverse {
		codes = append(codes, reverse)
	}
	if g.FgColor != nil {
		codes = append(codes, g.FgColor.FgCode())
	}
	if g.BgColor != nil {
		codes = append(codes, g.BgColor.BgCode())
	}

	return csi + strings.Join(codes, ";") + "m"
}

func (g *Graphic) Merge(other *Graphic) {
	if g.FgColor == nil {
		g.FgColor = other.FgColor
	}
	if g.BgColor == nil {
		g.BgColor = other.BgColor
	}
	g.Bold = g.Bold || other.Bold
	g.Reverse = g.Reverse || other.Reverse
}
