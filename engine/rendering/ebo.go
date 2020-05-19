package rendering

import (
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
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
func (b *EBO) Bind(m api.IMesh) {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.eboID)

	indicesSize := len(m.Indices()) * int(unsafe.Sizeof(uint32(0)))
	// fmt.Println("EBO.Bind: ", indicesSize)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indicesSize, gl.Ptr(m.Indices()), gl.STATIC_DRAW)
}
