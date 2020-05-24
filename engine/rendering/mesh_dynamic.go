// Package rendering defines Mesh features of shaders.
package rendering

import "github.com/wdevore/Ranger-Go-IGE/api"

// DynamicMesh combines a shader's VBO and EBO features,
// but has a single backing array
type DynamicMesh struct {
	// Backing array is for Dynamic VBOs
	vertices []float32

	// Indices are for EBO
	indices []uint32

	vbo *VBO
	ebo *EBO
}

func newDynamicMesh() api.IMesh {
	m := new(DynamicMesh)
	m.vertices = []float32{}
	m.indices = []uint32{}
	m.vbo = NewVBO(api.MeshDynamic)
	m.ebo = NewEBO()
	return m
}

// VboID return internal vbo id
func (m *DynamicMesh) VboID() uint32 {
	return m.vbo.VboID()
}

// BindVBO calls BindBuffer and BufferData
func (m *DynamicMesh) BindVBO() {
	m.vbo.Bind(m)
}

// BindEBO calls BindBuffer and BufferData
func (m *DynamicMesh) BindEBO() {
	m.ebo.Bind(m)
}

// GenNextBackingIndex generates next backing array index
func (m *DynamicMesh) GenNextBackingIndex() int {
	return 0
}

// AddArray creates a new backing array
func (m *DynamicMesh) AddArray() int {
	return 0
}

// ActivateArray activates a previous backing array
func (m *DynamicMesh) ActivateArray(backingIdx int) {
}

// AddVertex adds a vertex
func (m *DynamicMesh) AddVertex(x, y, z float32) {
	m.vertices = append(m.vertices, x, y, z)
}

// SetVertex sets an existing vertex components.
func (m *DynamicMesh) SetVertex(x, y, z float32, index int) {
	i := index * 3
	m.vertices[i] = x
	m.vertices[i+1] = y
	m.vertices[i+2] = z
}

// AddIndex adds an index for EBOs
func (m *DynamicMesh) AddIndex(index int) {
	m.indices = append(m.indices, uint32(index))
}

// Vertices returns the internal vertices
func (m *DynamicMesh) Vertices() []float32 {
	return m.vertices
}

// VerticesUsing returns a selected vertex array
// FIXME refactor to remove this unneeded method
func (m *DynamicMesh) VerticesUsing(backingArrayIdx int) []float32 {
	return m.vertices
}

// Indices returns the internal indices
func (m *DynamicMesh) Indices() []uint32 {
	return m.indices
}

// Update modifies the VBO buffer
func (m *DynamicMesh) Update(offset, size int) {
	m.vbo.Update(offset, size, m.vertices)
}

// UpdatePreScaled requires prescaled values
func (m *DynamicMesh) UpdatePreScaled(offset, size int) {
	m.vbo.UpdatePreScaled(offset, size, m.vertices)
}

// Discard deletes the backing array
func (m *DynamicMesh) Discard() {
	m.vertices = nil
}
