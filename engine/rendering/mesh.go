// Package rendering defines Mesh features of shaders.
package rendering

// Mesh combines a shader's VBO and EBO features.
type Mesh struct {
	// Vertices are for VBO
	Vertices []float32

	// Indices are for EBO
	Indices []uint32

	vbo VBO
	ebo EBO
}

// NewMesh creates a new Mesh object
func NewMesh() *Mesh {
	m := new(Mesh)
	m.Vertices = make([]float32, 100)
	m.Indices = make([]uint32, 200)
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
