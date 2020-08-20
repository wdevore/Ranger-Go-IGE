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

	shape        api.IAtlasShape
	textureIndex uint32
}

// NewTextureQuadObject creates a new vector object with an associated Mesh
func NewTextureQuadObject() api.IBufferObject {
	o := new(textureQuadObject)
	return o
}

// Construct configures a buffer object
func (b *textureQuadObject) Construct(meshType int, atlas api.IAtlas) {
}

// ConstructWithImage configures a buffer object
// meshType is ignored
func (b *textureQuadObject) ConstructWithImage(image *image.NRGBA, textureIndex uint32, smooth bool, atlas api.IAtlas) {
	b.textureIndex = textureIndex

	b.vao = NewVAO()
	b.vao.BindStart()

	b.vbo = NewVBO(api.MeshStatic)
	b.tbo = NewTBO()
	b.ebo = NewEBO()

	b.shape = atlas.Shapes()[0]
	vertices := b.shape.Vertices()
	indices := b.shape.Indices()

	vboBufferSize := len(*vertices) * api.XYZWComponentCount * int(unsafe.Sizeof(float32(0)))
	eboBufferSize := len(indices) * int(unsafe.Sizeof(uint32(0)))

	if vboBufferSize == 0 || eboBufferSize == 0 {
		panic("VBO.Construct: VBO/EBO buffers are zero")
	}

	b.vbo.BindTextureVbo2(vertices)

	b.ebo.Bind(eboBufferSize, indices)

	b.tbo.BindUsingImage(image, textureIndex, smooth)

	b.vao.UnUse()
}

// Use activates the VAO
func (b *textureQuadObject) Use() {
	b.vao.Use()
	b.tbo.Use(b.textureIndex)
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

	(*vertices)[2] = c[0]
	(*vertices)[3] = c[1]
	(*vertices)[6] = c[2]
	(*vertices)[7] = c[3]
	(*vertices)[10] = c[4]
	(*vertices)[11] = c[5]
	(*vertices)[14] = c[6]
	(*vertices)[15] = c[7]

	// stride := 4
	// i := 2
	// (*vertices)[i] = c[0]
	// (*vertices)[i+1] = c[1]
	// i += stride
	// (*vertices)[i] = c[2]
	// (*vertices)[i+1] = c[3]
	// i += stride
	// (*vertices)[i] = c[4]
	// (*vertices)[i+1] = c[5]
	// i += stride
	// (*vertices)[i] = c[6]
	// (*vertices)[i+1] = c[7]

	b.vbo.UpdateTexture(vertices)
}
