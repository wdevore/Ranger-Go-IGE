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
	index int
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
	a.AddShape(a.buildCross(atlasObj))
	a.AddShape(a.buildCrowBar(atlasObj))

	// Outline shapes
	a.AddShape(a.buildHLine(atlasObj))
	a.AddShape(a.buildVLine(atlasObj))
	a.AddShape(a.buildOutlineCenteredSquare(atlasObj))
	a.AddShape(a.buildCenteredOutlineTriangle(atlasObj))
	a.AddShape(a.buildCenteredOutlineCircle(atlasObj))
}

func (a *staticAtlas) buildPixel(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("Pixel")
	s.SetPrimitiveMode(gl.POINTS)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	v0 := atlasObj.AddVertex(0.0, 0.0, 0.0)

	atlasObj.AddIndex(v0)

	s.SetCount(atlasObj.End())

	return s
}

func (a *staticAtlas) GetNextIndex(glType int) int {
	id := a.index
	switch glType {
	case api.GLLines:
		a.index += 2
	}
	return id
}

func (a *staticAtlas) buildCross(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("Cross")
	s.SetPrimitiveMode(gl.LINES)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	v0 := atlasObj.AddVertex(-0.5, 0.0, 0.0)
	v1 := atlasObj.AddVertex(0.5, 0.0, 0.0)
	v2 := atlasObj.AddVertex(0.0, -0.5, 0.0)
	v3 := atlasObj.AddVertex(0.0, 0.5, 0.0)

	atlasObj.AddIndex(v0)
	atlasObj.AddIndex(v1)
	atlasObj.AddIndex(v2)
	atlasObj.AddIndex(v3)

	s.SetCount(atlasObj.End())

	return s
}

func (a *staticAtlas) buildCircle(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("Circle12Segments")
	s.SetPrimitiveMode(gl.TRIANGLE_FAN)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	v0 := atlasObj.AddVertex(0.0, 0.0, 0.0) // apex
	atlasObj.AddIndex(v0)

	segments := 12
	radius := 0.5 // diameter of 1.0
	step := math.Pi / float64(segments)

	for i := 0.0; i < 2.0*math.Pi; i += step {
		x := math.Cos(i) * radius
		y := math.Sin(i) * radius
		idx := atlasObj.AddVertex(float32(x), float32(y), 0.0)
		atlasObj.AddIndex(idx)
	}

	s.SetCount(atlasObj.End())

	return s
}

func (a *staticAtlas) buildSquare(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("Square")
	s.SetPrimitiveMode(gl.TRIANGLES)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	v0 := atlasObj.AddVertex(0.0, 0.0, 0.0)
	v1 := atlasObj.AddVertex(0.0, 1.0, 0.0)
	v2 := atlasObj.AddVertex(1.0, 1.0, 0.0)
	v3 := atlasObj.AddVertex(1.0, 0.0, 0.0)

	atlasObj.AddIndex(v0)
	atlasObj.AddIndex(v1)
	atlasObj.AddIndex(v3)
	atlasObj.AddIndex(v1)
	atlasObj.AddIndex(v2)
	atlasObj.AddIndex(v3)

	s.SetCount(atlasObj.End())

	return s
}

func (a *staticAtlas) buildCenteredSquare(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("CenteredSquare")
	s.SetPrimitiveMode(gl.TRIANGLES)

	s.SetOffset(atlasObj.Begin())

	const l float32 = 0.5 // side length

	// These vertices are specified in unit local-space
	v0 := atlasObj.AddVertex(l, l, 0.0)
	v1 := atlasObj.AddVertex(l, -l, 0.0)
	v2 := atlasObj.AddVertex(-l, -l, 0.0)
	v3 := atlasObj.AddVertex(-l, l, 0.0)

	atlasObj.AddIndex(v0)
	atlasObj.AddIndex(v3)
	atlasObj.AddIndex(v1)
	atlasObj.AddIndex(v1)
	atlasObj.AddIndex(v3)
	atlasObj.AddIndex(v2)

	s.SetCount(atlasObj.End())

	return s
}

func (a *staticAtlas) buildCenteredTriangle(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("CenteredTriangle")
	s.SetPrimitiveMode(gl.TRIANGLES)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	v0 := atlasObj.AddVertex(-0.5, -0.5, 0.0)
	v1 := atlasObj.AddVertex(0.5, -0.5, 0.0)
	v2 := atlasObj.AddVertex(0.0, 0.314, 0.0)

	atlasObj.AddIndex(v0)
	atlasObj.AddIndex(v1)
	atlasObj.AddIndex(v2)

	s.SetCount(atlasObj.End())

	return s
}

func (a *staticAtlas) buildCrowBar(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("CrowBar")
	s.SetPrimitiveMode(gl.TRIANGLES)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	v0 := atlasObj.AddVertex(-0.1, -0.5, 0.0)
	v1 := atlasObj.AddVertex(0.5, -0.5, 0.0)
	v2 := atlasObj.AddVertex(0.5, -0.4, 0.0)
	v3 := atlasObj.AddVertex(0.1, -0.4, 0.0)
	v4 := atlasObj.AddVertex(0.1, 0.5, 0.0)
	v5 := atlasObj.AddVertex(-0.5, 0.5, 0.0)
	v6 := atlasObj.AddVertex(-0.5, 0.4, 0.0)
	v7 := atlasObj.AddVertex(-0.1, 0.4, 0.0)

	atlasObj.AddIndex(v0)
	atlasObj.AddIndex(v1)
	atlasObj.AddIndex(v2)

	atlasObj.AddIndex(v2)
	atlasObj.AddIndex(v3)
	atlasObj.AddIndex(v0)

	atlasObj.AddIndex(v0)
	atlasObj.AddIndex(v3)
	atlasObj.AddIndex(v4)

	atlasObj.AddIndex(v0)
	atlasObj.AddIndex(v4)
	atlasObj.AddIndex(v7)

	atlasObj.AddIndex(v6)
	atlasObj.AddIndex(v7)
	atlasObj.AddIndex(v4)

	atlasObj.AddIndex(v6)
	atlasObj.AddIndex(v4)
	atlasObj.AddIndex(v5)

	s.SetCount(atlasObj.End())

	return s
}

// --------------------------------------------------------------------
// Outline shapes
// --------------------------------------------------------------------

func (a *staticAtlas) buildHLine(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("HLine")
	s.SetPrimitiveMode(gl.LINES)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	v0 := atlasObj.AddVertex(-0.5, 0.0, 0.0)
	v1 := atlasObj.AddVertex(0.5, 0.0, 0.0)

	atlasObj.AddIndex(v0)
	atlasObj.AddIndex(v1)

	s.SetCount(atlasObj.End())

	return s
}

func (a *staticAtlas) buildVLine(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("VLine")
	s.SetPrimitiveMode(gl.LINES)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	v0 := atlasObj.AddVertex(0.0, -0.5, 0.0)
	v1 := atlasObj.AddVertex(0.0, 0.5, 0.0)

	atlasObj.AddIndex(v0)
	atlasObj.AddIndex(v1)

	s.SetCount(atlasObj.End())

	return s
}

func (a *staticAtlas) buildOutlineCenteredSquare(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("CenteredOutlineSquare")
	s.SetPrimitiveMode(gl.LINE_LOOP)

	s.SetOffset(atlasObj.Begin())

	const l float32 = 0.5 // side length

	// These vertices are specified in unit local-space and
	// in CCW order.
	v0 := atlasObj.AddVertex(-l, -l, 0.0)
	v1 := atlasObj.AddVertex(l, -l, 0.0)
	v2 := atlasObj.AddVertex(l, l, 0.0)
	v3 := atlasObj.AddVertex(-l, l, 0.0)

	atlasObj.AddIndex(v0)
	atlasObj.AddIndex(v1)
	atlasObj.AddIndex(v2)
	atlasObj.AddIndex(v3)

	s.SetCount(atlasObj.End())

	return s
}

func (a *staticAtlas) buildCenteredOutlineTriangle(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("CenteredOutlineTriangle")
	s.SetPrimitiveMode(gl.LINE_LOOP)

	s.SetOffset(atlasObj.Begin())

	const l float32 = 0.5 // side length

	// These vertices are specified in unit local-space and
	// in CCW order.
	v0 := atlasObj.AddVertex(-l, -l, 0.0)
	v1 := atlasObj.AddVertex(l, -l, 0.0)
	v2 := atlasObj.AddVertex(0.0, 0.314, 0.0)

	atlasObj.AddIndex(v0)
	atlasObj.AddIndex(v1)
	atlasObj.AddIndex(v2)

	s.SetCount(atlasObj.End())

	return s
}

func (a *staticAtlas) buildCenteredOutlineCircle(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("Circle12SegmentsOutline")
	s.SetPrimitiveMode(gl.LINE_LOOP)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	// and in CCW
	segments := 12
	radius := 0.5 // diameter of 1.0
	step := math.Pi / float64(segments)

	for i := 0.0; i < 2.0*math.Pi; i += step {
		x := math.Cos(i) * radius
		y := math.Sin(i) * radius
		idx := atlasObj.AddVertex(float32(x), float32(y), 0.0)
		atlasObj.AddIndex(idx)
	}

	s.SetCount(atlasObj.End())

	return s
}
