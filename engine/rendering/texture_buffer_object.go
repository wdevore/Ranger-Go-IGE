package rendering

import (
	"image"
	"unsafe"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// TextureBufferObject associates an Atlas with a VAO
type TextureBufferObject struct {
	vao *VAO
	vbo *VBO
	tbo *TBO
	ebo *EBO

	floatSize int
	uintSize  int
}

// NewTextureBufferObject creates a new vector object with an associated Mesh
func NewTextureBufferObject() api.IBufferObject {
	o := new(TextureBufferObject)
	o.floatSize = int(unsafe.Sizeof(float32(0)))
	o.uintSize = int(unsafe.Sizeof(uint32(0)))
	return o
}

// Construct configures a buffer object
func (b *TextureBufferObject) Construct(meshType int, atlas api.IAtlas) {
}

// ConstructWithImage configures a buffer object
// meshType is ignored
func (b *TextureBufferObject) ConstructWithImage(image *image.NRGBA, smooth bool, atlas api.IAtlas) {
	b.vao = NewVAO()
	b.vao.BindStart()

	b.vbo = NewVBO(api.MeshStatic)
	b.tbo = NewTBO()
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

		for _, v := range shape.Vertices() {
			vertices = append(vertices, v)
		}

		for _, i := range shape.Indices() {
			indices = append(indices, indexOffset+uint32(i))
		}

		indexOffset = uint32(len(vertices) / api.XYZComponentCount)
	}

	vboBufferSize := len(vertices) * api.XYZWComponentCount * b.floatSize
	eboBufferSize := len(indices) * b.uintSize

	if vboBufferSize == 0 || eboBufferSize == 0 {
		panic("VBO.Construct: VBO/EBO buffers are zero")
	}

	b.tbo.BindTextureVbo2(vertices, b.vbo.VboID())

	b.ebo.Bind(eboBufferSize, indices)

	b.tbo.BindUsingImage(image, smooth)

	b.vao.UnUse()
}

// Use activates the VAO
func (b *TextureBufferObject) Use() {
	b.vao.Use()
	b.tbo.Use()
}

// UnUse deactivates the VAO
func (b *TextureBufferObject) UnUse() {
	b.tbo.UnUse()
	b.vao.UnUse()
}

// Update unused
func (b *TextureBufferObject) Update(shape api.IAtlasShape) {
}
