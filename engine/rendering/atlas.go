package rendering

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// Atlas is a map-collection of vector shapes managed by a vector object
type Atlas struct {
	shapes []api.IAtlasShape
}

// NewAtlas creates an atlas to be populated
func NewAtlas() api.IAtlas {
	o := new(Atlas)
	o.initialize()
	return o
}

// Initialize this embedded object
func (a *Atlas) initialize() {
	a.shapes = []api.IAtlasShape{}
}

// GenerateShape creates a shape, adds it to the atlas
func (a *Atlas) GenerateShape(atlasName string, primitiveMode uint32, bufferType bool) api.IAtlasShape {
	// If this shape is already present just return what's in the collection.
	for _, s := range a.shapes {
		if s.Name() == atlasName {
			return s
		}
	}
	s := NewAtlasShape()
	s.SetName(atlasName)
	s.SetPrimitiveMode(primitiveMode)
	// s.SetStatic(bufferType)
	a.shapes = append(a.shapes, s)
	return s
}

// HasShapes returns true if the atlas contains shapes
func (a *Atlas) HasShapes() bool {
	return len(a.shapes) > 0
}

// Shapes retuns shape collection
func (a *Atlas) Shapes() []api.IAtlasShape {
	return a.shapes
}

// AddShape adds a vector shape to the collection
func (a *Atlas) AddShape(shape api.IAtlasShape) {
}
