package geometry

import (
	"fmt"
	"math"

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

// SetByRectangle by rectangle
func (r *Rectangle) SetByRectangle(rect api.IRectangle) {
	r.left = rect.Left()
	r.bottom = rect.Bottom()
	r.top = rect.Top()
	r.right = rect.Right()
	r.width = rect.Width()
	r.height = rect.Height()
}

// SetBySize sets the bottom-left to origin and top-right to w/h.
func (r *Rectangle) SetBySize(width, height float32) {
	r.bottom = 0.0
	r.left = 0.0
	r.top = height
	r.right = width
	r.width = width
	r.height = height
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

// SetBounds2D set the min/max corners based on array of vertices.
// [x,y,x,y...]
func (r *Rectangle) SetBounds2D(vertices []float32) {
	minX := math.MaxFloat32
	minY := math.MaxFloat32
	maxX := -math.MaxFloat32
	maxY := -math.MaxFloat32

	for i := 0; i < len(vertices); i += 2 {
		x := float64(vertices[i])
		y := float64(vertices[i+1])

		minX = math.Min(minX, x)
		minY = math.Min(minY, y)

		maxX = math.Max(maxX, x)
		maxY = math.Max(maxY, y)
	}

	r.SetMinMax(float32(minX), float32(minY), float32(maxX), float32(maxY))
}

// SetBounds3D set the min/max corners based on array of vertices.
// [x,y,z,x,y,z...]
func (r *Rectangle) SetBounds3D(vertices []float32) {
	minX := math.MaxFloat32
	minY := math.MaxFloat32
	maxX := -math.MaxFloat32
	maxY := -math.MaxFloat32

	for i := 0; i < len(vertices); i += 3 {
		x := float64(vertices[i])
		y := float64(vertices[i+1])

		minX = math.Min(minX, x)
		minY = math.Min(minY, y)

		maxX = math.Max(maxX, x)
		maxY = math.Max(maxY, y)
	}

	r.SetMinMax(float32(minX), float32(minY), float32(maxX), float32(maxY))
	r.SetCenter(r.left+r.width/2.0, r.bottom+r.height/2.0)
}

// Expand resizes bounds based on x,y
func (r *Rectangle) Expand(x, y float32) {
	minX := float64(r.left)
	minY := float64(r.bottom)
	maxX := float64(r.right)
	maxY := float64(r.top)

	minX = math.Min(minX, float64(x))
	minY = math.Min(minY, float64(y))

	maxX = math.Max(maxX, float64(x))
	maxY = math.Max(maxY, float64(y))

	r.SetMinMax(float32(minX), float32(minY), float32(maxX), float32(maxY))
	r.SetCenter(r.left+r.width/2.0, r.bottom+r.height/2.0)
}

// Area return bounds area
func (r *Rectangle) Area() float32 {
	return r.width * r.height
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
	return fmt.Sprintf(
		"BL(%0.3f,%0.3f) - TR(%0.3f,%0.3f) - Size: %0.3f x %0.3f | At: (%0.3f,%0.3f)",
		r.left, r.bottom, r.right, r.top, r.width, r.height,
		r.center.X(), r.center.Y(),
	)
}
