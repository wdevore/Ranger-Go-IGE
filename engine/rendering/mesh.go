// Package rendering defines Mesh features of shaders.
package rendering

import "github.com/wdevore/Ranger-Go-IGE/api"

// Mesh combines a shader's VBO and EBO features.
type Mesh struct {
	// Backing array is for Dynamic VBOs
	// The total size of all arrays is the size the VBO
	// buffer is defined for.
	vertices []float32

	// Indices are for EBO
	indices []uint32

	vbo *VBO
	ebo *EBO
}

// NewMesh creates a new Mesh object
func NewMesh(isStatic bool) api.IMesh {
	m := new(Mesh)
	m.vertices = []float32{}
	m.indices = []uint32{}
	m.vbo = NewVBO(isStatic)
	m.ebo = NewEBO()
	return m
}

// VboID return internal vbo id
func (m *Mesh) VboID() uint32 {
	return m.vbo.VboID()
}

// BindVBO calls BindBuffer and BufferData
func (m *Mesh) BindVBO() {
	m.vbo.Bind(m)
}

// BindEBO calls BindBuffer and BufferData
func (m *Mesh) BindEBO() {
	m.ebo.Bind(m)
}

// AddVertex adds a vertex
func (m *Mesh) AddVertex(x, y, z float32) {
	m.vertices = append(m.vertices, x, y, z)
}

// SetVertex sets an existing vertex components.
func (m *Mesh) SetVertex(x, y, z float32, index int) {
	i := index * 3
	m.vertices[i] = x
	m.vertices[i+1] = y
	m.vertices[i+2] = z
}

// AddIndex adds an index for EBOs
func (m *Mesh) AddIndex(index int) {
	m.indices = append(m.indices, uint32(index))
}

// Vertices returns the internal vertices
func (m *Mesh) Vertices() []float32 {
	return m.vertices
}

// Indices returns the internal indices
func (m *Mesh) Indices() []uint32 {
	return m.indices
}

// Update modifies the VBO buffer
func (m *Mesh) Update(offset, size int) {
	m.vbo.Update(offset, size, m.vertices)
}

// UpdatePreScaled requires prescaled values
func (m *Mesh) UpdatePreScaled(offset, size int) {
	m.vbo.UpdatePreScaled(offset, size, m.vertices)
}

// Discard deletes the backing array
func (m *Mesh) Discard() {
	m.vertices = nil
}
