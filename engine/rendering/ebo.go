package rendering

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

// EBO represents a shader's EBO features.
type EBO struct {
	eboID uint32 // GLuint
}

// NewEBO creates a empty EBO
func NewEBO() *EBO {
	b := new(EBO)
	gl.GenBuffers(1, &b.eboID)
	return b
}

// Bind binds the buffer id against the mesh indices
func (b *EBO) Bind(bufferSize int, indices []uint32) {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.eboID)

	// indicesSize := len(indices) * int(unsafe.Sizeof(uint32(0)))
	// fmt.Println("EBO.Bind: ", indicesSize)
	// gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indicesSize, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bufferSize, gl.Ptr(indices), gl.STATIC_DRAW)
}
