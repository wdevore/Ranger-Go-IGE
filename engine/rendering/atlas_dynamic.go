package rendering

import (
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// The Atlas is pre-populated by
// objects that can be referenced by the following names:
// - Line

const (
	// VboElementCountPerLine - Each Line has 2 vertices (components) and each component
	// has 3 elements (xyz)
	VboElementCountPerLine = 6

	// EboElementCountPerComponent - Each line has 2 indices (elements) = 1 EBO component.
	// Each EBO component represents 2 VBO components where each
	// EBO element refers to 1 VBO component.
	EboElementCountPerComponent = 2
)

type dynamicAtlas struct {
	Atlas

	elementIndex int
}

// NewDynamicAtlas creates an atlas to be populated
func NewDynamicAtlas() api.IAtlas {
	o := new(dynamicAtlas)
	o.Initialize()
	return o
}

func (a *dynamicAtlas) Populate(atlasObj api.IAtlasObject) {
	// Add two lines

	backingArrayIdx := atlasObj.AddArray()
	a.AddShape(a.buildLine(atlasObj, backingArrayIdx))

	backingArrayIdx = atlasObj.AddArray()
	a.AddShape(a.buildLine(atlasObj, backingArrayIdx))

	atlasObj.AddIndex(0)
	atlasObj.AddIndex(1)
}

func (a *dynamicAtlas) buildLine(atlasObj api.IAtlasObject, backingArrayIdx int) api.IAtlasShape {
	s := NewAtlasShape(atlasObj)

	s.SetBackingArrayIdx(backingArrayIdx)
	s.SetName(fmt.Sprintf("Line-%d", backingArrayIdx))
	s.SetPrimitiveMode(gl.LINES)

	s.SetOffset(atlasObj.Begin())

	// These vertices are specified in unit local-space
	// These values are simply starter values. The app will changes these dynamically.
	atlasObj.AddVertex(0.0, 0.0, 0.0)
	atlasObj.AddVertex(0.0, 0.0, 0.0)

	// This is the total element count where an element is an "indice"
	s.SetCount(atlasObj.End())

	// Lines require 2 elements where "element" refers to an EBO
	// item. In this case a "begin" index and "end" index
	s.SetElementCount(EboElementCountPerComponent)

	return s
}

func (a *dynamicAtlas) GetNextIndex(glType int) int {
	id := a.elementIndex
	a.elementIndex += 2
	return id
}
