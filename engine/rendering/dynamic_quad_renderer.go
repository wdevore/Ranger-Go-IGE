package rendering

import (
	"log"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

type dynamicQuadRenderer struct {
	vao, vbo, ebo uint32

	shader     api.IShader
	textureMan api.ITextureManager
	atlasIndex int

	modelLoc, colorLoc int32

	vertices []float32
	indices  []uint32
}

// NewDynamicQuadRenderer creates a dynamic rectangle quad renderer.
// It is designed so the corners can be updated dynamically.
func NewDynamicQuadRenderer(shader api.IShader) api.IDynamicRenderer {
	o := new(dynamicQuadRenderer)

	o.shader = shader

	programID := shader.Program()
	shader.Use()

	o.modelLoc = gl.GetUniformLocation(programID, gl.Str("model\x00"))
	if o.modelLoc < 0 {
		panic("Couldn't find 'model' uniform variable")
	}

	o.colorLoc = gl.GetUniformLocation(programID, gl.Str("fragColor\x00"))
	if o.colorLoc < 0 {
		panic("Couldn't find 'fragColor' uniform variable")
	}

	return o
}

func (t *dynamicQuadRenderer) Build(atlasName string) {
	gl.GenVertexArrays(1, &t.vao)

	gl.GenBuffers(1, &t.vbo)

	// Activate VBO buffer while in the VAOs scope
	gl.BindVertexArray(t.vao)

	t.vertices = []float32{
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
	}

	// Indices defined in CCW order
	t.indices = []uint32{
		// 0, 1, 3, // first triangle
		// 1, 2, 3, // second triangle
		// OR
		0, 1, 2, // first triangle
		0, 2, 3, // second triangle
	}

	// Activate EBO buffer while in the VAOs scope
	gl.GenBuffers(1, &t.ebo)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, t.ebo)
	sizeOfUInt32 := int(unsafe.Sizeof(uint32(0)))
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizeOfUInt32*len(t.indices), gl.Ptr(t.indices), gl.STATIC_DRAW)

	if errNum := gl.GetError(); errNum != gl.NO_ERROR {
		log.Fatal("(ebo)GL Error: ", errNum)
	}

	gl.BindVertexArray(0) // close scope
	// --------- Scope capturing ENDs here -------------------
}

func (t *dynamicQuadRenderer) Use() {
	t.shader.Use()

	gl.BindVertexArray(t.vao)
}

func (t *dynamicQuadRenderer) UnUse() {
	gl.BindVertexArray(0)
}

func (t *dynamicQuadRenderer) Draw(model api.IMatrix4) {
	gl.UniformMatrix4fv(t.modelLoc, 1, false, &model.Matrix()[0])
	gl.DrawElements(gl.TRIANGLES, int32(len(t.indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
}

// SetColor sets the mix color on texture.
func (t *dynamicQuadRenderer) SetColor(color []float32) {
	gl.Uniform4fv(t.colorLoc, 1, &color[0])
}

func (t *dynamicQuadRenderer) Update() {
	// Update moves any modified data to the buffer.
	gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(t.vertices)*4, gl.Ptr(t.vertices))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (t *dynamicQuadRenderer) bindVbo() {
	gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(t.vertices), gl.Ptr(t.vertices), gl.DYNAMIC_DRAW)

	sizeOfFloat := int32(4)

	// If the data per-vertex is (x,y,z,s,t = 5) then
	// Stride = 5 * size of float

	// Our data layout is (x,y,z = 3)
	stride := 3 * sizeOfFloat

	// position attribute
	size := int32(3)   // x,y,z
	offset := int32(0) // position is first thus this attrib is offset by 0
	attribIndex := uint32(0)
	gl.VertexAttribPointer(attribIndex, size, gl.FLOAT, false, stride, gl.PtrOffset(int(offset)))
	gl.EnableVertexAttribArray(0)
}
