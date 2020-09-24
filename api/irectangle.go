package api

// IRectangle represents a 2D rectangle
type IRectangle interface {
	Set(x, y, w, h float32)
	SetMinMax(minX, minY, maxX, maxY float32)
	SetSize(w, h float32)
	SetCenter(x, y float32)

	Center() IPoint

	Width() float32
	Height() float32

	Left() float32
	Right() float32
	Top() float32
	Bottom() float32

	PointContained(p IPoint) bool
	PointInside(p IPoint) bool
	Intersects(o IRectangle) bool
	Contains(o IRectangle) bool
}
