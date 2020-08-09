package rendering

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// ShapeAtlas is a map-collection of vector shapes managed by a vector object
type ShapeAtlas struct {
	shapes []api.IAtlasShape
}

// NewShapeAtlas creates an atlas to be populated
func NewShapeAtlas() api.IAtlas {
	o := new(ShapeAtlas)
	o.initialize()
	return o
}

// Initialize this embedded object
func (a *ShapeAtlas) initialize() {
	a.shapes = []api.IAtlasShape{}
}

// GenerateShape creates a shape, adds it to the atlas
func (a *ShapeAtlas) GenerateShape(atlasName string, primitiveMode uint32) api.IAtlasShape {
	// If this shape is already present just return what's in the collection.
	for _, s := range a.shapes {
		if s.Name() == atlasName {
			return s
		}
	}
	s := NewAtlasShape()
	s.SetName(atlasName)
	s.SetPrimitiveMode(primitiveMode)
	a.shapes = append(a.shapes, s)
	return s
}

// HasShapes returns true if the atlas contains shapes
func (a *ShapeAtlas) HasShapes() bool {
	return len(a.shapes) > 0
}

// Shapes retuns shape collection
func (a *ShapeAtlas) Shapes() []api.IAtlasShape {
	return a.shapes
}

// AddShape adds a vector shape to the collection
func (a *ShapeAtlas) AddShape(shape api.IAtlasShape) {
}
