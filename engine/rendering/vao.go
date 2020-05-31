package rendering

import (
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// We only have one attribute (position) in the shader
// Some shaders may have normals and/or textures coords
const positionIndex uint32 = 0

// VAO defines a Vertex Array Object
type VAO struct {
	vaoID uint32
}

// NewVAO creates a new VAO
func NewVAO() *VAO {
	o := new(VAO)
	gl.GenVertexArrays(1, &o.vaoID)
	return o
}

// BindStart setups the VAO and Mesh
func (v *VAO) BindStart() {
	// Bind the Vertex Array Object first, then bind and set vertex buffer(s)
	// and attribute pointer(s).
	gl.BindVertexArray(v.vaoID) // Use
}

// BindComplete setups vertex-array-ptr and disable vertex-array-attr
func (v *VAO) BindComplete() {
	// Count == (xyz=3) * sizeof(float32)=4 == 12 thus each
	// vertex is 12 bytes
	arrayCount := int32(api.XYZComponentCount) * int32(unsafe.Sizeof(float32(0)))

	// We can link the attribute position with the data in the vertexData array
	gl.VertexAttribPointer(positionIndex, int32(api.XYZComponentCount), gl.FLOAT, false, arrayCount, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0) // UnUse
}

// Use bind vertex array to Id
func (v *VAO) Use() {
	gl.BindVertexArray(v.vaoID)
}

// UnUse removes the array binding (optional)
func (v *VAO) UnUse() {
	// See opengl wiki as to why "glBindVertexArray(0)" isn't really necessary here:
	// https://www.opengl.org/wiki/Vertex_Specification#Vertex_Buffer_Object
	// Note the line "Changing the GL_ARRAY_BUFFER binding changes nothing about vertex attribute 0..."
	gl.BindVertexArray(0)
}
