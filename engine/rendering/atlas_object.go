package rendering

import "github.com/wdevore/Ranger-Go-IGE/api"

// AtlasObject helps managing a Mesh. It is abstract and
// should be embedded.
type AtlasObject struct {
	hasNormals bool
	hasColors  bool
	isStatic   bool
	// vertexIdx          int
	// vertexSize         int

	prevComponentCount int
	// componentCount counts how many vertices have been added
	ComponentCount int
	Idx            int
	prevIndexCount int

	mesh api.IMesh
}

// No Allocator as this type is abstract and meant to
// be embedded

// Initialize sets defaults
func (va *AtlasObject) Initialize(isStatic, hasColors bool) {
	va.hasColors = hasColors
	va.isStatic = isStatic
	va.mesh = NewMesh()
}

// AddVertex adds a vertex to the mesh
func (va *AtlasObject) AddVertex(x, y, z float32) int {
	va.mesh.AddVertex(x, y, z)
	idx := va.ComponentCount
	va.ComponentCount++
	return idx
}

// AddIndex adds an index to the mesh
func (va *AtlasObject) AddIndex(index int) {
	va.mesh.AddIndex(index)
	va.Idx++
}

// Begin configures for a new sequence of vertices and indices
func (va *AtlasObject) Begin() int {
	va.prevComponentCount = va.ComponentCount
	va.prevIndexCount = va.Idx
	return va.prevIndexCount
}

// End closes sequence of vertices and indices
func (va *AtlasObject) End() int {
	return va.Idx - va.prevIndexCount
}

// Mesh returns this vector atlas's mesh
func (va *AtlasObject) Mesh() api.IMesh {
	return va.mesh
}
