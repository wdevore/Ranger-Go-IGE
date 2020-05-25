package api

// IPolygon represents 2D lines
type IPolygon interface {
	// AddVertex appends the point to vertices
	AddVertex(x, y float32)

	SetVertex(x, y float32, index int)

	PointInside(p IPoint) bool
}
