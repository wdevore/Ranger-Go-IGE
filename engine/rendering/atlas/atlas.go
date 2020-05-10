// Package atlas defines vector shape collections
package atlas

import (
	"math"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

// Atlas is a map-collection of vector shapes managed by a vector object
type Atlas struct {
	shapes map[string]api.IVectorShape
}

// NewAtlas creates a new atlas
func NewAtlas() api.IAtlas {
	o := new(Atlas)
	o.shapes = make(map[string]api.IVectorShape)
	return o
}

// Initialize adds a few basic shapes to atlas
func (a *Atlas) Initialize(vo api.IVectorObject) {
	uAtlas := vo.UniformAtlas()

	a.AddShape(buildSquare(uAtlas))
	a.AddShape(buildCenteredSquare(uAtlas))
	a.AddShape(buildCenteredTriangle(uAtlas))
}

// Shape returns a shape by name
func (a *Atlas) Shape(name string) api.IVectorShape {
	return a.shapes[name]
}

// Shapes returns all the shapes the atlas contains
func (a *Atlas) Shapes() map[string]api.IVectorShape {
	return a.shapes
}

// AddShape adds a vector shape to the collection
func (a *Atlas) AddShape(vs api.IVectorShape) {
	a.shapes[vs.Name()] = vs
}

func buildSquare(uAtlas api.IVectorAtlas) api.IVectorShape {
	s := rendering.NewVectorShape()
	s.SetName("Square")
	s.SetPrimitiveMode(gl.TRIANGLES)

	s.SetOffset(uAtlas.Begin())

	// These vertices are specified in unit local-space
	v0 := uAtlas.AddVertex(0.0, 0.0, 0.0)
	v1 := uAtlas.AddVertex(0.0, 1.0, 0.0)
	v2 := uAtlas.AddVertex(1.0, 1.0, 0.0)
	v3 := uAtlas.AddVertex(1.0, 0.0, 0.0)

	uAtlas.AddIndex(v0)
	uAtlas.AddIndex(v1)
	uAtlas.AddIndex(v3)
	uAtlas.AddIndex(v1)
	uAtlas.AddIndex(v2)
	uAtlas.AddIndex(v3)

	s.SetCount(int32(uAtlas.End()))

	return s
}

func buildCenteredSquare(uAtlas api.IVectorAtlas) api.IVectorShape {
	s := rendering.NewVectorShape()
	s.SetName("CenteredSquare")
	s.SetPrimitiveMode(gl.TRIANGLES)

	s.SetOffset(uAtlas.Begin())

	const l float32 = 0.5 // side length

	// These vertices are specified in unit local-space
	v0 := uAtlas.AddVertex(l, l, 0.0)
	v1 := uAtlas.AddVertex(l, -l, 0.0)
	v2 := uAtlas.AddVertex(-l, -l, 0.0)
	v3 := uAtlas.AddVertex(-l, l, 0.0)

	uAtlas.AddIndex(v0)
	uAtlas.AddIndex(v3)
	uAtlas.AddIndex(v1)
	uAtlas.AddIndex(v1)
	uAtlas.AddIndex(v3)
	uAtlas.AddIndex(v2)

	s.SetCount(int32(uAtlas.End()))

	return s
}

func buildCenteredTriangle(uAtlas api.IVectorAtlas) api.IVectorShape {
	s := rendering.NewVectorShape()
	s.SetName("CenteredTriangle")
	s.SetPrimitiveMode(gl.TRIANGLES)

	s.SetOffset(uAtlas.Begin())

	const l float32 = 0.25 // side length

	// 30 degrees yields triangle sides of equal length but the bbox is
	// rectangular not square.
	// 0 degrees yields a square bbox with unequal triangle sides.
	h := float32(0.5 * math.Cos(maths.DegreeToRadians*30.0))

	// These vertices are specified in unit local-space
	v0 := uAtlas.AddVertex(-l, -h, 0.0)
	v1 := uAtlas.AddVertex(l, -h, 0.0)
	v2 := uAtlas.AddVertex(0.0, h, 0.0)

	uAtlas.AddIndex(v0)
	uAtlas.AddIndex(v1)
	uAtlas.AddIndex(v2)

	s.SetCount(int32(uAtlas.End()))

	return s
}
