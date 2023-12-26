package terminal

import "github.com/134130/ftf/pkg/tree"

type State struct {
	Root         tree.TreeHandler
	Cursor       tree.TreeHandler
	Selection    []tree.TreeHandler
	SearchString string
	Rerender     chan bool
}
