package tree

type TreeHandler interface {
	GetID() string
	GetName() string
	GetParent() TreeHandler
	GetChildren() []TreeHandler
	GetChildrenByName(name string) []TreeHandler
	GetHighlightMatchedIndexes() []int
	SetHighlightMatchedIndexes([]int)
	HasPreview() bool
	GetPreview() string
	Traverse(func(t TreeHandler, depth int) error) error
	IsExpanded() bool
	Expand()
	Collapse()
	Prev() TreeHandler
	Next() TreeHandler
}
