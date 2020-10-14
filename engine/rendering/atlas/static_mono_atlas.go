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
// A mono Atlas uses a single color for the whole shape verse per vertex.

type staticMonoAtlas struct {
	world api.IWorld

	shapes map[int]*shape

	nextID int

	// For the baking process
	vertices []float32
	indices  []uint32

	// Buffers
	vaoID uint32
	vboID uint32
	eboID uint32

	shader api.IShader

	modelLoc int32
	colorLoc int32
}

// NewStaticMonoAtlas create atlas that holds static shapes
// have a single (i.e. mono) color.
func NewStaticMonoAtlas(world api.IWorld) api.IAtlasX {
	o := new(staticMonoAtlas)
	o.shapes = make(map[int]*shape)
	o.world = world
	return o
}

func (s *staticMonoAtlas) Configure() error {
	err := s.configureShaders(s.world.RelativePath(), s.world.Properties())
	if err != nil {
		return err
	}

	return nil
}

// AddShape adds a set of vertices and indices to the atlas.
func (s *staticMonoAtlas) AddShape(shapeName string, vertices []float32, indices []uint32, mode int) int {
	shape := shape{
		shapeName:     shapeName,
		vertices:      vertices,
		indices:       indices,
		indicesCount:  len(indices),
		primitiveMode: uint32(mode),
	}

	id := s.nextID
	s.shapes[id] = &shape

	s.nextID++

	return id
}

func (s *staticMonoAtlas) GetShapeByName(shapeName string) int {
	for id, shape := range s.shapes {
		if shape.shapeName == shapeName {
			return id
		}
	}

	return -1
}

func (s *staticMonoAtlas) Burn() error {
	err := s.Configure()
	if err != nil {
		panic(err)
	}

	s.Shake()
	err = s.Bake()
	if err != nil {
		return err
	}

	return nil
}

func (s *staticMonoAtlas) Shake() {
	// ---------------------------------------------------------
	// Collect all vertices and indices for buffers
	// ---------------------------------------------------------
	// The atlas has shapes and each shape has vertices. These need to be
	// combined into a single array and later copied into GL Buffer.
	// At the same time each shape needs to be updated
	// to adjust element offsets and counts.
	indicesOffset := 0
	indiceBlockOffset := uint32(0)

	for _, shape := range s.shapes {
		shape.indicesOffset = indicesOffset

		indicesOffset += len(shape.indices) * uintSize

		s.vertices = append(s.vertices, shape.vertices...)

		// The indice offset is always refering to a position within
		// the vertices array.
		for _, i := range shape.indices {
			s.indices = append(s.indices, uint32(i)+indiceBlockOffset)
		}

		// Offset the indices based on the vertice block position.
		indiceBlockOffset = uint32(len(s.vertices) / api.XYZComponentCount)
	}
}

// Bake finalizes the Atlas by "baking" the shapes into the buffers.
func (s *staticMonoAtlas) Bake() error {
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
	vboBufferSize := len(s.vertices) * api.XYZComponentCount * floatSize
	eboBufferSize := len(s.indices) * uintSize

	if vboBufferSize == 0 || eboBufferSize == 0 {
		return errors.New("staticMonoAtlas: VBO/EBO buffers are zero in size")
	}

	s.vboBind(vboBufferSize, s.vertices)

	s.eboBind(vboBufferSize, s.indices)

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

	// dispose redundant arrays.
	s.vertices = nil
	s.indices = nil

	return nil
}

func (s *staticMonoAtlas) Use() {
	s.shader.Use()
	gl.BindVertexArray(s.vaoID)
}

func (s *staticMonoAtlas) UnUse() {
	// See opengl wiki as to why "glBindVertexArray(0)" isn't really necessary here:
	// https://www.opengl.org/wiki/Vertex_Specification#Vertex_Buffer_Object
	// Note the line "Changing the GL_ARRAY_BUFFER binding changes nothing about vertex attribute 0..."
	gl.BindVertexArray(0)
}

// SetColor sets the shader's color
func (s *staticMonoAtlas) SetColor(color []float32) {
	gl.Uniform4fv(s.colorLoc, 1, &color[0])
}

func (s *staticMonoAtlas) Render(id int, model api.IMatrix4) {
	shape := s.shapes[id]

	gl.UniformMatrix4fv(s.modelLoc, 1, false, &model.Matrix()[0])

	gl.DrawElements(shape.primitiveMode, int32(shape.indicesCount), uint32(gl.UNSIGNED_INT), gl.PtrOffset(shape.indicesOffset))
}

func (s *staticMonoAtlas) configureShaders(relativePath string, config *configuration.Properties) error {
	dataPath, err := filepath.Abs(relativePath)
	if err != nil {
		return err
	}

	shaders := config.Shaders

	s.shader = rendering.NewShader(shaders.MonoVertexShaderFile, shaders.MonoFragmentShaderFile)
	err = s.shader.Load(dataPath)
	if err != nil {
		return err
	}

	return nil
}

func (s *staticMonoAtlas) configureUniforms() error {
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

func (s *staticMonoAtlas) vboBind(bufferSize int, vertices []float32) {
	// the buffer type of a vertex buffer object is GL_ARRAY_BUFFER
	// From this point on any buffer calls we make (on the GL_ARRAY_BUFFER target)
	// will be used to configure the currently bound buffer, which is VBO
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vboID)

	// Create and fill buffer
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize, gl.Ptr(vertices), gl.STATIC_DRAW)
}

func (s *staticMonoAtlas) eboBind(bufferSize int, indices []uint32) {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.eboID)

	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bufferSize, gl.Ptr(indices), gl.STATIC_DRAW)
}
