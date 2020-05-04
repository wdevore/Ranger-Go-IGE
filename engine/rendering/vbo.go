// Package rendering defines VBO features of shaders.
package rendering

import (
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
)

// VBO represents a shader's VBO features.
type VBO struct {
	// Indicate if an Id has been generated yet.
	genBound bool

	vboID uint32 // GLuint
}

// NewVBO creates a empty VBO
func NewVBO() *VBO {
	b := new(VBO)
	b.genBound = false
	return b
}

// GenBuffer generates a buffer id for buffer data -
// Call this BEFORE you call Bind.
func (b *VBO) GenBuffer() {
	gl.GenBuffers(1, &b.vboID)
	b.genBound = true
}

// Bind binds the buffer id against the mesh vertices
func (b *VBO) Bind(m *Mesh) {
	if !b.genBound {
		panic("A VBO buffer ID has not been generated. Call GenBuffer first.")
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, b.vboID)
	floatSize := int(unsafe.Sizeof(float32(0)))
	gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices)*floatSize, gl.Ptr(m.Vertices), gl.STATIC_DRAW)
}
