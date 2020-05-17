package rendering

import (
	"math"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// The Atlas is pre-populated by
// objects that can be referenced by the following names:
// - Pixel
// - Line
// - Cross
// - Circle12Segments
// - Square
// - CenteredSquare
// - CenteredTriangle
// - CrowBar

type staticAtlas struct {
	Atlas
}

// NewStaticAtlas creates an atlas to be populated
func NewStaticAtlas() api.IAtlas {
	o := new(staticAtlas)
	o.Initialize()
	return o
}

func (a *staticAtlas) Populate(atlasObj api.IAtlasObject) {
	a.AddShape(a.buildPixel(atlasObj))
	a.AddShape(a.buildSquare(atlasObj))
	a.AddShape(a.buildCenteredSquare(atlasObj))
	a.AddShape(a.buildCenteredTriangle(atlasObj))
	a.AddShape(a.buildCircle(atlasObj))
	a.AddShape(a.buildLine(atlasObj))
	a.AddShape(a.buildCross(atlasObj))
	a.AddShape(a.buildCrowBar(atlasObj))
}

func (a *staticAtlas) buildPixel(uAtlas api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape()
	s.SetName("Pixel")
	s.SetPrimitiveMode(gl.POINTS)

	s.SetOffset(uAtlas.Begin())

	// These vertices are specified in unit local-space
	v0 := uAtlas.AddVertex(0.0, 0.0, 0.0)

	uAtlas.AddIndex(v0)

	s.SetCount(int32(uAtlas.End()))

	return s
}

func (a *staticAtlas) buildLine(uAtlas api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape()
	s.SetName("Line")
	s.SetPrimitiveMode(gl.LINES)

	s.SetOffset(uAtlas.Begin())

	// These vertices are specified in unit local-space
	v0 := uAtlas.AddVertex(-0.5, 0.0, 0.0)
	v1 := uAtlas.AddVertex(0.5, 0.0, 0.0)

	uAtlas.AddIndex(v0)
	uAtlas.AddIndex(v1)

	s.SetCount(int32(uAtlas.End()))

	return s
}

func (a *staticAtlas) buildCross(uAtlas api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape()
	s.SetName("Cross")
	s.SetPrimitiveMode(gl.LINES)

	s.SetOffset(uAtlas.Begin())

	// These vertices are specified in unit local-space
	v0 := uAtlas.AddVertex(-0.5, 0.0, 0.0)
	v1 := uAtlas.AddVertex(0.5, 0.0, 0.0)
	v2 := uAtlas.AddVertex(0.0, -0.5, 0.0)
	v3 := uAtlas.AddVertex(0.0, 0.5, 0.0)

	uAtlas.AddIndex(v0)
	uAtlas.AddIndex(v1)
	uAtlas.AddIndex(v2)
	uAtlas.AddIndex(v3)

	s.SetCount(int32(uAtlas.End()))

	return s
}

func (a *staticAtlas) buildCircle(uAtlas api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape()
	s.SetName("Circle12Segments")
	s.SetPrimitiveMode(gl.TRIANGLE_FAN)

	s.SetOffset(uAtlas.Begin())

	// These vertices are specified in unit local-space
	v0 := uAtlas.AddVertex(0.0, 0.0, 0.0) // apex
	uAtlas.AddIndex(v0)

	segments := 12
	radius := 0.5 // diameter of 1.0
	step := math.Pi / float64(segments)

	for i := 0.0; i < 2.0*math.Pi; i += step {
		x := math.Cos(i) * radius
		y := math.Sin(i) * radius
		idx := uAtlas.AddVertex(float32(x), float32(y), 0.0)
		uAtlas.AddIndex(idx)
	}

	s.SetCount(int32(uAtlas.End()))

	return s
}

func (a *staticAtlas) buildSquare(uAtlas api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape()
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

func (a *staticAtlas) buildCenteredSquare(uAtlas api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape()
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

func (a *staticAtlas) buildCenteredTriangle(uAtlas api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape()
	s.SetName("CenteredTriangle")
	s.SetPrimitiveMode(gl.TRIANGLES)

	s.SetOffset(uAtlas.Begin())

	// These vertices are specified in unit local-space
	v0 := uAtlas.AddVertex(-0.5, -0.5, 0.0)
	v1 := uAtlas.AddVertex(0.5, -0.5, 0.0)
	v2 := uAtlas.AddVertex(0.0, 0.314, 0.0)

	uAtlas.AddIndex(v0)
	uAtlas.AddIndex(v1)
	uAtlas.AddIndex(v2)

	s.SetCount(int32(uAtlas.End()))

	return s
}

func (a *staticAtlas) buildCrowBar(uAtlas api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape()
	s.SetName("CrowBar")
	s.SetPrimitiveMode(gl.TRIANGLES)

	s.SetOffset(uAtlas.Begin())

	// These vertices are specified in unit local-space
	v0 := uAtlas.AddVertex(-0.1, -0.5, 0.0)
	v1 := uAtlas.AddVertex(0.5, -0.5, 0.0)
	v2 := uAtlas.AddVertex(0.5, -0.4, 0.0)
	v3 := uAtlas.AddVertex(0.1, -0.4, 0.0)
	v4 := uAtlas.AddVertex(0.1, 0.5, 0.0)
	v5 := uAtlas.AddVertex(-0.5, 0.5, 0.0)
	v6 := uAtlas.AddVertex(-0.5, 0.4, 0.0)
	v7 := uAtlas.AddVertex(-0.1, 0.4, 0.0)

	uAtlas.AddIndex(v0)
	uAtlas.AddIndex(v1)
	uAtlas.AddIndex(v2)

	uAtlas.AddIndex(v2)
	uAtlas.AddIndex(v3)
	uAtlas.AddIndex(v0)

	uAtlas.AddIndex(v0)
	uAtlas.AddIndex(v3)
	uAtlas.AddIndex(v4)

	uAtlas.AddIndex(v0)
	uAtlas.AddIndex(v4)
	uAtlas.AddIndex(v7)

	uAtlas.AddIndex(v6)
	uAtlas.AddIndex(v7)
	uAtlas.AddIndex(v4)

	uAtlas.AddIndex(v6)
	uAtlas.AddIndex(v4)
	uAtlas.AddIndex(v5)

	s.SetCount(int32(uAtlas.End()))

	return s
}
