package atlas

import (
	"errors"
	"path/filepath"

	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/configuration"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

// The atlas contains and renders shapes.
// A DyanmicPixel Atlas uses a single color for all pixels

type dynamicPixelAtlas struct {
	world api.IWorld

	// The backing array for the pixels. It is copied
	// each time the array changes.
	vertices []float32
	// This is updated with the number of pixels to rendered.
	indices []uint32

	indicesCount  int
	vboBufferSize int

	// Buffers
	vaoID uint32
	vboID uint32
	eboID uint32

	shader api.IShader

	modelLoc int32
	colorLoc int32

	dirty bool
}

// NewDynamicPixelAtlas create atlas that holds a bunch of pixels all
// of the same color.
func NewDynamicPixelAtlas(world api.IWorld) api.IAtlasX {
	o := new(dynamicPixelAtlas)
	o.world = world
	return o
}

func (s *dynamicPixelAtlas) Configure() error {
	err := s.configureShaders(s.world.RelativePath(), s.world.Properties())
	if err != nil {
		return err
	}

	return nil
}
func (s *dynamicPixelAtlas) SetData(vertices []float32, indices []uint32) {
	s.vertices = vertices
	s.indices = indices
}

func (s *dynamicPixelAtlas) Burn() error {
	err := s.Configure()
	if err != nil {
		return err
	}

	err = s.Bake()
	if err != nil {
		return err
	}

	return nil
}

func (s *dynamicPixelAtlas) Shake() {
}

// Bake finalizes the Atlas by "baking" the shapes into the buffers.
func (s *dynamicPixelAtlas) Bake() error {
	// ---------------------------------------------------------
	// BEGIN VAO Scope and generate buffer ids
	// ---------------------------------------------------------
	gl.GenVertexArrays(1, &s.vaoID)

	// This VAO bind stat the VAO Scope
	// Bind the Vertex Array Object first, then bind and set vertex buffer(s)
	// and attribute pointer(s).
	gl.BindVertexArray(s.vaoID)

	gl.GenBuffers(1, &s.vboID)

	gl.GenBuffers(1, &s.eboID)

	// The total buffer sizes are count of types (i.e floats or ints) times
	// the size of the type. Thus the size is in Bytes
	s.vboBufferSize = len(s.vertices) * api.XYZComponentCount * floatSize
	eboBufferSize := len(s.indices) * uintSize

	if s.vboBufferSize == 0 || eboBufferSize == 0 {
		return errors.New("dynamicPixelAtlas: VBO/EBO buffers are zero in size")
	}

	s.vboBind(s.vboBufferSize, s.vertices)

	s.eboBind(s.vboBufferSize, s.indices)

	// Count == (xyz=3) * sizeof(float32)=4 == 12 thus each
	// vertex is 12 bytes
	vertexSize := int32(api.XYZComponentCount) * int32(floatSize)

	// Some shaders may have normals and/or textures coords.
	// We only have one attribute (position) in the shader, so the
	// 'position' attribute defaults to zero.
	const positionIndex uint32 = 0

	// We can link the attribute position with the data in the vertexData array
	gl.VertexAttribPointer(positionIndex, int32(api.XYZComponentCount), gl.FLOAT, false, vertexSize, gl.PtrOffset(0))

	// ----------------------------------------------
	// END VAO Scope
	// ----------------------------------------------
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)

	err := s.configureUniforms()
	if err != nil {
		return err
	}

	return nil
}

func (s *dynamicPixelAtlas) Use() {
	s.shader.Use()
	gl.BindVertexArray(s.vaoID)
}

func (s *dynamicPixelAtlas) UnUse() {
	gl.BindVertexArray(0)
}

// SetColor sets the shader's color
func (s *dynamicPixelAtlas) SetColor(color []float32) {
	gl.Uniform4fv(s.colorLoc, 1, &color[0])
}

func (s *dynamicPixelAtlas) Update() {
	// Copy entire buffer even if just one element changed.
	if s.dirty {
		gl.BindBuffer(gl.ARRAY_BUFFER, s.vboID)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, s.vboBufferSize, gl.Ptr(s.vertices))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}
	s.dirty = false
}

func (s *dynamicPixelAtlas) SetPixelActiveCount(count int) {
	s.indicesCount = count
}

func (s *dynamicPixelAtlas) SetVertex(x, y float32, index int) {
	i := index * 3
	s.vertices[i] = x
	s.vertices[i+1] = y
	s.dirty = true
}

func (s *dynamicPixelAtlas) Render(id int, model api.IMatrix4) {
	gl.UniformMatrix4fv(s.modelLoc, 1, false, &model.Matrix()[0])

	gl.DrawElements(gl.POINTS, int32(s.indicesCount), uint32(gl.UNSIGNED_INT), gl.PtrOffset(0))
}

func (s *dynamicPixelAtlas) configureShaders(relativePath string, config *configuration.Properties) error {
	dataPath, err := filepath.Abs(relativePath)
	if err != nil {
		return err
	}

	shaders := config.Shaders

	s.shader = rendering.NewShader(shaders.DynamicPixelVertexShaderFile, shaders.DynamicPixelFragmentShaderFile)
	err = s.shader.Load(dataPath)
	if err != nil {
		return err
	}

	return nil
}

func (s *dynamicPixelAtlas) configureUniforms() error {
	s.shader.Use()

	program := s.shader.Program()

	s.modelLoc = gl.GetUniformLocation(program, gl.Str("model\x00"))
	if s.modelLoc < 0 {
		return errors.New("World: couldn't find 'model' uniform variable")
	}

	s.colorLoc = gl.GetUniformLocation(program, gl.Str("fragColor\x00"))
	if s.colorLoc < 0 {
		return errors.New("Couldn't find 'fragColor' uniform variable")
	}

	// Projection and View
	projLoc := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	if projLoc < 0 {
		return errors.New("NodeManager: couldn't find 'projection' uniform variable")
	}

	viewLoc := gl.GetUniformLocation(program, gl.Str("view\x00"))
	if viewLoc < 0 {
		return errors.New("NodeManager: couldn't find 'view' uniform variable")
	}

	pm := s.world.Projection().Matrix()
	gl.UniformMatrix4fv(projLoc, 1, false, &pm[0])

	vm := s.world.Viewspace().Matrix()
	gl.UniformMatrix4fv(viewLoc, 1, false, &vm[0])

	return nil
}

func (s *dynamicPixelAtlas) vboBind(bufferSize int, vertices []float32) {
	// the buffer type of a vertex buffer object is GL_ARRAY_BUFFER
	// From this point on any buffer calls we make (on the GL_ARRAY_BUFFER target)
	// will be used to configure the currently bound buffer, which is VBO
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vboID)

	// Create and fill buffer
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize, gl.Ptr(vertices), gl.DYNAMIC_DRAW)
}

func (s *dynamicPixelAtlas) eboBind(bufferSize int, indices []uint32) {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.eboID)

	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bufferSize, gl.Ptr(indices), gl.DYNAMIC_DRAW)
}
