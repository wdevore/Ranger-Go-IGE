package api

// IQuadTree is a spatial container for querying
type IQuadTree interface {
	Add(INode) bool
	Remove(INode) bool

	Clean()
	Clear()

	SetBoundary(x, y, w, h float32)
	Boundary() IRectangle

	SetCapacity(capacity int)
	Capacity() int

	SetMaxDepth(depth int)
	MaxDepth() int
}
