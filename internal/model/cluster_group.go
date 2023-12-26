package model

import (
	"github.com/134130/ftf/pkg/tree"
	"github.com/134130/ftf/pkg/util"
)

type ClusterGroup struct {
	UUID                    string
	Name                    string
	Preview                 string
	Parent                  tree.TreeHandler
	Children                []tree.TreeHandler
	expanded                bool
	highlightMatchedIndexes []int
}

var _ tree.TreeHandler = (*ClusterGroup)(nil)

func (c *ClusterGroup) GetID() string {
	return c.UUID
}

func (c *ClusterGroup) GetName() string {
	return c.Name
}

func (c *ClusterGroup) GetParent() tree.TreeHandler {
	return c.Parent
}

func (c *ClusterGroup) GetChildren() []tree.TreeHandler {
	return c.Children
}

func (c *ClusterGroup) GetChildrenByName(name string) []tree.TreeHandler {
	//TODO implement me
	panic("implement me")
}

func (c *ClusterGroup) GetHighlightMatchedIndexes() []int {
	return c.highlightMatchedIndexes
}

func (c *ClusterGroup) SetHighlightMatchedIndexes(ints []int) {
	c.highlightMatchedIndexes = ints
}

func (c *ClusterGroup) HasPreview() bool {
	return false
}

func (c *ClusterGroup) GetPreview() string {
	return c.Preview
}

func (c *ClusterGroup) Traverse(f func(tree.TreeHandler, int) error) error {
	return util.TraverseTree(c, f)
}

func (c *ClusterGroup) Prev() tree.TreeHandler {
	return util.PrevTree(c)
}

func (c *ClusterGroup) Next() tree.TreeHandler {
	return util.NextTree(c)
}

func (c *ClusterGroup) IsExpanded() bool {
	return c.expanded
}

func (c *ClusterGroup) Expand() {
	c.expanded = true
}

func (c *ClusterGroup) Collapse() {
	c.expanded = false
}
