package view

import (
	"fmt"
	"ftf/internal/config"
	term "ftf/pkg/terminal"
	"math"
	"os/exec"
	"strings"
)

type previewView struct {
	graphics    config.Graphics
	state       *term.State
	lastId      string
	lastPreview []string
	scroll      int
	noLines     int
}

func NewPreviewView(graphics config.Graphics, state *term.State) term.ViewRenderer {
	return &previewView{
		graphics: graphics,
		state:    state,
	}
}

var _ term.ViewRenderer = (*previewView)(nil)

func (v *previewView) Position(totalRows, totalCols int) term.Position {
	return term.Position{
		Top:  1,
		Left: int(math.Ceil(float64(totalCols)/2.0)) + 1,
		Rows: totalRows - 1,
		Cols: int(math.Floor(float64(totalCols) / 2.0)),
	}
}

func (v *previewView) HasBorder() bool {
	return true
}

func (v *previewView) ShouldRender() bool {
	return v.state.Cursor.HasPreview()
}

func (v *previewView) Render(p term.Position) []term.LineAppender {
	if v.lastId != v.state.Cursor.GetID() {
		v.lastId = v.state.Cursor.GetID()
		v.scroll = 0

		preview, err := getPreview("cat {}", v.state.Cursor.GetID())
		if err != nil {
			preview = err.Error()
		}
		preview = strings.ReplaceAll(preview, "\t", "    ")
		v.lastPreview = strings.Split(preview, "\n")
	}

	lines := v.lastPreview
	if v.scroll > len(lines)-p.Rows {
		if len(lines) < p.Rows {
			v.scroll = 0
		} else {
			v.scroll = len(lines) - p.Rows
		}
	}

	var termLines []term.LineAppender
	for i := v.scroll; i-v.scroll < p.Rows && i < len(lines); i++ {
		termLine := term.NewLine(p.Cols, &term.Graphic{})
		termLine.AppendRaw(lines[i])
		termLines = append(termLines, termLine)
	}
	return termLines
}

func getPreview(cmdTemplate string, path string) (string, error) {
	//return "preview", nil
	var stdout, stderr strings.Builder
	preview := exec.Command("bash", "-c", "cat /Users/cooper/.ssh/chequer-dev.pem")
	preview.Stdout = &stdout
	preview.Stderr = &stderr
	err := preview.Run()
	if err != nil {
		return stdout.String(), fmt.Errorf("%w %s", err, stderr.String())
	} else {
		return stdout.String(), nil
	}
}

func (v *previewView) Commands() map[string]term.Command {
	return map[string]term.Command{
		"preview:down": v.down,
		"preview:up":   v.up,
	}
}

func (v *previewView) up(helper term.Helper, args ...interface{}) error {
	if v.scroll > 0 {
		v.scroll -= 1
	}
	return nil
}

func (v *previewView) down(helper term.Helper, args ...interface{}) error {
	v.scroll += 1
	return nil
}
