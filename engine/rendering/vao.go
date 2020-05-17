package rendering

import (
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// We only have one attribute in this engine
const attributeIndex uint32 = 0

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
	arrayCount := int32(api.XYZComponentCount) * int32(unsafe.Sizeof(float32(0)))
	gl.VertexAttribPointer(attributeIndex, int32(api.XYZComponentCount), gl.FLOAT, false, arrayCount, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0) // UnUse
}

// Render shape using VAO
func (v *VAO) Render(vs api.IAtlasShape) {
	// The signature of glDrawElements was defined back before there were buffer objects;
	// originally you'd pass an actual pointer to data in a client-side vertex array.
	// When device-side buffers were introduced, this function was extended to support them
	// as well, by shoehorning a buffer offset into the address argument.
	// Because we are using VBOs we need to awkwardly cast the offset value into a
	// pointer to void.
	// If we weren't using VBOs then we would use client-side addresses: &_mesh->indices[offset]
	gl.DrawElements(vs.PrimitiveMode(), int32(vs.Count()), uint32(gl.UNSIGNED_INT), gl.PtrOffset(vs.Offset()))
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

// void VAO::draw(int primitiveType, int offset, int count) {
// 	glBindVertexArray(_vaoId);
// 	glDrawElements(primitiveType, count, GL_UNSIGNED_INT, (const GLvoid*)(offset));
// 	glBindVertexArray(0);
// }
