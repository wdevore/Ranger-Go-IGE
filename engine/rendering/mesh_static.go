// Package rendering defines Mesh features of shaders.
package rendering

import "github.com/wdevore/Ranger-Go-IGE/api"

// StaticMesh combines a shader's VBO and EBO features.
type StaticMesh struct {
	// Backing array is for Dynamic VBOs
	// The total size of all arrays is the size the VBO
	// buffer is defined for.
	vertices []float32

	// Indices are for EBO
	indices []uint32

	vbo *VBO
	ebo *EBO
}

func newStaticMesh() api.IMesh {
	m := new(StaticMesh)
	m.vertices = []float32{}
	m.indices = []uint32{}
	m.vbo = NewVBO(api.MeshStatic)
	m.ebo = NewEBO()
	return m
}

// VboID return internal vbo id
func (m *StaticMesh) VboID() uint32 {
	return m.vbo.VboID()
}

// BindVBO calls BindBuffer and BufferData
func (m *StaticMesh) BindVBO() {
	m.vbo.Bind(m)
}

// BindEBO calls BindBuffer and BufferData
func (m *StaticMesh) BindEBO() {
	m.ebo.Bind(m)
}

// GenNextBackingIndex generates next backing array index
func (m *StaticMesh) GenNextBackingIndex() int {
	return 0
}

// AddArray creates a new backing array
func (m *StaticMesh) AddArray() int {
	return 0
}

// ActivateArray activates a previous backing array
func (m *StaticMesh) ActivateArray(backingIdx int) {
}

// AddVertex adds a vertex
func (m *StaticMesh) AddVertex(x, y, z float32) {
	m.vertices = append(m.vertices, x, y, z)
}

// SetVertex sets an existing vertex components.
func (m *StaticMesh) SetVertex(x, y, z float32, index int) {
	i := index * 3
	m.vertices[i] = x
	m.vertices[i+1] = y
	m.vertices[i+2] = z
}

// AddIndex adds an index for EBOs
func (m *StaticMesh) AddIndex(index int) {
	m.indices = append(m.indices, uint32(index))
}

// VerticesUsing returns a selected vertex array
func (m *StaticMesh) VerticesUsing(backingArrayIdx int) []float32 {
	return m.vertices
}

// Vertices returns the internal vertices
func (m *StaticMesh) Vertices() []float32 {
	return m.vertices
}

// Indices returns the internal indices
func (m *StaticMesh) Indices() []uint32 {
	return m.indices
}

// Update modifies the VBO buffer
func (m *StaticMesh) Update(offset, size int) {
	m.vbo.Update(offset, size, m.vertices)
}

// UpdatePreScaled requires prescaled values
func (m *StaticMesh) UpdatePreScaled(offset, size int) {
	m.vbo.UpdatePreScaled(offset, size, m.vertices)
}

// Discard deletes the backing array
func (m *StaticMesh) Discard() {
	m.vertices = nil
}
