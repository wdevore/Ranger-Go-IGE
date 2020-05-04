package atlas

import (
	"math"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

// BasicAtlas is an Atlas of basic vector shapes, for example, Square or Circle.
type BasicAtlas struct {
	Atlas
}

// NewBasicAtlas creates a basic atlas.
func NewBasicAtlas() *BasicAtlas {
	ba := new(BasicAtlas)
	return ba
}

// Populate loads atlas with a starter set of basic shapes
func (ba *BasicAtlas) Populate(vo *rendering.VectorObject) {
	ba.initialize()
	uAtlas := vo.UniAtlas
	ba.AddShape(buildSquare(uAtlas))
	ba.AddShape(buildCenteredSquare(uAtlas))
	ba.AddShape(buildCenteredTriangle(uAtlas))
}

func buildSquare(uAtlas *rendering.VectorUniformAtlas) *rendering.VectorShape {
	s := rendering.NewVectorShape()
	s.Name = "Square"
	s.PrimitiveMode = gl.TRIANGLES

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

	s.Count = int32(uAtlas.End())

	return s
}

func buildCenteredSquare(uAtlas *rendering.VectorUniformAtlas) *rendering.VectorShape {
	s := rendering.NewVectorShape()
	s.Name = "CenteredSquare"
	s.PrimitiveMode = gl.TRIANGLES

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

	s.Count = int32(uAtlas.End())

	return s
}

func buildCenteredTriangle(uAtlas *rendering.VectorUniformAtlas) *rendering.VectorShape {
	s := rendering.NewVectorShape()
	s.Name = "CenteredTriangle"
	s.PrimitiveMode = gl.TRIANGLES

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

	s.Count = int32(uAtlas.End())

	return s
}
