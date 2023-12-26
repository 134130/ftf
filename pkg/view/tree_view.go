package view

import (
	"github.com/134130/ftf/pkg/config"
	term "github.com/134130/ftf/pkg/terminal"
	"github.com/134130/ftf/pkg/tree"
	"math"
	"strings"
)

type treeView struct {
	graphics config.Graphics
	state    *term.State
	rows     int
	scroll   int
	lineById map[string]int
}

func NewTreeView(graphics config.Graphics, state *term.State) term.ViewRenderer {
	return &treeView{
		graphics: graphics,
		state:    state,
	}
}

var _ term.ViewRenderer = (*treeView)(nil)

func (v treeView) Position(totalRows, totalCols int) term.Position {
	if v.state.Cursor.HasPreview() {
		return term.Position{
			Top:  1,
			Left: 1,
			Rows: totalRows - 1,
			Cols: int(math.Ceil(float64(totalCols) / 2.0)),
		}
	} else {
		return term.Position{
			Top:  1,
			Left: 1,
			Rows: totalRows - 1,
			Cols: totalCols,
		}
	}
}

func (v treeView) HasBorder() bool {
	return false
}

func (v treeView) ShouldRender() bool {
	return true
}

func (v treeView) Render(pos term.Position) []term.LineAppender {
	var lines []term.LineAppender
	v.rows = pos.Rows
	v.lineById = make(map[string]int)
	v.state.Root.Traverse(func(t tree.TreeHandler, depth int) error {
		line := v.renderNode(t, depth, pos.Cols)
		v.lineById[t.GetID()] = len(lines)
		lines = append(lines, line)
		return nil
	})
	v.scroll = v.scrollForId(v.state.Cursor.GetID())
	return lines[v.scroll:]
}

func (v treeView) scrollForId(id string) int {
	targetLine := v.lineById[id]
	if targetLine < v.scroll {
		return targetLine
	} else if targetLine >= v.scroll+v.rows {
		return targetLine - v.rows + 1
	} else {
		return v.scroll
	}
}

func (v treeView) Commands() map[string]term.Command {
	return map[string]term.Command{
		"tree:prev":   v.prev,
		"tree:next":   v.next,
		"tree:open":   v.open,
		"tree:close":  v.close,
		"tree:parent": v.parent,
	}
}

func (v treeView) prev(helper term.Helper, args ...interface{}) error {
	prev := v.state.Cursor.Prev()
	if prev != nil {
		v.state.Cursor = prev
	}
	return nil
}

func (v treeView) next(helper term.Helper, args ...interface{}) error {
	next := v.state.Cursor.Next()
	if next != nil {
		v.state.Cursor = next
	}
	return nil
}

func (v treeView) open(helper term.Helper, args ...interface{}) error {
	v.state.Cursor.Expand()
	return nil
}

func (v treeView) close(helper term.Helper, args ...interface{}) error {
	v.state.Cursor.Collapse()
	return nil
}

func (v treeView) parent(helper term.Helper, args ...interface{}) error {
	parent := v.state.Cursor.GetParent()
	if parent != nil {
		v.state.Cursor = parent
	}

	return nil
}

func (v treeView) renderNode(node tree.TreeHandler, indent, maxLength int) term.LineAppender {
	line := term.NewLine(maxLength, &term.Graphic{})
	line.Append(strings.Repeat("  ", indent), &term.Graphic{})

	graphic := term.Graphic{}
	if node == v.state.Cursor {
		if g, ok := v.graphics["tree:cursor"]; ok {
			graphic.Merge(g)
		}
	}

	if node.IsExpanded() {
		line.Append("▼ ", &graphic)
	} else {
		line.Append("▶ ", &graphic)
	}
	line.Append(node.GetName(), &graphic)
	return line
}