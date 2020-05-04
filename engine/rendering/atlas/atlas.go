// Package atlas defines vector shape collections
package atlas

import "github.com/wdevore/Ranger-Go-IGE/engine/rendering"

// Atlas is a map-collection of vector shapes managed by a vector object
type Atlas struct {
	Shapes map[string]*rendering.VectorShape
}

func (a *Atlas) initialize() {
	a.Shapes = make(map[string]*rendering.VectorShape)
}

// AddShape adds a vector shape to the collection
func (a *Atlas) AddShape(vs *rendering.VectorShape) {
	a.Shapes[vs.Name] = vs
}
