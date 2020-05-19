package geometry

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
)

type point struct {
	x, y float32
}

// NewPoint constructs a new IPoint
func NewPoint() api.IPoint {
	o := new(point)
	return o
}

// NewPointUsing constructs a new IPoint using components
func NewPointUsing(x, y float32) api.IPoint {
	o := new(point)
	o.x = x
	o.y = y
	return o
}

func (p *point) Components() (x, y float32) {
	return p.x, p.y
}

func (p *point) ComponentsAsInt32() (x, y int32) {
	return int32(p.x), int32(p.y)
}

func (p *point) X() float32 {
	return p.x
}

func (p *point) Y() float32 {
	return p.y
}

func (p *point) SetByComp(x, y float32) {
	p.x = x
	p.y = y
}

func (p *point) SetByPoint(ip api.IPoint) {
	p.x = ip.X()
	p.y = ip.Y()
}

// MulPoint left-multiplies the point by the given matrix
// assuming the 3rd component = 0, fourth (w) component = 1.
// |M00 M01 M02 M03|
// |M10 M11 M12 M13|
// |M20 M21 M22 M23|
// |M30 M31 M32 M33|
func (p *point) MulPoint(m api.IMatrix4) {
	me := m.Matrix()

	// vector = x, y, 0, 1
	p.SetByComp(
		p.x*me[maths.M00]+p.y*me[maths.M01]+me[maths.M03],
		p.x*me[maths.M10]+p.y*me[maths.M11]+me[maths.M13])
}

func (p point) String() string {
	return fmt.Sprintf("(%0.3f,%0.3f)", p.x, p.y)
}
