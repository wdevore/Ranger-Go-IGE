package rendering

import (
	"image"
	"unsafe"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// BufferObject associates an Atlas with a VAO
type BufferObject struct {
	vao *VAO
	vbo *VBO
	ebo *EBO

	floatSize int
	uintSize  int
}

// NewBufferObject creates a new vector object with an associated Mesh
func NewBufferObject() api.IBufferObject {
	o := new(BufferObject)
	o.floatSize = int(unsafe.Sizeof(float32(0)))
	o.uintSize = int(unsafe.Sizeof(uint32(0)))
	return o
}

// Construct configures a buffer object
func (b *BufferObject) Construct(meshType int, atlas api.IAtlas) {
	b.vao = NewVAO()
	b.vao.BindStart()

	b.vbo = NewVBO(meshType)
	b.ebo = NewEBO()

	// The atlas has shapes and each shape has vertices.
	// These need to be combined into a single array
	// and copied into GL Buffer.
	// At the same time each shape needs to be updated
	// to adjust offsets and counts for the EBO
	vertices := []float32{}
	indices := []uint32{}

	elementOffset := 0
	indexOffset := uint32(0)

	for _, shape := range atlas.Shapes() {
		shape.SetElementOffset(elementOffset)
		elementOffset += len(shape.Indices()) * b.uintSize

		for _, v := range *shape.Vertices() {
			vertices = append(vertices, v)
		}

		for _, i := range shape.Indices() {
			indices = append(indices, indexOffset+uint32(i))
		}

		indexOffset = uint32(len(vertices) / api.XYZComponentCount)
	}

	vboBufferSize := len(vertices) * api.XYZComponentCount * b.floatSize
	eboBufferSize := len(indices) * b.uintSize

	if vboBufferSize == 0 || eboBufferSize == 0 {
		panic("BO.Construct: VBO/EBO buffers are zero")
	}

	b.vbo.Bind(vboBufferSize, vertices)

	b.ebo.Bind(eboBufferSize, indices)

	b.vao.BindComplete()
}

// ConstructWithImage note need to rethink the api
func (b *BufferObject) ConstructWithImage(image *image.NRGBA, textureIndex uint32, smooth bool, atlas api.IAtlas) {
}

// Use activates the VAO
func (b *BufferObject) Use() {
	b.vao.Use()
}

// UnUse deactivates the VAO
func (b *BufferObject) UnUse() {
	b.vao.UnUse()
}

// Update modifies the VBO buffer
func (b *BufferObject) Update(shape api.IAtlasShape) {
	// b.mesh.Update(shape)
	b.vbo.Update(shape.Offset(), shape.Count(), shape.Vertices())
}
