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
	a.AddShape(a.buildPixelBuffer(atlasObj))
}

func (a *dynamicAtlas) buildLine(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("Line")
	s.SetPrimitiveMode(gl.LINES)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	// These values are simply starter values. The app will changes these dynamically.
	v0 := atlasObj.AddVertex(-1.0, -1.0, 0.0)
	v1 := atlasObj.AddVertex(1.0, 1.0, 0.0)

	atlasObj.AddIndex(v0)
	atlasObj.AddIndex(v1)

	s.SetCount(atlasObj.End())

	return s
}

func (a *dynamicAtlas) buildPixelBuffer(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("PixelBuffer")
	s.SetPrimitiveMode(gl.POINTS)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	// The app will changes these dynamically.
	for i := 0; i < 1000; i++ {
		vi := atlasObj.AddVertex(0.0, 0.0, 0.0)
		atlasObj.AddIndex(vi)
	}

	s.SetCount(atlasObj.End())

	return s
}
