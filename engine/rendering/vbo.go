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
	o := new(VBO)
	o.genBound = false
	return o
}

// GenBuffer generates a buffer id for buffer data -
// Call this BEFORE you call Bind.
func (v *VBO) GenBuffer() {
	if !v.genBound {
		gl.GenBuffers(1, &v.vboID)
		v.genBound = true
	}
}

// Bind binds the buffer id against the mesh vertices
func (v *VBO) Bind(m *Mesh) {
	if !v.genBound {
		panic("A VBO buffer ID has not been generated. Call GenBuffer first.")
	}

	// the buffer type of a vertex buffer object is GL_ARRAY_BUFFER
	// From this point on any buffer calls we make (on the GL_ARRAY_BUFFER target)
	// will be used to configure the currently bound buffer, which is VBO
	gl.BindBuffer(gl.ARRAY_BUFFER, v.vboID)
	floatSize := int(unsafe.Sizeof(float32(0)))

	// Then we can make a call to the glBufferData function that copies the previously
	// defined vertex data into the buffer's memory. glBufferData is a function
	// specifically targeted to copy user-defined data into the currently bound buffer.
	gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices())*floatSize, gl.Ptr(m.Vertices()), gl.STATIC_DRAW)
}
