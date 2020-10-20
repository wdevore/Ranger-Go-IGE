package atlas

import (
	"errors"
	"image"
	"path/filepath"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

type singleTextureAtlas struct {
	world api.IWorld
	burnt bool

	atlasName string

	vao, tbo, vbo, ebo uint32

	// For the Shaking process
	quad    []float32
	indices []uint32

	shader      api.IShader
	spriteSheet api.ISpriteSheet

	modelLoc, colorLoc int32
}

// ###################################################################
// TODO improve this. We need to create a buffer large enough for a few
// hundred characters so we would only need to copy the buffer once instead
// for every character in a string.
// ###################################################################

// NewSingleTextureAtlas creates a specific texture atlas that renders to a single quad
// This quad could represent a single character is a text string.
func NewSingleTextureAtlas(atlasName string, spriteSheet api.ISpriteSheet, world api.IWorld) api.IAtlasX {
	o := new(singleTextureAtlas)

	o.spriteSheet = spriteSheet
	o.world = world
	o.atlasName = atlasName

	return o
}

func (t *singleTextureAtlas) Burnt() bool {
	return t.burnt
}

func (t *singleTextureAtlas) Burn() error {
	err := t.Configure()
	if err != nil {
		return err
	}

	err = t.Bake()
	if err != nil {
		return err
	}

	t.burnt = true

	return nil
}

func (t *singleTextureAtlas) Configure() error {
	err := t.configureShaders(t.world)
	if err != nil {
		return err
	}

	return nil
}

func (t *singleTextureAtlas) SpriteSheet() api.ISpriteSheet {
	return t.spriteSheet
}

func (t *singleTextureAtlas) Shake() {
	// Indices defined in CCW order
	t.indices = []uint32{
		// 0, 1, 3, // first triangle
		// 1, 2, 3, // second triangle
		// OR
		0, 1, 2, // first triangle
		0, 2, 3, // second triangle
	}

	coords := t.spriteSheet.TextureSTCoordsByIndex(0)
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
}

func (t *singleTextureAtlas) Bake() error {
	gl.GenVertexArrays(1, &t.vao)

	gl.GenBuffers(1, &t.vbo)

	gl.GenBuffers(1, &t.ebo)

	// Activate VBO buffer while in the VAOs scope
	gl.BindVertexArray(t.vao)

	t.Shake()

	t.bindVbo()

	// Activate EBO buffer while in the VAOs scope
	t.bindEbo()

	gl.GenTextures(1, &t.tbo)

	t.bindTbo(t.spriteSheet.SheetImage())

	gl.BindVertexArray(0) // close scope
	// --------- Scope capturing ENDs here -------------------

	err := t.configureUniforms()
	if err != nil {
		return err
	}

	return nil
}

func (t *singleTextureAtlas) Use() {
	t.shader.Use()

	gl.BindVertexArray(t.vao)

	gl.ActiveTexture(gl.TEXTURE0)

	gl.BindTexture(gl.TEXTURE_2D, t.tbo)
}

func (t *singleTextureAtlas) UnUse() {
	gl.BindVertexArray(0)
}

func (t *singleTextureAtlas) Render(id int, model api.IMatrix4) {
	gl.UniformMatrix4fv(t.modelLoc, 1, false, &model.Matrix()[0])
	gl.DrawElements(gl.TRIANGLES, int32(len(t.indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
}

// SetColor sets the mix color on texture.
func (t *singleTextureAtlas) SetColor(color []float32) {
	gl.Uniform4fv(t.colorLoc, 1, &color[0])
}

// SelectCoordsByIndex implements: ISingleTextureAtlasX
func (t *singleTextureAtlas) SelectCoordsByIndex(index int) {
	coords := t.spriteSheet.TextureSTCoordsByIndex(index)
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

func (t *singleTextureAtlas) configureShaders(world api.IWorld) error {
	dataPath, err := filepath.Abs(world.RelativePath())
	if err != nil {
		return err
	}

	shaders := world.Properties().Shaders

	t.shader = rendering.NewShader(shaders.TextureVertexShaderFile, shaders.TextureFragmentShaderFile)
	err = t.shader.Load(dataPath)
	if err != nil {
		return err
	}

	return nil
}

func (t *singleTextureAtlas) configureUniforms() error {
	t.shader.Use()

	program := t.shader.Program()

	t.modelLoc = gl.GetUniformLocation(program, gl.Str("model\x00"))
	if t.modelLoc < 0 {
		return errors.New("Couldn't find 'model' uniform variable")
	}

	t.colorLoc = gl.GetUniformLocation(program, gl.Str("fragColor\x00"))
	if t.colorLoc < 0 {
		return errors.New("Couldn't find 'fragColor' uniform variable")
	}

	// Projection and View
	projLoc := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	if projLoc < 0 {
		return errors.New("SingleTextureAtlas: couldn't find 'projection' uniform variable")
	}

	viewLoc := gl.GetUniformLocation(program, gl.Str("view\x00"))
	if viewLoc < 0 {
		return errors.New("SingleTextureAtlas: couldn't find 'view' uniform variable")
	}

	pm := t.world.Projection().Matrix()
	gl.UniformMatrix4fv(projLoc, 1, false, &pm[0])

	vm := t.world.Viewspace().Matrix()
	gl.UniformMatrix4fv(viewLoc, 1, false, &vm[0])

	return nil
}

func (t *singleTextureAtlas) bindVbo() {
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

func (t *singleTextureAtlas) bindTbo(texture *image.NRGBA) {
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

func (t *singleTextureAtlas) bindEbo() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, t.ebo)

	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, uintSize*len(t.indices), gl.Ptr(t.indices), gl.STATIC_DRAW)

	// if errNum := gl.GetError(); errNum != gl.NO_ERROR {
	// 	return errors.New("(ebo)GL Error: ", errNum)
	// }
}
