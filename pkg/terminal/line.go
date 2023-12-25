package terminal

import (
	"golang.org/x/text/width"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var escapeRegex, graphicsEscapeRegex *regexp.Regexp

func init() {
	escapeRegex = regexp.MustCompile("\x1b\\[[0-9;]*[a-zA-Z]")
	graphicsEscapeRegex = regexp.MustCompile("\x1b\\[[0-9;]*m")
}

type LineAppender interface {
	Append(string, *Graphic) LineAppender
	AppendRaw(string) LineAppender
	Length() int
	Text() string
}

type line struct {
	line            strings.Builder
	length          int
	maxLength       int
	defaultGraphics *Graphic
}

var _ LineAppender = &line{}

func NewLine(maxLength int, defaultGraphics *Graphic) LineAppender {
	return &line{
		maxLength:       maxLength,
		defaultGraphics: defaultGraphics,
	}
}

func (l *line) Append(s string, graphics *Graphic) LineAppender {
	if l.length >= l.maxLength {
		return l
	}

	if graphics != nil {
		l.line.WriteString(graphics.ToEscapeCode())
	}
	l.appendText(s)
	l.line.WriteString(resetGraphics)
	l.line.WriteString(l.defaultGraphics.ToEscapeCode())
	return l
}

func (l *line) AppendRaw(s string) LineAppender {
	matches := escapeRegex.FindAllStringIndex(s, -1)
	prevIndex := 0
	for i := 0; i < len(matches)+1; i++ {
		piece := ""
		if i < len(matches) {
			piece = s[prevIndex:matches[i][0]]
		} else {
			piece = s[prevIndex:]
		}
		l.appendText(piece)
		if i < len(matches) {
			escapeCode := s[matches[i][0]:matches[i][1]]
			if graphicsEscapeRegex.MatchString(escapeCode) {
				l.line.WriteString(escapeCode)
			}
			prevIndex = matches[i][1]
		}
	}
	l.line.WriteString(resetGraphics)
	l.line.WriteString(l.defaultGraphics.ToEscapeCode())
	return l
}

func (l *line) Length() int {
	return l.length
}

func (l *line) Text() string {
	return l.line.String()
}

func (l *line) appendText(s string) {
	for len(s) > 0 && l.length < l.maxLength {
		r, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		if r == utf8.RuneError || unicode.IsMark(r) || unicode.IsControl(r) {
			continue
		}
		termWidth := 1
		runeKind := width.LookupRune(r).Kind()
		if runeKind == width.EastAsianWide || runeKind == width.EastAsianFullwidth {
			termWidth = 2
		}
		if l.length+termWidth <= l.maxLength {
			l.length += termWidth
		} else {
			break
		}
		l.line.WriteString(string(r))
	}
}
