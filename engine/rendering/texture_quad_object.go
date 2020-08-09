package rendering

import (
	"image"
	"unsafe"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

type textureQuadObject struct {
	vao *VAO
	vbo *VBO
	tbo *TBO
	ebo *EBO

	floatSize int
	uintSize  int

	shape api.IAtlasShape
}

// NewTextureQuadObject creates a new vector object with an associated Mesh
func NewTextureQuadObject() api.IBufferObject {
	o := new(textureQuadObject)
	o.floatSize = int(unsafe.Sizeof(float32(0)))
	o.uintSize = int(unsafe.Sizeof(uint32(0)))
	return o
}

// Construct configures a buffer object
func (b *textureQuadObject) Construct(meshType int, atlas api.IAtlas) {
}

// ConstructWithImage configures a buffer object
// meshType is ignored
func (b *textureQuadObject) ConstructWithImage(image *image.NRGBA, smooth bool, atlas api.IAtlas) {
	b.vao = NewVAO()
	b.vao.BindStart()

	b.vbo = NewVBO(api.MeshStatic)
	b.tbo = NewTBO()
	b.ebo = NewEBO()

	b.shape = atlas.Shapes()[0]
	vertices := b.shape.Vertices()
	indices := b.shape.Indices()

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
func (b *textureQuadObject) Use() {
	b.vao.Use()
	b.tbo.Use(0)
}

// UnUse deactivates the VAO
func (b *textureQuadObject) UnUse() {
	b.tbo.UnUse()
	b.vao.UnUse()
}

// Update unused
func (b *textureQuadObject) Update(shape api.IAtlasShape) {
}

// TextureUpdate ...
func (b *textureQuadObject) TextureUpdate(coords *[]float32) {
	vertices := b.shape.Vertices()
	c := *coords
	stride := 4

	i := 2
	vertices[i] = c[0]
	vertices[i+1] = c[1]
	i += stride
	vertices[i] = c[2]
	vertices[i+1] = c[3]
	i += stride
	vertices[i] = c[4]
	vertices[i+1] = c[5]
	i += stride
	vertices[i] = c[6]
	vertices[i+1] = c[7]

	b.vbo.UpdateTexture(vertices)
}
