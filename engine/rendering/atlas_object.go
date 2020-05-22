package rendering

import "github.com/wdevore/Ranger-Go-IGE/api"

// AtlasObject helps managing a Mesh. It is abstract and
// should be embedded.
type AtlasObject struct {
	// isStatic bool

	prevVertexCount int
	vertexCount     int
	indexCount      int
	prevIndexCount  int

	mesh api.IMesh
}

// No Allocator as this type is abstract and meant to
// be embedded

// Initialize sets defaults
func (a *AtlasObject) Initialize(isStatic bool) {
	// a.isStatic = isStatic
	a.mesh = NewMesh(isStatic)
}

// AddVertex adds a vertex to the mesh
func (a *AtlasObject) AddVertex(x, y, z float32) int {
	a.mesh.AddVertex(x, y, z)
	idx := a.vertexCount
	a.vertexCount++
	return idx
}

// SetVertex modifies a vertex in a mesh. The vertices still need
// to be copied to the graphics buffer using the VBOs Update(...)
func (a *AtlasObject) SetVertex(x, y, z float32, index int) {
	a.mesh.SetVertex(x, y, z, index)
}

// AddIndex adds an index to the mesh
func (a *AtlasObject) AddIndex(index int) {
	a.mesh.AddIndex(index)
	a.indexCount++
}

// Begin configures for a new sequence of vertices and indices
func (a *AtlasObject) Begin() int {
	a.prevVertexCount = a.vertexCount
	a.prevIndexCount = a.indexCount
	return a.prevIndexCount
}

// End closes sequence of vertices and indices
func (a *AtlasObject) End() int {
	return a.indexCount - a.prevIndexCount
}

// Mesh returns this vector atlas's mesh
func (a *AtlasObject) Mesh() api.IMesh {
	return a.mesh
}

// Update modifies the VBO buffer
func (a *AtlasObject) Update(offset, vertexCount int) {
	a.mesh.Update(offset, vertexCount)
}

func (a *AtlasObject) UpdatePreScaled(offset, vertexCount int) {
	a.mesh.UpdatePreScaled(offset, vertexCount)
}
