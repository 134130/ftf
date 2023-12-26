package util

import (
	"github.com/134130/ftf/pkg/tree"
)

func TraverseTree(t tree.TreeHandler, f func(tree.TreeHandler, int) error) error {
	type treeWithDepth struct {
		tree  tree.TreeHandler
		depth int
	}

	stack := []treeWithDepth{{t, 0}}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		err := f(current.tree, current.depth)
		if err != nil {
			return err
		}

		if current.tree.IsExpanded() {
			children := current.tree.GetChildren()
			for i := len(children) - 1; i >= 0; i-- {
				stack = append(stack, treeWithDepth{children[i], current.depth + 1})
			}
		}
	}
	return nil
}

func PrevTree(t tree.TreeHandler) tree.TreeHandler {
	parent := t.GetParent()
	if parent == nil {
		return nil
	}

	siblings := parent.GetChildren()
	if t == siblings[0] {
		return parent
	}

	var prevSibling tree.TreeHandler
	for i, sibling := range siblings {
		if sibling == t {
			prevSibling = siblings[i-1]
		}
	}

	node := prevSibling
	for {
		if !node.IsExpanded() {
			return node
		}
		children := node.GetChildren()
		if len(children) == 0 {
			return node
		} else {
			node = children[len(children)-1]
		}
	}
}

func NextTree(t tree.TreeHandler) tree.TreeHandler {
	children := t.GetChildren()
	if t.IsExpanded() && len(children) > 0 {
		return children[0]
	}

	var node = t
	for node.GetParent() != nil {
		siblings := node.GetParent().GetChildren()
		for i, sibling := range siblings {
			if sibling == node && i < len(siblings)-1 {
				return siblings[i+1]
			}
		}
		node = node.GetParent()
	}
	return nil
}
