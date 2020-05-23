package rendering

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// This atlas has one item--a pixel buffer.

type pixelAtlas struct {
	Atlas
}

// NewPixelAtlas creates an atlas to be populated
func NewPixelAtlas() api.IAtlas {
	o := new(pixelAtlas)
	o.Initialize()
	return o
}

func (p *pixelAtlas) Populate(atlasObj api.IAtlasObject) {
	s := NewAtlasShape(atlasObj)
	s.SetName("PixelBuffer")
	s.SetPrimitiveMode(gl.POINTS)

	s.SetOffset(atlasObj.Begin())
	s.SetMaxSize(500)

	// These vertices are specified in unit local-space
	// The app will change these dynamically.
	for i := 0; i < s.MaxSize(); i++ {
		vi := atlasObj.AddVertex(0.0, 0.0, 0.0)
		atlasObj.AddIndex(vi)
	}

	s.SetCount(atlasObj.End())
	p.AddShape(s)
}

func (p *pixelAtlas) GetNextIndex(glType int) int {
	return 0
}
