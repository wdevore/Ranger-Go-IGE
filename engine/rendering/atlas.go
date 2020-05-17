package rendering

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// Atlas is a map-collection of vector shapes managed by a vector object
type Atlas struct {
	shapes map[string]api.IAtlasShape
}

// NewAtlas creates a new dictionary of shapes. The Atlas is pre-populated by
// objects that can be referenced by the following names:
// - Pixel
// - Line
// - Cross
// - Circle12Segments
// - Square
// - CenteredSquare
// - CenteredTriangle
// - CrowBar
// func NewAtlas() api.IAtlas {
// 	o := new(Atlas)
// 	return o
// }

// Initialize this embedded object
func (a *Atlas) Initialize() {
	a.shapes = make(map[string]api.IAtlasShape)
}

// Shape returns a shape by name
func (a *Atlas) Shape(name string) api.IAtlasShape {
	return a.shapes[name]
}

// Shapes returns all the shapes the atlas contains
func (a *Atlas) Shapes() map[string]api.IAtlasShape {
	return a.shapes
}

// AddShape adds a vector shape to the collection
func (a *Atlas) AddShape(shape api.IAtlasShape) {
	a.shapes[shape.Name()] = shape
}
