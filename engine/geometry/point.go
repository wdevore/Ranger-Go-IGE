package geometry

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
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

func (p point) String() string {
	return fmt.Sprintf("(%0.3f,%0.3f)", p.x, p.y)
}
