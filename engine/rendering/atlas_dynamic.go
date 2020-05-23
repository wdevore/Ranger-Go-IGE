package rendering

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// The Atlas is pre-populated by
// objects that can be referenced by the following names:
// - Line

const (
	// Each Line has 2 vertices (components) and each component
	// has 3 elements (xyz)
	vboElementCountPerLine = 6

	// Each line has 2 indices (elements) = 1 EBO component.
	// Each EBO component represents 2 VBO components where each
	// EBO element refers to 1 VBO component.
	eboElementCountPerComponent = 2
)

type dynamicAtlas struct {
	Atlas

	index int
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

func (a *dynamicAtlas) buildLine(atlasObj api.IAtlasObject) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)
	s.SetName("Line")
	s.SetPrimitiveMode(gl.LINES)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	// These values are simply starter values. The app will changes these dynamically.
	for i := 0; i < 2; i++ {
		v0 := atlasObj.AddVertex(0.0, 0.0, 0.0)
		atlasObj.AddIndex(v0) // Begin
		v1 := atlasObj.AddVertex(0.0, 0.0, 0.0)
		atlasObj.AddIndex(v1) // End
	}

	// This is the total element count where an element is an "indice"
	s.SetCount(atlasObj.End())

	// Lines require 2 elements where "element" refers to an EBO
	// item. In this case a "begin" index and "end" index
	s.SetElementCount(eboElementCountPerComponent)

	return s
}

func (a *dynamicAtlas) GetNextIndex(glType int) int {
	id := a.index
	a.index += 2
	return id
}
