package tree

type TreeHandler interface {
	GetID() string
	GetName() string
	GetParent() TreeHandler
	GetChildren() []TreeHandler
	GetChildrenByName(name string) []TreeHandler
	HasPreview() bool
	Traverse(func(TreeHandler, int) error) error
	IsExpanded() bool
	Expand()
	Collapse()
	Prev() TreeHandler
	Next() TreeHandler
}
