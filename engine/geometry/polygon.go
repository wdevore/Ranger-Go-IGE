package geometry

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// Polygon is a mesh with additional methods
type Polygon struct {
	vertices []api.IPoint
}

// NewPolygon constructs a new IPolygon
func NewPolygon() api.IPolygon {
	o := new(Polygon)
	o.vertices = []api.IPoint{}
	return o
}

// AddVertex adds a point to the mesh
func (p *Polygon) AddVertex(x, y float32) {
	p.vertices = append(p.vertices, NewPointUsing(x, y))
}

// SetVertex updates a point on the mesh
func (p *Polygon) SetVertex(x, y float32, index int) {
	p.vertices[index].SetByComp(x, y)
}

// PointInside will return false if the point is on the right/bottom edge and/or outward
// if point is on the right and/or bottom edge it is considered outside and
// the left/top edge is considered inside.
// This is consistant with polygon filling.
func (p *Polygon) PointInside(po api.IPoint) bool {
	i := 0
	c := false
	vertices := p.vertices
	nvert := len(vertices)
	j := 1

	for j < nvert {
		if ((vertices[i].Y() > po.Y()) != (vertices[j].Y() > po.Y())) &&
			(po.X() < (vertices[j].X()-vertices[i].X())*(po.Y()-vertices[i].Y())/
				(vertices[j].Y()-vertices[i].Y())+vertices[i].X()) {
			c = !c
		}
		i++
		j++
	}

	// Last edge to close loop
	i = j - 1
	j = 0
	if ((vertices[i].Y() > po.Y()) != (vertices[j].Y() > po.Y())) &&
		(po.X() < (vertices[j].X()-vertices[i].X())*(po.Y()-vertices[i].Y())/
			(vertices[j].Y()-vertices[i].Y())+vertices[i].X()) {
		c = !c
	}

	return c
}

func (p Polygon) String() string {
	return fmt.Sprintf("<%7.3f,%7.3f>", p.vertices[0].X(), p.vertices[0].Y())
}
