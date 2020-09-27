package api

// QuadrantBoundsFunc is a functor that returns the current rectangle
// during a traversal.
type QuadrantBoundsFunc func(bounds IRectangle)

// IQuadTree is a spatial container for querying
type IQuadTree interface {
	Add(INode) bool
	Remove(INode) bool
	Query(IRectangle, *[]INode)
	Traverse(quadrantCB QuadrantBoundsFunc)

	Clean()
	Clear()

	SetBoundary(x, y, w, h float32)
	SetBoundaryByMinMax(minX, minY, maxX, maxY float32)
	Boundary() IRectangle

	SetCapacity(capacity int)
	Capacity() int

	SetMaxDepth(depth int)
	MaxDepth() int
}
