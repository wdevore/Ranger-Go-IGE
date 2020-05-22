// Package rendering defines VBO features of shaders.
package rendering

import (
	"log"
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
func NewVBO(isStatic bool) *VBO {
	o := new(VBO)

	gl.GenBuffers(1, &o.vboID)
	o.staticDraw = isStatic
	o.floatSize = int(unsafe.Sizeof(float32(0)))

	return o
}

func (v *VBO) VboID() uint32 {
	return v.vboID
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
		// gl.BufferStorage()
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices())*v.floatSize, gl.Ptr(m.Vertices()), v.bufferUsage)
}

// Update moves any modified vertices to the buffer.
// Note: Be sure to call the AtlasObject SetVertex prior.
func (v *VBO) Update(offset, count int, vertices []float32) {
	if v.bufferUsage == gl.STATIC_DRAW {
		panic("VBO is not configured as DYNAMIC_DRAW.")
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, v.vboID)
	gl.BufferSubData(gl.ARRAY_BUFFER, offset*v.floatSize, count*v.floatSize, gl.Ptr(vertices))
	// gl.BufferSubData(gl.ARRAY_BUFFER, offset, count, gl.Ptr(m.Vertices()))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

// UpdatePreScaled expects parms to be prescaled by data-type.
func (v *VBO) UpdatePreScaled(offset, size int, vertices []float32) {
	// if v.bufferUsage == gl.STATIC_DRAW {
	// 	panic("VBO is not configured as DYNAMIC_DRAW.")
	// }

	gl.BindBuffer(gl.ARRAY_BUFFER, v.vboID)

	// The last parameter should be a separate buffer
	// The 'offset' and 'size' are parameters for the destination
	// buffer and in is bytes
	// In other words the source buffer is captured as a whole,
	// the destination buffer is piece-meal based on offset/size
	gl.BufferSubData(gl.ARRAY_BUFFER, offset, size, gl.Ptr(vertices))

	if errNum := gl.GetError(); errNum != gl.NO_ERROR {
		log.Fatal("(vbo)GL Error: ", errNum)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
