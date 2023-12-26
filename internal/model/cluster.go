package model

import (
	"github.com/134130/ftf/pkg/tree"
	"github.com/134130/ftf/pkg/util"
)

type Cluster struct {
	UUID                    string
	Name                    string
	Preview                 string
	Parent                  tree.TreeHandler
	Children                []tree.TreeHandler
	expanded                bool
	highlightMatchedIndexes []int
}

var _ tree.TreeHandler = (*Cluster)(nil)

func (c *Cluster) GetID() string {
	return c.UUID
}

func (c *Cluster) GetName() string {
	return c.Name
}

func (c *Cluster) GetParent() tree.TreeHandler {
	return c.Parent
}

func (c *Cluster) GetChildren() []tree.TreeHandler {
	return nil
}

func (c *Cluster) GetChildrenByName(name string) []tree.TreeHandler {
	//TODO implement me
	panic("implement me")
}

func (c *Cluster) GetHighlightMatchedIndexes() []int {
	return c.highlightMatchedIndexes
}

func (c *Cluster) SetHighlightMatchedIndexes(ints []int) {
	c.highlightMatchedIndexes = ints
}

func (c *Cluster) HasPreview() bool {
	return true
}

func (c *Cluster) GetPreview() string {
	return c.Preview
}

func (c *Cluster) Traverse(f func(tree.TreeHandler, int) error) error {
	return util.TraverseTree(c, f)
}

func (c *Cluster) Prev() tree.TreeHandler {
	return util.PrevTree(c)
}

func (c *Cluster) Next() tree.TreeHandler {
	return util.NextTree(c)
}

func (c *Cluster) IsExpanded() bool {
	return c.expanded
}

func (c *Cluster) Expand() {
	c.expanded = true
}

func (c *Cluster) Collapse() {
	c.expanded = false
}
