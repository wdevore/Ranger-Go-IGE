// Package rendering defines Mesh features of shaders.
package rendering

import "github.com/wdevore/Ranger-Go-IGE/api"

// DynamicMultiMesh uses multiple backing arrays, one for each objec type.
type DynamicMultiMesh struct {
	// Backing array is for Dynamic VBOs
	// The total size of all arrays is the size the VBO
	// buffer is defined for.
	vertices              [][]float32
	backingIndexCnt       int
	activeBackingArrayIdx int

	// Indices are for EBO
	indices []uint32

	vbo *VBO
	ebo *EBO
}

func newDynamicMultiMesh() api.IMesh {
	m := new(DynamicMultiMesh)
	m.vertices = [][]float32{}
	m.indices = []uint32{}
	m.vbo = NewVBO(api.MeshDynamicMulti)
	m.ebo = NewEBO()
	return m
}

// VboID return internal vbo id
func (m *DynamicMultiMesh) VboID() uint32 {
	return m.vbo.VboID()
}

// BindVBO calls BindBuffer and BufferData
func (m *DynamicMultiMesh) BindVBO() {
	m.vbo.Bind(m)
}

// BindEBO calls BindBuffer and BufferData
func (m *DynamicMultiMesh) BindEBO() {
	m.ebo.Bind(m)
}

// GenNextBackingIndex generates next backing array index
func (m *DynamicMultiMesh) GenNextBackingIndex() int {
	idx := m.backingIndexCnt
	m.backingIndexCnt++
	return idx
}

// AddArray creates a new backing array
func (m *DynamicMultiMesh) AddArray() int {
	m.activeBackingArrayIdx = m.GenNextBackingIndex()

	m.vertices = append(m.vertices, []float32{})

	return m.activeBackingArrayIdx
}

// ActivateArray activates a previous backing array
func (m *DynamicMultiMesh) ActivateArray(backingIdx int) {
	m.activeBackingArrayIdx = backingIdx
}

// AddVertex adds a vertex
func (m *DynamicMultiMesh) AddVertex(x, y, z float32) {
	v := m.vertices[m.activeBackingArrayIdx]
	v = append(v, x, y, z)
	m.vertices[m.activeBackingArrayIdx] = v
}

// SetVertex sets an existing vertex components.
func (m *DynamicMultiMesh) SetVertex(x, y, z float32, index int) {
	i := index * 3
	v := m.vertices[m.activeBackingArrayIdx]
	v[i] = x
	v[i+1] = y
	v[i+2] = z
}

// AddIndex adds an index for EBOs
func (m *DynamicMultiMesh) AddIndex(index int) {
	m.indices = append(m.indices, uint32(index))
}

// Vertices returns the internal vertices
func (m *DynamicMultiMesh) Vertices() []float32 {
	v := m.vertices[m.activeBackingArrayIdx]
	return v
}

// VerticesUsing returns a selected vertex array
func (m *DynamicMultiMesh) VerticesUsing(backingArrayIdx int) []float32 {
	v := m.vertices[backingArrayIdx]
	return v
}

// Indices returns the internal indices
func (m *DynamicMultiMesh) Indices() []uint32 {
	return m.indices
}

// Update modifies the VBO buffer
func (m *DynamicMultiMesh) Update(offset, size int) {
	v := m.vertices[m.activeBackingArrayIdx]
	m.vbo.Update(offset, size, v)
}

// UpdatePreScaled requires prescaled values
func (m *DynamicMultiMesh) UpdatePreScaled(offset, size int) {
	v := m.vertices[m.activeBackingArrayIdx]
	m.vbo.UpdatePreScaled(offset, size, v)
}

// Discard deletes the backing array
func (m *DynamicMultiMesh) Discard() {
	m.vertices = nil
}
