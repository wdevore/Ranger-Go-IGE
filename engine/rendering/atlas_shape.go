package rendering

import (
	"unsafe"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// AtlasShape defines shape element attributes
//
type AtlasShape struct {
	// ---------------------------------------------------------
	// VBO vars
	// ---------------------------------------------------------
	// Backing array is copied to GL buffer.
	// If the atlas is static then the backing array is copied only
	// once. Otherwise the backing array is copied continuously.
	vertices   []float32
	bufferSize int

	// ---------------------------------------------------------
	// EBO vars
	// ---------------------------------------------------------
	// Indices are for EBO. They point to entries in VBO buffer
	// As the vbo is being populated the indices are updated
	// accordingly.
	// The indices are set during GL configuration. They reference
	// vertices in the GL buffer
	indices []uint32

	name string
	// GL_LINES, GL_LINE_LOOP etc...
	primitiveMode uint32

	// Offset is multiplied by the size of an uint32 in preparation for
	// calling DrawElements()
	elementByteOffset int
	// A Line for example has an element count of 2
	elementCount int

	count int // Used for vbo updates
}

// NewAtlasShape creates a new vector shape
func NewAtlasShape() api.IAtlasShape {
	o := new(AtlasShape)
	return o
}

// SetIndices sets the indices
func (a *AtlasShape) SetIndices(indices []uint32) {
	a.indices = indices
}

// Indices returns the indices
func (a *AtlasShape) Indices() []uint32 {
	return a.indices
}

// SetBufferSize sets the expected vbo buffer size
func (a *AtlasShape) SetBufferSize(size int) {
	a.bufferSize = size
}

// BufferSize returns the expected vbo buffer size
func (a *AtlasShape) BufferSize() int {
	return a.bufferSize
}

// SetVertex2D updates a vertex in a backing array
func (a *AtlasShape) SetVertex2D(x, y float32, index int) {
	i := index * 3
	a.vertices[i] = x
	a.vertices[i+1] = y
}

// SetVertices sets the backing array
func (a *AtlasShape) SetVertices(vertices []float32) {
	a.vertices = vertices
}

// Vertices returns a selected vertex array from Mesh
func (a *AtlasShape) Vertices() *[]float32 {
	return &a.vertices
}

// SetOffset scales offset by size of an uint32
func (a *AtlasShape) SetOffset(offset int) {
	a.elementByteOffset = offset * int(unsafe.Sizeof(uint32(0)))
}

// SetElementOffset sets the EBO offset without considering data-type size
// The value should be in bytes
func (a *AtlasShape) SetElementOffset(offset int) {
	a.elementByteOffset = offset
}

// Offset returns calculated offset
func (a *AtlasShape) Offset() int {
	return a.elementByteOffset
}

// SetElementCount specifies how many elements are need to draw based
// primitive mode type.
func (a *AtlasShape) SetElementCount(count int) {
	a.elementCount = count
}

// ElementCount returns how many elements are need to draw
func (a *AtlasShape) ElementCount() int {
	return a.elementCount
}

// Name returns name
func (a *AtlasShape) Name() string {
	return a.name
}

// SetName sets name
func (a *AtlasShape) SetName(n string) {
	a.name = n
}

// PrimitiveMode returns OpenGL primative mode
func (a *AtlasShape) PrimitiveMode() uint32 {
	return a.primitiveMode
}

// SetPrimitiveMode sets name
func (a *AtlasShape) SetPrimitiveMode(m uint32) {
	a.primitiveMode = m
}

// Count returns the total count of elements
func (a *AtlasShape) Count() int {
	return a.count
}

// SetCount sets total element count
func (a *AtlasShape) SetCount(c int) {
	a.count = c
}
