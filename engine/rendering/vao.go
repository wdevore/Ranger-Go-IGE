package rendering

import (
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

const xyzComponentCount int32 = 3
const attributeIndex uint32 = 0

// VAO defines a Vertex Array Object
type VAO struct {
	// Indicates if an Id has been generated
	genBound bool
	vaoID    uint32
	mesh     api.IMesh
}

// NewVAO creates a new VAO
func NewVAO(m api.IMesh) *VAO {
	v := new(VAO)
	v.mesh = m
	return v
}

// Bind setups the VAO and Mesh
func (v *VAO) Bind() {
	if !v.genBound {
		gl.GenVertexArrays(1, &v.vaoID)
		v.genBound = true
	}

	v.mesh.GenBuffers()

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s)
	// and attribute pointer(s).
	gl.BindVertexArray(v.vaoID)

	v.mesh.Bind()

	arrayCount := xyzComponentCount * int32(unsafe.Sizeof(float32(0)))
	gl.VertexAttribPointer(attributeIndex, int32(xyzComponentCount), gl.FLOAT, false, arrayCount, gl.PtrOffset(0))

	// vertexInputAttrb := uint32(0) // aPos
	// sizeOfInputAttrb := int32(3)  // 3 floats. The 4th is supplied by the shader
	// stride := int32(0)            // All the vertices will tightly packed
	// gl.VertexAttribPointer(vertexInputAttrb, sizeOfInputAttrb, gl.FLOAT, false, stride, nil)
	// gl.EnableVertexAttribArray(vertexInputAttrb)

	gl.EnableVertexAttribArray(0)

	// Note that this is allowed, the call to glVertexAttribPointer registered VBO as the currently bound
	// vertex buffer object so afterwards we can safely unbind
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// Unbind VAO (it's always a good thing to unbind any buffer/array to prevent strange bugs),
	// remember: do NOT unbind the EBO, keep it bound to this VAO
	gl.BindVertexArray(0)

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

	// Rather than multiply repeatedly
	//glDrawElements(_shape->primitiveType, _shape->count, GL_UNSIGNED_INT, (const GLvoid*)(_shape->offset * sizeof(unsigned int)));
	// we use a pre computed version.
	gl.DrawElements(vs.PrimitiveMode(), vs.Count(), uint32(gl.UNSIGNED_INT), gl.PtrOffset(vs.Offset()))
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
