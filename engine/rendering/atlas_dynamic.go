package rendering

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// The Atlas is pre-populated by
// objects that can be referenced by the following names:
// - Line

type dynamicAtlas struct {
	Atlas
}

// NewDynamicAtlas creates an atlas to be populated
func NewDynamicAtlas() api.IAtlas {
	o := new(dynamicAtlas)
	o.Initialize()
	return o
}

func (a *dynamicAtlas) Populate(atlasObj api.IAtlasObject) {
	a.AddShape(a.buildLine(atlasObj))
}

func (a *dynamicAtlas) buildLine(uAtlas api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape()
	s.SetName("Line")
	s.SetPrimitiveMode(gl.LINES)

	s.SetOffset(uAtlas.Begin())

	// These vertices are specified in unit local-space
	// These values are simply starter values. The app will changes these dynamically.
	v0 := uAtlas.AddVertex(-0.5, 0.0, 0.0)
	v1 := uAtlas.AddVertex(0.5, 0.0, 0.0)

	uAtlas.AddIndex(v0)
	uAtlas.AddIndex(v1)

	s.SetCount(int32(uAtlas.End()))

	return s
}
