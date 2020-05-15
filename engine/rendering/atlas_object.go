package rendering

import "github.com/wdevore/Ranger-Go-IGE/api"

// AtlasObject helps managing a Mesh. It is abstract and
// should be embedded.
type AtlasObject struct {
	isStatic bool

	prevComponentCount int
	// componentCount counts how many vertices have been added
	ComponentCount int
	idx            int
	prevIndexCount int

	mesh api.IMesh
}

// No Allocator as this type is abstract and meant to
// be embedded

// Initialize sets defaults
func (va *AtlasObject) Initialize(isStatic bool) {
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

// SetVertex modifies a vertex in a mesh. The vertices still need
// to be copied to the graphics buffer using the VBOs Update(...)
func (va *AtlasObject) SetVertex(x, y, z float32, index int) {
	va.mesh.SetVertex(x, y, z, index)
}

// AddIndex adds an index to the mesh
func (va *AtlasObject) AddIndex(index int) {
	va.mesh.AddIndex(index)
	va.idx++
}

// Begin configures for a new sequence of vertices and indices
func (va *AtlasObject) Begin() int {
	va.prevComponentCount = va.ComponentCount
	va.prevIndexCount = va.idx
	return va.prevIndexCount
}

// End closes sequence of vertices and indices
func (va *AtlasObject) End() int {
	return va.idx - va.prevIndexCount
}

// Mesh returns this vector atlas's mesh
func (va *AtlasObject) Mesh() api.IMesh {
	return va.mesh
}
