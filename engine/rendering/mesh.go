// Package rendering defines Mesh features of shaders.
package rendering

import "github.com/wdevore/Ranger-Go-IGE/api"

// Mesh combines a shader's VBO and EBO features.
type Mesh struct {
	// Vertices are for VBO
	vertices []float32

	// Indices are for EBO
	indices []uint32

	vbo *VBO
	ebo *EBO
}

// NewMesh creates a new Mesh object
func NewMesh() api.IMesh {
	m := new(Mesh)
	m.vertices = []float32{}
	m.indices = []uint32{}
	m.vbo = NewVBO()
	m.ebo = NewEBO()
	return m
}

// Bind binds this Mesh to a VBO and EBO
func (m *Mesh) Bind() {
	m.vbo.Bind(m)
	m.ebo.Bind(m)
}

// GenBuffers generates buffers for VBO and EBO
func (m *Mesh) GenBuffers() {
	m.vbo.GenBuffer()
	m.ebo.GenBuffer()
}

func (m *Mesh) AddVertex(x, y, z float32) {
	m.vertices = append(m.vertices, x, y, z)
}

func (m *Mesh) SetVertex(x, y, z float32, index int) {
	m.vertices[index] = x
	m.vertices[index+1] = y
	m.vertices[index+2] = z
}

func (m *Mesh) AddIndex(index int) {
	m.indices = append(m.indices, uint32(index))
}

func (m *Mesh) Vertices() []float32 {

	return m.vertices
}

func (m *Mesh) Indices() []uint32 {
	return m.indices
}
