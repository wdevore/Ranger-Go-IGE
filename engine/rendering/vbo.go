// Package rendering defines VBO features of shaders.
package rendering

import (
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// VBO represents a shader's VBO features.
type VBO struct {
	vboID uint32 // GLuint

	staticDraw  int
	bufferUsage uint32
	floatSize   int
}

// NewVBO creates a empty VBO defaulting to STATIC_DRAW
func NewVBO(meshType int) *VBO {
	o := new(VBO)

	gl.GenBuffers(1, &o.vboID)
	o.staticDraw = meshType
	o.floatSize = int(unsafe.Sizeof(float32(0)))

	return o
}

// VboID returns internal vbo buffer id.
func (v *VBO) VboID() uint32 {
	return v.vboID
}

// Use bind vertex array to Id
func (v *VBO) Use() {
}

// SetDrawUsage changes buffer usage style between static or dynamic.
// This MUST be called prior to Bind() call.
func (v *VBO) SetDrawUsage(meshType int) {
	v.staticDraw = meshType
}

// Bind binds the buffer id against the mesh vertices
func (v *VBO) Bind(bufferSize int, vertices []float32) {
	// the buffer type of a vertex buffer object is GL_ARRAY_BUFFER
	// From this point on any buffer calls we make (on the GL_ARRAY_BUFFER target)
	// will be used to configure the currently bound buffer, which is VBO
	gl.BindBuffer(gl.ARRAY_BUFFER, v.vboID)

	// Then we can make a call to the glBufferData function that copies the previously
	// defined vertex data into the buffer's memory. glBufferData is a function
	// specifically targeted to copy user-defined data into the currently bound buffer.
	v.bufferUsage = uint32(gl.STATIC_DRAW)
	if v.staticDraw != api.MeshStatic {
		v.bufferUsage = gl.DYNAMIC_DRAW
	}

	// Create and fill buffer
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize, gl.Ptr(vertices), v.bufferUsage)

	// Static vbo's don't need the backing array once the data has
	// been copied to the GL buffer.
	// TODO we may not want to discard the backing data incase this is
	// needed by other code. For now it is disabled
	// if v.bufferUsage == gl.STATIC_DRAW {
	// 	m.Discard()
	// }
}

// Update moves any modified vertices to the buffer.
// Note: Be sure to call the AtlasObject SetVertex prior.
func (v *VBO) Update(offset, count int, vertices []float32) {
	if v.bufferUsage == gl.STATIC_DRAW {
		panic("VBO is not configured as DYNAMIC_DRAW.")
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, v.vboID)
	gl.BufferSubData(gl.ARRAY_BUFFER, offset*v.floatSize, count*v.floatSize*api.XYZComponentCount, gl.Ptr(vertices))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
