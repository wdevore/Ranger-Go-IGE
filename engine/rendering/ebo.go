package rendering

import (
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
)

// EBO represents a shader's EBO features.
type EBO struct {
	// Indicate if an Id has been generated yet.
	genBound bool

	eboID uint32 // GLuint
}

// NewEBO creates a empty EBO
func NewEBO() *EBO {
	b := new(EBO)
	b.genBound = false
	return b
}

// GenBuffer generates a buffer id for buffer data -
// Call this BEFORE you call Bind.
func (b *EBO) GenBuffer() {
	gl.GenBuffers(1, &b.eboID)
	b.genBound = true
}

// Bind binds the buffer id against the mesh indices
func (b *EBO) Bind(m *Mesh) {
	if !b.genBound {
		panic("An EBO buffer ID has not been generated. Call GenBuffer first.")
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.eboID)
	indicesCount := len(m.Indices) * int(unsafe.Sizeof(uint32(0)))
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indicesCount, gl.Ptr(m.Indices), gl.STATIC_DRAW)
}
