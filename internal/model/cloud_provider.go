package model

import (
	"github.com/134130/ftf/pkg/tree"
	"github.com/134130/ftf/pkg/util"
)

type CloudProvider struct {
	UUID                    string
	Name                    string
	Children                []tree.TreeHandler
	expanded                bool
	highlightMatchedIndexes []int
}

var _ tree.TreeHandler = (*CloudProvider)(nil)

func (c *CloudProvider) GetID() string {
	return c.UUID
}

func (c *CloudProvider) GetName() string {
	return c.Name
}

func (c *CloudProvider) GetParent() tree.TreeHandler {
	return nil
}

func (c *CloudProvider) GetChildren() []tree.TreeHandler {
	return c.Children
}

func (c *CloudProvider) GetChildrenByName(name string) []tree.TreeHandler {
	//TODO implement me
	panic("implement me")
}

func (c *CloudProvider) GetHighlightMatchedIndexes() []int {
	return c.highlightMatchedIndexes
}

func (c *CloudProvider) SetHighlightMatchedIndexes(ints []int) {
	c.highlightMatchedIndexes = ints
}

func (c *CloudProvider) HasPreview() bool {
	return false
}

func (c *CloudProvider) Traverse(f func(tree.TreeHandler, int) error) error {
	return util.TraverseTree(c, f)
}

func (c *CloudProvider) Prev() tree.TreeHandler {
	return util.PrevTree(c)
}

func (c *CloudProvider) Next() tree.TreeHandler {
	return util.NextTree(c)
}

func (c *CloudProvider) IsExpanded() bool {
	return c.expanded
}

func (c *CloudProvider) Expand() {
	c.expanded = true
}

func (c *CloudProvider) Collapse() {
	c.expanded = false
}
