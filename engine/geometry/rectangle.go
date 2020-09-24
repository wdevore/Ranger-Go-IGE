package geometry

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// Ranger's device coordinate space is oriented the same as OpenGL
// ^ +Y
// |
// |
// |
// .--------> +X

// Rectangle is a square with behaviours.
//
//                    >           maxX, maxY
//          .--------------------.
//          |        Top         |
//          |                    |
//          |                    |
//       >  | Left    .   right  |  <
//          |                    |
//          |                    |
//          |      bottom        |
//          .--------------------.
//     minX,minY      <
//
//
type Rectangle struct {
	left, top, bottom, right float32
	width, height            float32
	center                   api.IPoint
}

// NewRectangle creates a circle
func NewRectangle() api.IRectangle {
	o := new(Rectangle)
	o.center = NewPoint()
	return o
}

// Set sets top/left and width/height
func (r *Rectangle) Set(x, y, w, h float32) {
	r.left = x
	r.bottom = y
	r.top = r.bottom + h
	r.right = r.left + w

	r.width = w
	r.height = h
}

// SetMinMax sets the top/left and bottom/right corners
func (r *Rectangle) SetMinMax(minX, minY, maxX, maxY float32) {
	r.left = minX
	r.right = maxX
	r.top = maxY
	r.bottom = minY

	r.width = maxX - minX
	r.height = maxY - minY
}

// SetSize sets just the width/height
func (r *Rectangle) SetSize(w, h float32) {
	r.width = w
	r.height = h
}

// SetCenter of circle
func (r *Rectangle) SetCenter(x, y float32) {
	r.center.SetByComp(x, y)
}

// Center of circle
func (r *Rectangle) Center() api.IPoint {
	return r.center
}

// Width returns right-left length
func (r *Rectangle) Width() float32 {
	return r.right - r.left
}

// Height returns top-bottom length
func (r *Rectangle) Height() float32 {
	return r.top - r.bottom
}

// Left returns left component
func (r *Rectangle) Left() float32 {
	return r.left
}

// Right returns right component
func (r *Rectangle) Right() float32 {
	return r.right
}

// Top returns top component
func (r *Rectangle) Top() float32 {
	return r.top
}

// Bottom returns bottom component
func (r *Rectangle) Bottom() float32 {
	return r.bottom
}

// PointContained checks point using left-top rule. Point is NOT
// inside if on an Edge.
func (r *Rectangle) PointContained(p api.IPoint) bool {
	return p.X() > r.left && p.X() < (r.right) && p.Y() > r.bottom && p.Y() < (r.top)
}

// PointInside checks point using left-top rule. Point is
// inside if on an top or left Edge and NOT on a bottom or right edge
func (r *Rectangle) PointInside(p api.IPoint) bool {
	return p.X() >= r.left && p.X() < (r.right) && p.Y() >= r.bottom && p.Y() < (r.top)
}

// Intersects returns true if other (o) intersects this rectangle.
func (r *Rectangle) Intersects(o api.IRectangle) bool {
	return r.left <= o.Right() &&
		o.Left() <= r.right &&
		r.bottom <= o.Top() &&
		o.Bottom() <= r.top
}

// Contains returns true if r completely contains o.
func (r *Rectangle) Contains(o api.IRectangle) bool {
	return r.left <= o.Left() &&
		r.right >= o.Right() &&
		r.top >= o.Top() &&
		r.bottom <= o.Bottom()
}

func (r Rectangle) String() string {
	return fmt.Sprintf("LB(%0.3f,%0.3f) - RT(%0.3f,%0.3f) - Size: %0.3f x %0.3f", r.left, r.bottom, r.right, r.top, r.width, r.height)
}
