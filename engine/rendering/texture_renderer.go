package rendering

import (
	"image"
	"log"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

type textureRenderer struct {
	vao, tbo, vbo, ebo uint32

	shader     api.IShader
	textureMan api.ITextureManager
	atlasIndex int

	modelLoc, colorLoc int32

	quad    []float32
	indices []uint32
}

// NewTextureRenderer creates a specific texture renderer
func NewTextureRenderer(textureMan api.ITextureManager, shader api.IShader) api.ITextureRenderer {
	o := new(textureRenderer)

	o.textureMan = textureMan
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

func (t *textureRenderer) Build(atlasName string) {
	var textureAtlas api.ITextureAtlas

	t.atlasIndex, textureAtlas = t.textureMan.GetAtlasPairByName(atlasName)

	gl.GenVertexArrays(1, &t.vao)

	gl.GenBuffers(1, &t.vbo)

	// Activate VBO buffer while in the VAOs scope
	gl.BindVertexArray(t.vao)

	// Indices defined in CCW order
	t.indices = []uint32{
		// 0, 1, 3, // first triangle
		// 1, 2, 3, // second triangle
		// OR
		0, 1, 2, // first triangle
		0, 2, 3, // second triangle
	}

	coords := t.textureMan.GetSTCoords(t.atlasIndex, 0)
	if coords == nil {
		panic("Sub texture not found")
	}
	c := *coords

	t.quad = append(t.quad, -0.5, -0.5, 0.0) // xy = aPos
	t.quad = append(t.quad, c[0], c[1])      // uv = aTexCoord
	t.quad = append(t.quad, 0.5, -0.5, 0.0)
	t.quad = append(t.quad, c[2], c[3]) // uv
	t.quad = append(t.quad, 0.5, 0.5, 0.0)
	t.quad = append(t.quad, c[4], c[5]) // uv
	t.quad = append(t.quad, -0.5, 0.5, 0.0)
	t.quad = append(t.quad, c[6], c[7]) // uv

	t.bindTextureVbo()

	// Activate EBO buffer while in the VAOs scope
	gl.GenBuffers(1, &t.ebo)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, t.ebo)
	sizeOfUInt32 := int(unsafe.Sizeof(uint32(0)))
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizeOfUInt32*len(t.indices), gl.Ptr(t.indices), gl.STATIC_DRAW)

	if errNum := gl.GetError(); errNum != gl.NO_ERROR {
		log.Fatal("(ebo)GL Error: ", errNum)
	}

	gl.GenTextures(1, &t.tbo)

	t.bindTbo(textureAtlas.AtlasImage())

	gl.BindVertexArray(0) // close scope
	// --------- Scope capturing ENDs here -------------------
}

func (t *textureRenderer) Use() {
	t.shader.Use()

	gl.BindVertexArray(t.vao)

	gl.ActiveTexture(gl.TEXTURE0)

	gl.BindTexture(gl.TEXTURE_2D, t.tbo)
}

func (t *textureRenderer) UnUse() {
	gl.BindVertexArray(0)
}

func (t *textureRenderer) Draw(model api.IMatrix4) {
	gl.UniformMatrix4fv(t.modelLoc, 1, false, &model.Matrix()[0])
	gl.DrawElements(gl.TRIANGLES, int32(len(t.indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
}

// SetColor sets the mix color on texture.
func (t *textureRenderer) SetColor(color []float32) {
	gl.Uniform4fv(t.colorLoc, 1, &color[0])
}

func (t *textureRenderer) SelectCoordsByIndex(index int) {
	coords := t.textureMan.GetSTCoords(t.atlasIndex, index)
	c := *coords

	i := 3
	t.quad[i] = c[0]
	t.quad[i+1] = c[1]
	i += 3 + 2
	t.quad[i] = c[2]
	t.quad[i+1] = c[3]
	i += 3 + 2
	t.quad[i] = c[4]
	t.quad[i+1] = c[5]
	i += 3 + 2
	t.quad[i] = c[6]
	t.quad[i+1] = c[7]

	// Update moves any modified data to the buffer.
	gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(t.quad)*4, gl.Ptr(t.quad))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (t *textureRenderer) bindTextureVbo() {
	gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(t.quad), gl.Ptr(t.quad), gl.DYNAMIC_DRAW)

	sizeOfFloat := int32(4)

	// If the data per-vertex is (x,y,z,s,t = 5) then
	// Stride = 5 * size of float
	// OR
	// If the data per-vertex is (x,y,z,r,g,b,s,t = 8) then
	// Stride = 8 * size of float

	// Our data layout is x,y,z,s,t
	stride := 5 * sizeOfFloat

	// position attribute
	size := int32(3)   // x,y,z
	offset := int32(0) // position is first thus this attrib is offset by 0
	attribIndex := uint32(0)
	gl.VertexAttribPointer(attribIndex, size, gl.FLOAT, false, stride, gl.PtrOffset(int(offset)))
	gl.EnableVertexAttribArray(0)

	// texture coord attribute is offset by 3 (i.e. x,y,z)
	size = int32(2)   // s,t
	offset = int32(3) // the preceeding component size = 3, thus this attrib is offset by 3
	attribIndex = uint32(1)
	gl.VertexAttribPointer(attribIndex, size, gl.FLOAT, false, stride, gl.PtrOffset(int(offset*sizeOfFloat)))
	gl.EnableVertexAttribArray(1)
}

func (t *textureRenderer) bindTbo(texture *image.NRGBA) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.tbo)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	width := int32(texture.Bounds().Dx())
	height := int32(texture.Bounds().Dy())

	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	// Give the image to OpenGL
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(texture.Pix))
	// gl.GenerateMipmap(gl.TEXTURE_2D)
}
