package view

import (
	"github.com/134130/ftf/pkg/config"
	term "github.com/134130/ftf/pkg/terminal"
	"github.com/134130/ftf/pkg/tree"
	"github.com/fatih/color"
	"github.com/sahilm/fuzzy"
	"strings"
)

type searchbarView struct {
	graphics config.Graphics
	state    *term.State
}

func NewSearchbarView(graphics config.Graphics, state *term.State) term.ViewRenderer {
	return &searchbarView{
		graphics: graphics,
		state:    state,
	}
}

var _ term.ViewRenderer = (*searchbarView)(nil)

func (v *searchbarView) Position(totalRows int, totalCols int) term.Position {
	return term.Position{
		Top:  1,
		Left: 1,
		Rows: 2,
		Cols: totalCols,
	}
}

func (v *searchbarView) Commands() map[string]term.Command {
	return map[string]term.Command{
		"searchbar:backspace": func(helper term.Helper, args ...interface{}) error {
			if len(v.state.SearchString) > 0 {
				v.state.SearchString = v.state.SearchString[:len(v.state.SearchString)-1]
			}
			v.search()
			return nil
		},
		"searchbar:append": func(helper term.Helper, args ...interface{}) error {
			v.state.SearchString += string(args[0].(rune))
			v.search()
			return nil
		},
	}
}

type treeToDepths []treeToDepth
type treeToDepth struct {
	tree  tree.TreeHandler
	depth int
}

func (t treeToDepths) String(i int) string {
	return t[i].tree.GetName()
}

func (t treeToDepths) Len() int {
	return len(t)
}

func (v *searchbarView) search() {
	searchString := v.state.SearchString

	// Get all nodes in the tree
	var nodes treeToDepths
	stack := []treeToDepth{{v.state.Root, 0}}
	for len(stack) > 0 {
		// Pop node from stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Add children to stack
		children := current.tree.GetChildren()
		for i := len(children) - 1; i >= 0; i-- {
			stack = append(stack, treeToDepth{children[i], current.depth + 1})
		}

		current.tree.SetHighlightMatchedIndexes(nil)
		current.tree.Collapse()
		nodes = append(nodes, current)
	}

	// Filter nodes that match the search string
	matches := fuzzy.FindFrom(searchString, nodes)
	matchNodes := make([]tree.TreeHandler, len(matches))
	for i, match := range matches {
		matchNode := nodes[match.Index].tree
		matchNode.SetHighlightMatchedIndexes(match.MatchedIndexes)
		matchNodes[i] = matchNode
	}

	// Expand all matched nodes and it's parent that match the search string
	for _, node := range matchNodes {
		node.Expand()
		for parent := node.GetParent(); parent != nil; parent = parent.GetParent() {
			parent.Expand()
		}
	}
}

func (v *searchbarView) HasBorder() bool {
	return false
}

func (v *searchbarView) ShouldRender() bool {
	return true
}

func (v *searchbarView) Render(position term.Position) []term.LineAppender {
	searchLine := term.NewLine(position.Cols, &term.Graphic{}).
		AppendRaw(color.BlueString("> ")).
		Append(v.state.SearchString, &term.Graphic{})

	separatorLine := term.NewLine(position.Cols, &term.Graphic{}).
		AppendRaw(color.HiBlackString(strings.Repeat("─", 50)))

	return []term.LineAppender{searchLine, separatorLine}
}
