package rendering

import (
	"strings"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// Atlas is a map-collection of vector shapes managed by a vector object
type Atlas struct {
	shapes map[string]api.IAtlasShape
}

// Initialize this embedded object
func (a *Atlas) Initialize() {
	a.shapes = make(map[string]api.IAtlasShape)
}

// Shape returns a shape by name
func (a *Atlas) Shape(name string) api.IAtlasShape {
	return a.shapes[name]
}

// GetNextShape returns the next available shape by category
func (a *Atlas) GetNextShape(category string) api.IAtlasShape {
	var shape api.IAtlasShape = nil

	for name, shp := range a.shapes {
		cat := strings.Split(name, "-")
		if cat[0] == category && !shp.InUse() {
			shp.SetInUse(true)
			shape = shp
			break
		}
	}

	return shape
}

// Shapes returns all the shapes the atlas contains
func (a *Atlas) Shapes() map[string]api.IAtlasShape {
	return a.shapes
}

// AddShape adds a vector shape to the collection
func (a *Atlas) AddShape(shape api.IAtlasShape) {
	a.shapes[shape.Name()] = shape
}
