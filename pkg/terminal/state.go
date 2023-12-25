package terminal

import "ftf/pkg/tree"

type State struct {
	Root   tree.TreeHandler
	Cursor tree.TreeHandler
}
