// Package rendering defines VBO features of shaders.
package rendering

import (
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
)

// VBO represents a shader's VBO features.
type VBO struct {
	vboID uint32 // GLuint

	staticDraw  bool
	bufferUsage uint32
	floatSize   int
}

// NewVBO creates a empty VBO defaulting to STATIC_DRAW
func NewVBO() *VBO {
	o := new(VBO)

	gl.GenBuffers(1, &o.vboID)
	o.staticDraw = true
	o.floatSize = int(unsafe.Sizeof(float32(0)))

	return o
}

// SetDrawUsage changes buffer usage style between static or dynamic.
// This MUST be called prior to Bind() call.
func (v *VBO) SetDrawUsage(usage bool) {
	v.staticDraw = usage
}

// Bind binds the buffer id against the mesh vertices
func (v *VBO) Bind(m *Mesh) {
	// the buffer type of a vertex buffer object is GL_ARRAY_BUFFER
	// From this point on any buffer calls we make (on the GL_ARRAY_BUFFER target)
	// will be used to configure the currently bound buffer, which is VBO
	gl.BindBuffer(gl.ARRAY_BUFFER, v.vboID)

	// Then we can make a call to the glBufferData function that copies the previously
	// defined vertex data into the buffer's memory. glBufferData is a function
	// specifically targeted to copy user-defined data into the currently bound buffer.
	v.bufferUsage = uint32(gl.STATIC_DRAW)
	if !v.staticDraw {
		v.bufferUsage = gl.DYNAMIC_DRAW
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices())*v.floatSize, gl.Ptr(m.Vertices()), v.bufferUsage)
}

// Update moves any modified vertices to the buffer.
// Note: Be sure to call the AtlasObject SetVertex prior.
func (v *VBO) Update(offset, vertexCount int, m *Mesh) {
	if v.bufferUsage == gl.STATIC_DRAW {
		panic("VBO is not configured as DYNAMIC_DRAW.")
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, v.vboID)
	gl.BufferSubData(gl.ARRAY_BUFFER, offset, vertexCount*v.floatSize, gl.Ptr(m.Vertices()))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
